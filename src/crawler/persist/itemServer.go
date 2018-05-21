package persist

import (
	"log"
	"gopkg.in/olivere/elastic.v3"
	"golang.org/x/net/context"
	"crawler/model"
)

func ItemServer(index string) (chan model.Item,error) {

	out := make(chan model.Item)

	client, err := elastic.NewClient(

		elastic.SetSniff(false),
		elastic.SetURL("http://192.168.99.100:9200"),
	)

	if err != nil {

		return nil, err
	}

	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("get Item %d %v", itemCount, item)
			itemCount++

			err := SaveItem(client,index,item)
			if err != nil {
				panic(err)
			}
		}
	}()

	return out,nil
}

func SaveItem(client *elastic.Client,index string,item model.Item) error {

	_, err := client.Index().
		Index(index).
		Type(item.Type).
		Id(item.Id).
		BodyJson(item).
		DoC(context.Background())

	return err

}
