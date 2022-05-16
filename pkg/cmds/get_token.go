package cmds

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/Litekube/network-controller/contant"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	certutil "github.com/rancher/dynamiclistener/cert"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"text/template"
)

var gettokenTemplate = template.Must(template.New("network-controller GetToken").Parse(`
------------------------------------------------
network-controller:
    BootstrapToken: {{.BootstrapToken}}
    NetworkServerIp: {{.NetworkServerIp}}
    NetworkServerPort: {{.NetworkServerPort}}
    GrpcServerIp: {{.GrpcServerIp}}
    GrpcServerPort: {{.GrpcServerPort}}
    NetworkCertsDir: {{.NetworkCertsDir}}
    GrpcCertsDir: {{.GrpcCertsDir}}
------------------------------------------------
`))

func NewGetTokenCommand() *cli.Command {
	return &cli.Command{
		Name:      "get-token",
		Usage:     "get grpc server ip/port/certs",
		UsageText: "ncadm [global options] get-token [options]",
		Action:    getToken,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "bootstrap-token",
				Usage:       "uid of bootstrap token",
				Destination: &bootstrapToken,
			},
			&cli.StringFlag{
				Name:        "network-certs-dir",
				Usage:       "dir of network ca/client cert/client key file",
				Destination: &networkCertsDir,
			},
			&cli.StringFlag{
				Name:        "grpc-certs-dir",
				Usage:       "dir of grpc ca/client cert/client key file",
				Destination: &grpcCertsDir,
			},
		},
	}
}

func getToken(ctx *cli.Context) error {
	client := NewClient()
	if client == nil {
		return errors.New("fail to init gRPC client")
	}

	req := &pb_gen.GetTokenRequest{
		BootStrapToken: bootstrapToken,
	}
	resp, err := client.BootstrapClient.BootstrapC.GetToken(context.Background(), req)
	if err != nil {
		return err
	} else if resp.Code != "200" {
		return errors.New(resp.Message)
	}

	caBytes, err := base64.StdEncoding.DecodeString(resp.GrpcCaCert)
	certBytes, err := base64.StdEncoding.DecodeString(resp.GrpcClientCert)
	keyBytes, err := base64.StdEncoding.DecodeString(resp.GrpcClientKey)
	certutil.WriteCert(filepath.Join(grpcCertsDir, contant.CAFile), caBytes)
	certutil.WriteCert(filepath.Join(grpcCertsDir, contant.ClientCertFile), certBytes)
	certutil.WriteKey(filepath.Join(grpcCertsDir, contant.ClientKeyFile), keyBytes)

	caBytes, err = base64.StdEncoding.DecodeString(resp.NetworkCaCert)
	certBytes, err = base64.StdEncoding.DecodeString(resp.NetworkClientCert)
	keyBytes, err = base64.StdEncoding.DecodeString(resp.NetworkClientKey)
	certutil.WriteCert(filepath.Join(networkCertsDir, contant.CAFile), caBytes)
	certutil.WriteCert(filepath.Join(networkCertsDir, contant.ClientCertFile), certBytes)
	certutil.WriteKey(filepath.Join(networkCertsDir, contant.ClientKeyFile), keyBytes)

	data := struct {
		BootstrapToken    string
		NetworkServerIp   string
		NetworkServerPort string
		GrpcServerIp      string
		GrpcServerPort    string
		NetworkCertsDir   string
		GrpcCertsDir      string
	}{
		BootstrapToken:    bootstrapToken,
		NetworkServerIp:   resp.NetworkServerIp,
		NetworkServerPort: resp.NetworkServerPort,
		GrpcServerIp:      resp.GrpcServerIp,
		GrpcServerPort:    resp.GrpcServerPort,
		NetworkCertsDir:   networkCertsDir,
		GrpcCertsDir:      grpcCertsDir,
	}

	gettokenTemplate.Execute(os.Stdout, &data)

	return nil
}
