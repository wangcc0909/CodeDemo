package main

import (
	pb "shippy2/consignment-server"
	"golang.org/x/net/context"
	"log"
	"github.com/micro/go-micro"
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

func (s *service) CreateConsignment(ctx context.Context, in *pb.Consignment,resp *pb.Respone) error{
	consignment,err := s.repo.Create(in)
	if err != nil {
		return err
	}

	resp = &pb.Respone{
		Created:true,
		Consignment:consignment,
	}
	return nil
}

func (s *service) GetConsignments(ctx context.Context,in *pb.Request,resp *pb.Respone) error {
	consignments := s.repo.GetAll()
	resp = &pb.Respone{Consigments:consignments}
	return nil
}

const (
	Address = ":50051"
)

func main() {
	server := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	server.Init()
	repo := Repository{}

	pb.RegisterShippingServiceHandler(server.Server(),&service{repo:repo})
	if err := server.Run();err != nil {
		log.Fatalf("server serve error %v",err)
	}
}
