package grpc_server

import (
	"context"
	"ws-vpn/grpc/pb_gen"
	"ws-vpn/internal"
)

type grpcServer struct {
	*pb_gen.UnimplementedLiteKubeVpnServiceServer
	port         int
	unRegisterCh chan string
	service      *internal.LiteVpnService
}

func newGrpcServer(port int, unRegisterCh chan string) *grpcServer {
	return &grpcServer{
		port:         port,
		unRegisterCh: unRegisterCh,
		service:      internal.NewLiteVpnService(unRegisterCh),
	}
}

func (s *grpcServer) HelloWorld(ctx context.Context, req *pb_gen.HelloWorldRequest) (*pb_gen.HelloWorldResponse, error) {
	logger.Infof("get helloworld request: %+v", req)
	reply := &pb_gen.HelloWorldResponse{ThanksText: "hello,this wanna"}
	return reply, nil
}

func (s *grpcServer) CheckConnState(ctx context.Context, req *pb_gen.CheckConnStateRequest) (*pb_gen.CheckConnResponse, error) {
	logger.Infof("get CheckConnState request: %+v", req)
	return s.service.CheckConnState(ctx, req)
}

func (s *grpcServer) UnRegister(ctx context.Context, req *pb_gen.UnRegisterRequest) (*pb_gen.UnRegisterResponse, error) {
	logger.Infof("get CheckConnState request: %+v", req)
	return s.service.UnRegister(ctx, req)
}
