package model

import (
	"log"
	"github.com/globalsign/mgo"
)

const (
	host = "192.168.99.100:27017"
	source = "admin"
	user = "peaut"
	pass = "123456"
)

var globals *mgo.Session

func init() {
	dialInfo := &mgo.DialInfo{
		Addrs:[]string{host},
		Source:source,
		Username:user,
		Password:pass,
	}

	s,err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalln("create session error ",err)
	}

	globals = s
}

func Connect(db,collection string) (*mgo.Session,*mgo.Collection) {
	s := globals.Copy()
	c := s.DB(db).C(collection)
	return s,c
}

func Insert(db, collection string, docs ...interface{}) error {
	ms,c := Connect(db,collection)
	defer ms.Close()
	return c.Insert(docs...)
}

func FindOne(db, collection string, query, selector, result interface{}) error {
	ms,c := Connect(db,collection)
	defer ms.Close()

	return c.Find(query).Select(selector).One(result)
}

func FindAll(db, collection string, query, selector, result interface{}) error {
	ms,c := Connect(db,collection)
	defer ms.Close()
	return c.Find(query).Select(selector).All(result)
}

func Update(db, collection string, selector, update interface{}) error {
	ms,c := Connect(db,collection)
	defer ms.Close()
	return c.Update(selector,update)
}

func Remove(db, collection string, query interface{}) error {
	ms,c := Connect(db,collection)
	defer ms.Close()
	return c.Remove(query)
}
