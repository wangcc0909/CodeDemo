package main

import (
	"net/http"
	"log"
	"flag"
	"io/ioutil"
	"crypto/x509"
	"crypto/tls"
	"golang.org/x/net/http2"
	"fmt"
)

var url = "https://localhost:8000"

var httpVersion = flag.Int("version",2,"HTTP version")

func main() {
	flag.Parse()
	client := http.Client{}
	caCert,err := ioutil.ReadFile("src/http/server.crt")
	if err != nil {
		log.Panic(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		RootCAs:certPool,
	}

	switch *httpVersion {
	case 1:
		client.Transport = &http.Transport{
			TLSClientConfig:config,
		}
	case 2:
		client.Transport = &http2.Transport{
			TLSClientConfig:config,
		}
	}
	r,err := client.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer r.Body.Close()
	result,err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(result))
}
