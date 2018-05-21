package main

import (
	"interfaces/mock"
	"fmt"
)

type Retrieve interface {
	Get(url string) string
}

func download(r Retrieve) string {
	return r.Get(address)
}

type Poster interface {
	Post(url string,form map[string]string) string
}

func post(poster Poster) {
	result := poster.Post(address,map[string]string{
		"contents": "Contents",
		"name":"names",
	})

	fmt.Println(result)
}

type RetrievePost interface {
	Retrieve
	Poster
}

func session(rp RetrievePost) string {
	rp.Post(address,map[string]string{
		"contents":"Contents new",
		"name":"name new",
	})
	return rp.Get(address)
}

const address = "http://www.imooc.com"

func main() {
	r := mock.Retrieve{Contents: "this is a contents"}
	fmt.Println(download(&r))

	//r2 := real2.Retrieve{}
	//fmt.Println(download(r2))

	post(&r)

	fmt.Println()

	fmt.Println(session(&r))


}
