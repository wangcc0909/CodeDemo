package engine

import (
	"log"
	"crawler/fetcher"
)

func Worker(r Request) (ParserResult,error)  {

	log.Printf("fetcher url %s",r.Url)
	contents,err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("fetcher code error %s",err)
		return  ParserResult{},err
	}

	//这里解析是用的取出来的Request的解析器
	return r.Parser.Parse(contents,r.Url),nil
}