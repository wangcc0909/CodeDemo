package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

//将一个int64数转换成[]byte
func IntToHex(num int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer,binary.BigEndian,num)//将num的binary编码写入buffer中
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
