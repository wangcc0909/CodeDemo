package model

import "github.com/globalsign/mgo/bson"

type Duanzi struct {
	Id      bson.ObjectId `bson:"_id" json:"id"`
	Title   string        `bson:"title" json:"title"`
	Content string        `bson:"content" json:"content"`
}

const (
	db = "duanzi"
	collection = "duanziModel"
)

func InsertDuanzi(duanzi Duanzi) error {
	return Insert(db,collection,duanzi)
}

func FindAllDuanzi() ([]Duanzi, error) {
	var dzs []Duanzi
	err := FindAll(db,collection,nil,nil,&dzs)
	return dzs,err
}