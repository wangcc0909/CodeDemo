package main

import (
	"testing"
	"crawler_distributed/supportrpc"
	"crawler/model"
	"time"
)

func TestItemServer(t *testing.T) {
	const host = ":2345"
	//support server
	go saveRpc(host,"test1")

	time.Sleep(time.Second)

	//support client
	client,err := supportrpc.ClientRpc(host)
	if err != nil {
		panic(err)
	}
	//call save
	expected := model.Profile{
		Name:"负二代",
		Age:25,
		Height:171,
		Weight:0,
		InComing:"12001-20000元",
		Gender:"男",
		Occupation:"--",
		Education:"中专",
		HuKou:"广东深圳",
		Car:"未购车",
		Horse:"--",
		Marriage:"未婚",
		XinZuo:"巨蟹座",
	}

	item := model.Item{
		Url:"http://album.zhenai.com/u/105993901",
		Id:"105993901",
		Type:"zhenai",
		Profile:expected,
	}

	var result = ""
	err = client.Call("ItemSaveServer.Save",item,&result)
	if err != nil || result != "ok"{
		t.Errorf("Save item server error %s , item %v",err,item)
	}

}