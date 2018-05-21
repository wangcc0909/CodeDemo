package worker

import "crawler/engine"

type WorkServer struct {

}

func (w WorkServer) Worker(request Request,result *ParserResult) error {
	r,err := DeSerializeRequest(request)
	if err != nil {
		return err
	}
	p,err := engine.Worker(r)
	if err != nil {
		return err
	}
	*result = SerializeParserResult(p)
	return nil
}
