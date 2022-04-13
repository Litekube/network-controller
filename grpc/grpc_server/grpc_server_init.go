package grpc_server

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"strings"
	"ws-vpn/grpc/pb_gen"
	"ws-vpn/utils"
)

var LiteVpnSocket = "unix://litevpn.sock"
var logger = utils.GetLogger()
var gServer *grpcServer

func StartGrpcServer(port int, unRegisterCh chan string) {
	gServer = newGrpcServer(port, unRegisterCh)
	gServer.startGrpcServerTcp()
}

func (s *grpcServer) startGrpcServerTcp() error {
	tcpAddr := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		logger.Errorf("tcp failed to listen: %v", err)
		return err
	}
	server := grpc.NewServer()
	// 注册 grpcurl 所需的 reflection 服务
	reflection.Register(server)
	// 注册业务服务
	pb_gen.RegisterLiteKubeVpnServiceServer(server, s)
	logger.Infof("grpc server ready to serve at %+v", tcpAddr)
	if err := server.Serve(lis); err != nil {
		logger.Errorf("grpc server failed to serve: %v", err)
		return err
	}
	return nil
}

func (s *grpcServer) startGrpcServerUDS() error {
	os.Remove("/tmp/litevpn.sock")
	server_addr, err := net.ResolveUnixAddr("unix", "/tmp/litevpn.sock")
	if err != nil {
		logger.Errorf("failed to resolve unix addr err:%+v")
		return err
	}
	fmt.Println(server_addr)
	lis, err := net.ListenUnix("unix", server_addr)
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		return err
	}

	gs := grpc.NewServer()
	pb_gen.RegisterLiteKubeVpnServiceServer(gs, s)
	err = gs.Serve(lis)
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		return err
	}
	return nil
}

// createListener returns a listener bound to the requested protocol and address.
func (s *grpcServer) createListener(listener string) (ret net.Listener, rerr error) {
	if listener == "" {
		listener = LiteVpnSocket
	}
	network, address := s.networkAndAddress(listener)

	if network == "unix" {
		if err := os.Remove(address); err != nil && !os.IsNotExist(err) {
			logger.Warningf("failed to remove socket %s: %v", address, err)
		}
		defer func() {
			if err := os.Chmod(address, 0600); err != nil {
				rerr = err
			}
		}()
	} else {
		network = "tcp"
	}

	return net.Listen(network, address)
}

// networkAndAddress crudely splits a URL string into network (scheme) and address,
// where the address includes everything after the scheme/authority separator.
func (s *grpcServer) networkAndAddress(str string) (string, string) {
	parts := strings.SplitN(str, "://", 2)
	if len(parts) > 1 {
		return parts[0], parts[1]
	}
	return "", parts[0]
}
