package worker

import (
	"crawler/model"
	"crawler/engine"
	parser2 "crawler/zhenai/parser"
	"fmt"
	"log"
)

type SerializeParser struct {
	Name string
	Args interface{}
}

type Request struct {
	Url             string
	SerializeParser SerializeParser
}

type ParserResult struct {
	Items    []model.Item
	Requests []Request
}

func SerializeRequest(r engine.Request) Request {
	name, args := r.Parser.Name()
	return Request{
		Url: r.Url,
		SerializeParser: SerializeParser{
			Name: name,
			Args: args,
		},
	}
}

func DeSerializeRequest(r Request) (engine.Request,error) {

	parser,err := deserializeParser(r.SerializeParser)
	if err != nil {
		return engine.Request{},err
	}

	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	},nil
}

func SerializeParserResult(r engine.ParserResult) ParserResult {
	result := ParserResult{
		Items:r.Items,
	}

	for _,request := range r.Requests{
		result.Requests = append(result.Requests,SerializeRequest(request))
	}

	return result
}

func DeserializeParserResult(p ParserResult) engine.ParserResult {
	parserResult := engine.ParserResult{
		Items:p.Items,
	}

	for _,r := range p.Requests{
		request,err := DeSerializeRequest(r)
		if err != nil {
			log.Printf("deserialize parserResult error: %v,%v",r,err)
		}
		parserResult.Requests = append(parserResult.Requests,request)
	}


	return parserResult
}

func deserializeParser(parser SerializeParser) (engine.Parser,error) {
	switch parser.Name {
	case "ParserCityList":
		return engine.NewFuncParser(parser2.ParserCityList,"ParserCityList"),nil
	case "ParserCity":
		return engine.NewFuncParser(parser2.ParserCity,"ParserCity"),nil
	case "NilParser":
		return engine.NilParser{},nil
	case "ProfileParser":

		userName,ok := parser.Args.(string)
		if ok {
			return parser2.NewProfileParser(userName),nil
		}else {
			return nil,fmt.Errorf("error parser args %v",parser.Args)
		}
	default:
		return nil,fmt.Errorf("no find the parser %v",parser.Name)

	}
}
