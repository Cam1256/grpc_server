package main

import (
	"context"
	"net"

	pb "github.com/SalviCF/authorization-server/proto"
	"google.golang.org/grpc"
)

func generateID() string {
	rand.Seed(time.Now().Unix())
	return "ID: " + strconv.Itoa(rand.Int())
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		panic("cannot connect with server " + err.Error())
	}

	serviceClient := pb.NewKahosServiceClient(conn)

	err := serviceClient.Create(context.Background(), &pb.ReadPoliciesTenantReq{
	

	
	}
}