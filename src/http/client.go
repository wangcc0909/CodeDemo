package main

import (
	"net/http"
	"net/http/httputil"
	"fmt"
)

func main() {

	//创建一个请求
	request,err := http.NewRequest(http.MethodGet,"http://www.baidu.com",nil)
	if err != nil {
		panic(err)
	}

	//添加请求头
	request.Header.Add("User-Agent",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")

	//创建一个用户
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("Redirect :",req)
			return nil
		},
	}

	//发送请求
	response,err:= client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	//解析数据
	result,err := httputil.DumpResponse(response,true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n",result)

}
