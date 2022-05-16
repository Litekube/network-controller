package cmds

import (
	"github.com/Litekube/network-controller/grpc/grpc_client"
)

type Client struct {
	GClient         *grpc_client.GrpcClient
	BootstrapClient *grpc_client.GrpcBootStrapClient
}

var life int64
var nodeToken string
var bootstrapToken string
var networkCertsDir string
var grpcCertsDir string

func NewClient() *Client {
	client := &Client{
		GClient: &grpc_client.GrpcClient{
			Ip:       GlobalConfig.ip,
			Port:     GlobalConfig.port,
			CAFile:   GlobalConfig.CAFile,
			CertFile: GlobalConfig.CertFile,
			KeyFile:  GlobalConfig.KeyFile,
		},
		BootstrapClient: &grpc_client.GrpcBootStrapClient{
			Ip:            GlobalConfig.ip,
			BootstrapPort: GlobalConfig.bootStrapPort,
		},
	}
	// init grpc client
	if err := client.GClient.InitGrpcClientConn(); err != nil {
		panic(err)
	}

	// init bootstrap grpc client
	if err := client.BootstrapClient.InitGrpcBootstrapClientConn(); err != nil {
		panic(err)
	}

	return client
}

func (c *Client) GRPC() *grpc_client.GrpcClient {
	return c.GClient
}
