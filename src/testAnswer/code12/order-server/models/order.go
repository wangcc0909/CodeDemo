package models

import "github.com/globalsign/mgo/bson"

type Order struct {
	OrderId            bson.ObjectId
	OrderNo            string
	ProductId          uint32
	Count              uint32
	UserId             uint32
	Status             uint32
	ProductAmountTotal float32
	OrderAmountTotal   float32
	LogisticsFee       float32
	AddressId          uint32
	Remark             string
}
