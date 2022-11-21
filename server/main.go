package main

import (
	"context"
	"net"

	pb "github.com/SalviCF/authorization-server/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedKahosServiceServer
}

func (s *server) read(context.Context, *pb.ReadPoliciesTenantReq) (*pb.ReadPoliciesTenantResp, error) {
	return nil, nil
}

func (s *server) create(context.Context, *pb.CreatePolicyTenantReq) (*pb.CreatePolicyTenantResp, error) {
	return nil, nil
}

func (s *server) delete(context.Context, *pb.DeletePolicyTenantReq) (*pb.DeletePolicyTenantResp, error) {
	return nil, nil
}

func (s *server) update(context.Context, *pb.UpdatePolicyTenantReq) (*pb.UpdatePolicyTenantResp, error) {
	return nil, nil
}

func (s *server) readroles(context.Context, *pb.ReadRolesTenantReq) (*pb.ReadRolesTenantResp, error) {
	return nil, nil
}

func (s *server) createrole(context.Context, *pb.CreateRoleTenantReq) (*pb.CreateRoleTenantResp, error) {
	return nil, nil

}
func (s *server) deleterole(context.Context, *pb.DeleteRoleTenantReq) (*pb.DeletePolicyTenantResp, error) {
	return nil, nil
}

func (s *server) updaterole(context.Context, *pb.UpdateRoleTenantReq) (*pb.UpdatePolicyTenantResp, error) {
	return nil, nil
}

func main() {
	listner, err := net.Listen("tcp", ":50051")

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	serv := grpc.NewServer()
	pb.RegisterKahosServiceServer(serv, &server{})
	if err = serv.Serve(listner); err != nil {
		panic("cannot initialize the server" + err.Error())
	}

}
