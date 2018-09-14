package main

import (
	"net/http"
	"log"
)

func main() {
	hand := handleStruct{}
	ser := &http.Server{Addr:":8000",Handler:hand}
	log.Fatal(ser.ListenAndServeTLS("src/http/server.crt","src/http/server.key"))
}

type handleStruct struct {
}

func (h handleStruct) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("HELLO"))
}
