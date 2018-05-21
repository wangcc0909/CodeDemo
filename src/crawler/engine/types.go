package engine

import "crawler/model"

type ParserFunc func([]byte,string) ParserResult

type Parser interface {
	Parse([]byte,string) ParserResult
	Name() (name string,args interface{})
}

type Request struct {
	Url string
	Parser Parser
}

type ParserResult struct {
	Requests []Request
	Items []model.Item
}

type NilParser struct {

}

func (NilParser) Parse([]byte, string) ParserResult {
	return ParserResult{}
}

func (NilParser) Name() (name string, args interface{}) {
	return "NilParser", nil
}

type FuncParser struct {
	ParserFunc ParserFunc
	name string
}

func NewFuncParser(parse ParserFunc,name string) *FuncParser{
	return &FuncParser{
		ParserFunc:parse,
		name:name,
	}
}

func (f *FuncParser) Parse(content []byte,url string) ParserResult {
	return f.ParserFunc(content,url)
}

func (f *FuncParser) Name() (name string, args interface{}) {
	return f.name,nil
}


