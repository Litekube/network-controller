package cmds

import (
	"context"
	"errors"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	"github.com/urfave/cli/v2"
	"os"
	"text/template"
)

var unregisterTemplate = template.Must(template.New("network-controller UnRegister").Parse(`
------------------------------------------------
network-controller:
    node-token: {{.NodeToken}}
    Success: {{.Result}}
------------------------------------------------
`))

func NewUnRegisterCommand() *cli.Command {
	return &cli.Command{
		Name:      "unregister",
		Usage:     "close network connection unregister bind ip",
		UsageText: "ncadm [global options] unregister [options]",
		Action:    unregister,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "node-token",
				Usage:       "uid of node conn",
				Destination: &nodeToken,
			},
		},
	}
}

func unregister(ctx *cli.Context) error {
	client := NewClient()
	if client == nil {
		return errors.New("fail to init gRPC client")
	}

	req := &pb_gen.UnRegisterRequest{
		Token: nodeToken,
	}
	resp, err := client.GClient.C.UnRegister(context.Background(), req)
	if err != nil {
		return err
	} else if resp.Code != "200" {
		return errors.New(resp.Message)
	}

	data := struct {
		NodeToken string
		Result    bool
	}{
		NodeToken: nodeToken,
		Result:    resp.Result,
	}

	unregisterTemplate.Execute(os.Stdout, &data)

	return nil
}
