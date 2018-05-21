package client

import (
	"crawler/engine"
	"crawler_distributed/worker"
	"net/rpc"
)

func WorkProcessor(client chan *rpc.Client) engine.Processor  {

	return func(request engine.Request) (engine.ParserResult, error) {

		/*client,err := supportrpc.ClientRpc(host)
		if err != nil {
			return engine.ParserResult{}, nil
		}*/

		result := worker.ParserResult{}

		sReq := worker.SerializeRequest(request)
		c := <- client
		err := c.Call("WorkServer.Worker",sReq,&result)
		if err != nil {
			return engine.ParserResult{}, err
		}
		return worker.DeserializeParserResult(result),nil
	}
}