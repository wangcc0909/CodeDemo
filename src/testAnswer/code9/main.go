package main

import (
	"testAnswer/code9/crawler"
	"fmt"
)

func main() {
	result := crawler.CrawlerJianshu("https://www.jianshu.com/p/340c3f251d24")
	fmt.Println(result)
}
