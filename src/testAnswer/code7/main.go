package main

import (
	"net/http"
	"github.com/skip2/go-qrcode"
	"log"
	"net/url"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"encoding/json"
	"time"
)

func main() {

	http.HandleFunc("/",handle)
	http.HandleFunc("/crawel",handleCrawel)
	http.ListenAndServe(":8080",nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	qr,err := qrcode.Encode("今天周三",qrcode.Highest,256)
	if err != nil {
		log.Panic(err)
		w.Write([]byte("err"))
		return
	}
	w.Write(qr)
}

func handleCrawel(w http.ResponseWriter, r *http.Request)  {
	_,err := url.Parse("https://www.tianapi.com/")
	if err != nil {
		log.Panic(err)
		return
	}
	resp,err := http.Get("https://www.tianapi.com/")
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println(resp.Body)
	defer resp.Body.Close()
	document,err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Panic(err)
		return
	}
	var datas = make(map[string]interface{})
	document.Find("li a").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
		datas[selection.Text()] = selection.Text()
	})
	time.Sleep(5 * time.Second)
	var result []byte
	result,err = json.Marshal(datas)
	if err != nil {
		log.Panic(err)
		return
	}
	w.Write(result)
}