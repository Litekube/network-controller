package grpc_client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/Litekube/litekube-vpn/grpc/pb_gen"
	"github.com/Litekube/litekube-vpn/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"path/filepath"
)

type GrpcClient struct {
	c           pb_gen.LiteKubeVpnServiceClient
	Ip          string
	Port        string
	GrpcCertDir string
	CAFile      string
	CertFile    string
	KeyFile     string
}

var logger = utils.GetLogger()

func (c *GrpcClient) InitGrpcClientConn() error {
	// Set up a connection to the server.
	var address string
	if len(c.Ip) == 0 || len(c.Port) == 0 {
		logger.Error("ip and port can't be empty")
		return errors.New("ip and port can't be empty")
	}
	address = fmt.Sprintf("%s:%s", c.Ip, c.Port)

	var dialOpt []grpc.DialOption
	cert, err := tls.LoadX509KeyPair(filepath.Join(c.GrpcCertDir, c.CertFile), filepath.Join(c.GrpcCertDir, c.KeyFile))
	if err != nil {
		logger.Errorf("tls.LoadX509KeyPair err: %v", err)
		return err
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(filepath.Join(c.GrpcCertDir, c.CAFile))
	if err != nil {
		logger.Errorf("ioutil.ReadFile err: %v", err)
		return err
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		logger.Errorf("certPool.AppendCertsFromPEM err")
		return err
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   c.Ip,
		RootCAs:      certPool,
	})
	dialOpt = append(dialOpt, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(address, dialOpt...)
	if err != nil {
		logger.Errorf("can't connect: %v", err)
		return err
	}
	// init grpc client
	c.c = pb_gen.NewLiteKubeVpnServiceClient(conn)

	return nil
}
