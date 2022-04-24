package grpc_server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Litekube/network-controller/certs"
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/contant"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	"github.com/Litekube/network-controller/internal"
	"github.com/Litekube/network-controller/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
)

type GrpcServer struct {
	*pb_gen.UnimplementedLiteKubeVpnServiceServer
	port          int
	UnRegisterCh  chan string
	service       *internal.LiteVpnService
	grpcTlsConfig config.TLSConfig
	vpnTlsConfig  config.TLSConfig
}

var LiteVpnSocket = "unix://litevpn.sock"
var logger = utils.GetLogger()
var gServer *GrpcServer

func newGrpcServer(cfg config.ServerConfig, unRegisterCh chan string) *GrpcServer {
	s := &GrpcServer{
		port:         cfg.GrpcPort,
		UnRegisterCh: unRegisterCh,
		grpcTlsConfig: config.TLSConfig{
			CAFile:         filepath.Join(cfg.GrpcCertDir, contant.CAFile),
			CAKeyFile:      filepath.Join(cfg.GrpcCertDir, contant.CAKeyFile),
			ServerCertFile: filepath.Join(cfg.GrpcCertDir, contant.ServerCertFile),
			ServerKeyFile:  filepath.Join(cfg.GrpcCertDir, contant.ServerKeyFile),
			ClientCertFile: filepath.Join(cfg.GrpcCertDir, contant.ClientCertFile),
			ClientKeyFile:  filepath.Join(cfg.GrpcCertDir, contant.ClientKeyFile),
		},
		vpnTlsConfig: config.TLSConfig{
			CAFile:         filepath.Join(cfg.VpnCertDir, contant.CAFile),
			CAKeyFile:      filepath.Join(cfg.VpnCertDir, contant.CAKeyFile),
			ServerCertFile: filepath.Join(cfg.VpnCertDir, contant.ServerCertFile),
			ServerKeyFile:  filepath.Join(cfg.VpnCertDir, contant.ServerKeyFile),
			ClientCertFile: filepath.Join(cfg.VpnCertDir, contant.ClientCertFile),
			ClientKeyFile:  filepath.Join(cfg.VpnCertDir, contant.ClientKeyFile),
		},
	}
	s.service = internal.NewLiteVpnService(unRegisterCh, s.grpcTlsConfig, s.vpnTlsConfig)
	return s
}

func StartGrpcServer(cfg config.ServerConfig, unRegisterCh chan string) {
	gServer = newGrpcServer(cfg, unRegisterCh)
	utils.CreateDir(cfg.GrpcCertDir)
	err := certs.CheckGrpcCertConfig(gServer.grpcTlsConfig)
	if err != nil {
		logger.Error(err)
	}
	err = gServer.startGrpcServerTcp()
	if err != nil {
		logger.Error(err)
	}
}

func (s *GrpcServer) startGrpcServerTcp() error {
	tcpAddr := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		logger.Errorf("tcp failed to listen: %v", err)
		return err
	}

	gopts := []grpc.ServerOption{}
	if len(s.grpcTlsConfig.ServerCertFile) != 0 && len(s.grpcTlsConfig.ServerKeyFile) != 0 {
		creds, err := credentials.NewServerTLSFromFile(s.grpcTlsConfig.ServerCertFile, s.grpcTlsConfig.ServerKeyFile)
		if err != nil {
			logger.Error(err)
			return err
		}
		gopts = append(gopts, grpc.Creds(creds))
	}
	cert, err := tls.LoadX509KeyPair(s.grpcTlsConfig.ServerCertFile, s.grpcTlsConfig.ServerKeyFile)
	//cert, err := certificate.LoadCertificate(s.CertFile)
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(s.grpcTlsConfig.CAFile)
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	gopts = append(gopts, grpc.Creds(creds))

	server := grpc.NewServer(gopts...)
	// register reflection for grpcurl service
	reflection.Register(server)
	// register service
	pb_gen.RegisterLiteKubeVpnServiceServer(server, s)
	logger.Infof("grpc server ready to serve at %+v", tcpAddr)
	if err := server.Serve(lis); err != nil {
		logger.Errorf("grpc server failed to serve: %v", err)
		return err
	}
	return nil
}

func (s *GrpcServer) startGrpcServerUDS() error {
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
