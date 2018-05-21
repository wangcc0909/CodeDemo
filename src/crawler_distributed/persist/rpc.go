package persist

import (
	"crawler/model"
	"crawler/persist"
	"gopkg.in/olivere/elastic.v3"
	"log"
)

type ItemSaveServer struct {
	Client *elastic.Client
	Index  string
}

func (s *ItemSaveServer) Save(item model.Item, result *string) error {
	err := persist.SaveItem(s.Client, s.Index, item)

	if err == nil {
		*result = "ok"
	}else {
		log.Printf("Save Server err %v",err)
	}
	log.Printf("item %v saved",item)
	return err
}
