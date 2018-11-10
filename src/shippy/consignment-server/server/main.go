package main

import (
	pb "shippy/consignment-server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment,error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	Consignments []*pb.Consignment
}

func (r *Repository) Create(consignment *pb.Consignment) (*pb.Consignment,error) {
	r.Consignments = append(r.Consignments,consignment)
	return consignment,nil
}

func (r *Repository) GetAll() []*pb.Consignment {
	return r.Consignments
}

type service struct {
	repo Repository
}

func (s *service) CreateConsignment(ctx context.Context, in *pb.Consignment) (*pb.Respone,error) {
	consignment,err := s.repo.Create(in)
	if err != nil {
		return nil, err
	}

	resp := &pb.Respone{
		Created:true,
		Consignment:consignment,
	}
	return resp,nil
}

func (s *service) GetConsignments(ctx context.Context,in *pb.Request) (*pb.Respone,error) {
	consignments := s.repo.GetAll()
	return &pb.Respone{Consigments:consignments},nil
}

const (
	Address = ":8899"
)

func main() {
	listener,err := net.Listen("tcp",Address)
	if err != nil {
		log.Fatalf("net listener error %v",err)
	}

	s := grpc.NewServer()

	repo := Repository{}

	pb.RegisterShippingServiceServer(s,&service{repo:repo})
	if err := s.Serve(listener);err != nil {
		log.Fatalf("server serve error %v",err)
	}
}
