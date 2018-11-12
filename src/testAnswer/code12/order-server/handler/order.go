package handler

import (
	"golang.org/x/net/context"
	pb "testAnswer/code12/order-server/proto"
	"testAnswer/code12/common"
	"testAnswer/code12/order-server/models"
	"github.com/globalsign/mgo/bson"
	"time"
)

const (
	db         = "order"
	collection = "orderModel"
)

type Service struct {
}

func (s Service) CreateOrder(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	//这里通过mongodb创建数据
	var order models.Order
	order.UserId = request.UserId
	order.Count = request.Count
	order.AddressId = request.AddressId
	order.LogisticsFee = request.LogisticsFee
	order.OrderAmountTotal = request.OrderAmountTotal
	order.ProductAmountTotal = request.ProductAmountTotal
	order.ProductId = request.ProductId
	order.Remark = request.Remark
	order.OrderId = bson.NewObjectId()
	order.OrderNo = "ff12"
	order.Status = 0
	err := common.Insert(db,collection,&order)
	if err != nil {
		return nil,err
	}

	return &pb.Response{
		OrderId:order.OrderId.String(),
		OrderNo:order.OrderNo,
		UserId:order.UserId,
		Status:order.Status,
		CreatedAt:time.Now().String(),
	},nil

}
