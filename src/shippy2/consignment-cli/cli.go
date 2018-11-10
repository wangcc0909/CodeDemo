package main

import (
	pb "shippy2/consignment-server"
	"log"
	"golang.org/x/net/context"
	"io/ioutil"
	"encoding/json"
	"os"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/client"
)

func parseFile(filePath string) (*pb.Consignment, error) {
	bytes,err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil,err
	}

	var consignment *pb.Consignment
	err = json.Unmarshal(bytes,&consignment)
	if err != nil {
		return nil, err
	}
	return consignment,nil
}

func main() {

	dir,err := os.Getwd()
	if err != nil {
		log.Fatalf("os getWd error %v",err)
	}
	consignment,err := parseFile(dir + "\\src\\shippy\\consignment-cli\\consignment.json")
	if err != nil {
		log.Fatalf("parsefile error %v",err)
	}

	cmd.Init()
	c := pb.NewShippingServiceClient("go.micro.srv.consignment",client.DefaultClient)
	resp,err := c.CreateConsignment(context.Background(),consignment)
	if err != nil {
		log.Fatalf("create consignment error %v",err)
	}
	log.Printf("created = %v",resp.Created)

	resp1,err := c.GetConsignments(context.Background(),&pb.Request{})
	if err != nil {
		log.Fatalf("get all error %v",err)
	}

	log.Println(resp1.Consigments)
}
