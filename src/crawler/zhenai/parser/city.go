package parser

import (
	"crawler/engine"
	"regexp"
)

var (
	cityRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
	cityUrlRe = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/beijing/[^"]+)"`)
)
func ParserCity(contents []byte,_ string) engine.ParserResult {

	matches := cityRe.FindAllSubmatch(contents,-1)
	result := engine.ParserResult{}

	for _,m := range matches {
		//result.Items = append(result.Items,"User " + string(m[2]))
		result.Requests = append(result.Requests,engine.Request{
			Url:string(m[1]),
			Parser: NewProfileParser(string(m[2])),
		})
	}

	urls := cityUrlRe.FindAllSubmatch(contents,-1)
	for _,m := range urls{
		result.Requests = append(result.Requests,engine.Request{
			Url:string(m[1]),
			Parser: engine.NewFuncParser(ParserCity,"ParserCity"),
		})

	}

	return result
}
