package main

import (
	"fmt"
	"net/http"
	"io"
	"regexp"
	"strings"
	"strconv"
	"os"
)

//https://www.pengfu.com/xiaohua_3.html

func httpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}

	defer resp.Body.Close()

	var buf = make([]byte, 1024*4)
	for {
		n, err2 := resp.Body.Read(buf)

		if n == 0 {
			if err2 == io.EOF {
				break
			} else {
				err = err2
				return
			}
		}
		result += string(buf[:n])
	}
	return
}

func CrawlerOne(url string) (title, content string, err error) {
	result, err1 := httpGet(url)
	if err1 != nil {
		fmt.Println("httpGet err")
		err = err1
		return
	}

	reg := regexp.MustCompile(`<h1>(.*)</h1>`)
	temp := reg.FindAllStringSubmatch(result, 1)
	for _, tempTitle := range temp {

		title = tempTitle[1]
		title = strings.Replace(title, "\t", "", -1)
	}

	reg2 := regexp.MustCompile(`<div class="content-txt pt10">([^<]+)<a id="prev" href="`)
	temp2 := reg2.FindAllStringSubmatch(result, -1)
	for _, tempContent := range temp2 {

		content = tempContent[1]
		content = strings.Replace(content,"\t","",-1)
		content = strings.Replace(content,"\r","",-1)
		content = strings.Replace(content,"\n","",-1)
	}
	return title, content, nil
}

//<h1 class="dp-b"><a href="https://www.pengfu.com/content_1835036_1.html" target="_blank">多人正在等车</a>
func CrawlerPage(i int, page chan int) {
	url := "https://www.pengfu.com/xiaohua_" + fmt.Sprintf("%d", i) + ".html"
	result, err := httpGet(url)
	if err != nil {
		fmt.Println("httpGet err")
		page <- i
		return
	}
	fileName := strconv.Itoa(i) + ".txt"
	file,err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	reg := regexp.MustCompile(`<h1 class="dp-b"><a href="(.*)" target="_blank">`)
	temp := reg.FindAllStringSubmatch(result, -1)
	for _, content := range temp {
		title, joyContent, err := CrawlerOne(content[1])
		if err != nil {
			fmt.Println(err)
			continue
		}

		file.Write([]byte(title+ "\n"))
		file.Write([]byte(joyContent + "\n"))
		file.Write([]byte("================================\n"))
	}

	page <- i
}

func DoWork(start, end int) {
	var page = make(chan int)

	for i := start; i <= end; i++ {
		go CrawlerPage(i, page)
	}

	for i := start; i <= end; i++ {
		fmt.Printf("第%d页爬取完成\n", <-page)
	}
}

func main() {
	var start, end int
	fmt.Println("请输入开始页面( >=1)")

	fmt.Scan(&start)

	fmt.Println("请输入结束页面( >= 开始页面)")
	fmt.Scan(&end)

	DoWork(start, end)
}
