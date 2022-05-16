package cmds

import (
	"context"
	"errors"
	"fmt"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	"github.com/urfave/cli/v2"
	"os"
	"text/template"
)

var printTemplate = template.Must(template.New("network-controller bootstrap").Parse(`
------------------------------------------------
network-controller:
    token: {{.BootstrapToken}}@{{.IP}}:{{.Port}}
    ExpireMsg: {{.ExpireMsg}}
------------------------------------------------
`))

func NewCreateTokenCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-bootstrap-token",
		Usage:     "create network bootstrap token info",
		UsageText: "ncadm [global options] create-bootstrap-token [options]",
		Action:    createBootstrapToken,
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "life",
				Usage:       "how long will bootstrap token info be valid",
				Destination: &life,
				Value:       10,
			},
		},
	}
}

func createBootstrapToken(ctx *cli.Context) error {
	client := NewClient()
	if client == nil {
		return errors.New("fail to init gRPC client")
	}

	req := &pb_gen.GetBootStrapTokenRequest{
		ExpireTime: life,
	}
	resp, err := client.GClient.C.GetBootStrapToken(context.Background(), req)
	if err != nil {
		return err
	} else if resp.Code != "200" {
		return errors.New(resp.Message)
	}

	var expireMsg string
	if life > 0 {
		expireMsg = fmt.Sprintf("expire in %d min", life)
	} else if life < 0 {
		expireMsg = "no expire"
	} else {
		expireMsg = "no msg"
	}

	data := struct {
		BootstrapToken string
		IP             string
		Port           string
		ExpireMsg      string
	}{
		BootstrapToken: resp.BootStrapToken,
		IP:             resp.CloudIp,
		Port:           resp.Port,
		ExpireMsg:      expireMsg,
	}

	printTemplate.Execute(os.Stdout, &data)

	return nil
}
