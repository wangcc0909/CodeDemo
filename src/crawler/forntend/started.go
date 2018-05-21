package main

import (
	"net/http"
	"crawler/forntend/controller"
)

func main() {

	http.Handle("/",http.FileServer(http.Dir("src/crawler/forntend/view")))

	handle := controller.CreateSearchResultHandle(
		"src/crawler/forntend/view/template.html")
	http.Handle("/search",handle)
	err := http.ListenAndServe(":8888",nil)
	if err != nil {
		panic(err)
	}
}
