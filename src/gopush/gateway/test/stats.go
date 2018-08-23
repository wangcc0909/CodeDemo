package main

import (
	"net/http"
	"crypto/tls"
	"os"
	"fmt"
	"io/ioutil"
)

func main() {
	var (
		err    error
		resp   *http.Response
		buf    []byte
		client *http.Client
	)
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	if resp, err = client.Get("https://localhost:7788/stats"); err != nil {
		goto ERR
	}
	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		goto ERR
	}
	fmt.Println("返回值:" + string(buf))
	return
ERR:
	fmt.Println(err)
	os.Exit(-1)
	return
}
