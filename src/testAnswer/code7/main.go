package main

import (
	"net/http"
	"github.com/skip2/go-qrcode"
	"log"
	"net/url"
	"io/ioutil"
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

//1.获取歌曲名和演唱者
// 2.将演唱者和歌曲名以参数的形式替换ur中的keyword,获取FileHash
// 3.FileHash就是 一下注释url中的hash参数
var ur = "http://songsearch.kugou.com/song_search_v2?callback=jQuery1124006980366032059648_1518578518932&keyword=%E5%88%9A%E5%88%9A%E5%A5%BD&page=1&pagesize=30&userid=-1&clientver=&platform=WebFilter&tag=em&filter=2&iscorrection=1&privilege_filter=0&_=1518578518934"
//http://www.kugou.com/yy/index.php?r=play/getdata&hash=8E5DDAC9C06A6469ED500F18985D56D6
func handleCrawel(w http.ResponseWriter, r *http.Request)  {
	_,err := url.Parse(ur)
	if err != nil {
		log.Panic(err)
		return
	}
	resp,err := http.Get(ur)
	if err != nil {
		log.Panic(err)
		return
	}
	defer resp.Body.Close()
	result,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	log.Println(string(result))
	w.Write(result)
}