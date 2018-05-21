package main

import (
	"crawler/persist"
	"crawler/engine"
	"crawler/simple"
	"crawler/zhenai/parser"
)

func main() {

	itemChan,err := persist.ItemServer("dating_profile")
	if err != nil {
		panic(err)
	}

	e := engine.CurrentEngine{
		Scheduler:&simple.QueuedScheduler{},
		WorkerCount:100,
		ItemServer:itemChan,
		Processor:engine.Worker,
	}

	e.Run(engine.Request{
		Url:"http://www.zhenai.com/zhenghun",
		Parser:engine.NewFuncParser(parser.ParserCityList,"ParserCityList"),
	})

	/*e.Run(engine.Request{
		Url:"http://www.zhenai.com/zhenghun/beijing",
		ParserFunc:parser.ParserCity,
	})*/
	/*client, err := elastic.NewClient(

		elastic.SetSniff(false),
		elastic.SetURL("http://192.168.99.100:9200"),
	)

	if err != nil {
		panic(err)
	}

	_,err = client.DeleteIndex("dating_profile").DoC(context.Background())
	if err != nil {
		panic(err)
	}*/
}
