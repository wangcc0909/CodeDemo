package persist

import (
	"testing"
	"crawler/model"
	"gopkg.in/olivere/elastic.v3"
	"golang.org/x/net/context"
	"encoding/json"
)

func TestItemServer(t *testing.T)  {

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

	newEx := model.Item{
		Url:"http://album.zhenai.com/u/105993901",
		Id:"105993901",
		Type:"zhenai",
		Profile:expected,
	}



	client,err:= elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://192.168.99.100:9200"),
	)
	if err != nil {
		panic(err)
	}

	err = SaveItem(client,"test",newEx)
	if err != nil {
		panic(err)
	}

	reslut,err := client.Get().Index("test").
		Type("zhenai").Id(newEx.Id).DoC(context.Background())

	if err != nil {
		panic(err)
	}

	actual := model.Item{}

	err = json.Unmarshal(*reslut.Source,&actual)
	if err != nil {
		panic(err)
	}
	if newEx != actual{
		t.Errorf("expected %v,actual %v",expected,actual)
	}



}