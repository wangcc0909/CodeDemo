package main

import (
	"net/http"
	"errhandingserver/filelistingserver/filelisting"
	"log"
	"os"
)

type appHander func(http.ResponseWriter, *http.Request) error

type UserError interface {
	error
	Message() string
}

func errWrapper(hander appHander) func(http.ResponseWriter, *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic %v",r)
				http.Error(writer,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
			}
		}()

		err := hander(writer,request)

		if err != nil {
			log.Printf("Error handing request %s",err.Error())

			if userErr,ok := err.(UserError); ok {
				http.Error(writer,userErr.Message(),http.StatusBadRequest)
				return
			}

			code := http.StatusOK

			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}

			http.Error(writer,http.StatusText(code),code)
		}
	}
}

func main() {
	http.HandleFunc("/", errWrapper(filelisting.Hander))

	err := http.ListenAndServe(":8888",nil)
	if err != nil {
		panic(err)
	}

}
