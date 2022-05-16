package cmds

import (
	"context"
	"errors"
	"github.com/Litekube/network-controller/contant"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	"github.com/urfave/cli/v2"
	"os"
	"text/template"
)

var checkTemplate = template.Must(template.New("network-controller CheckConnState").Parse(`
------------------------------------------------
network-controller:
    node-token: {{.NodeToken}}
    BindIp: {{.BindIp}}
    ConnState: {{.StateMsg}}
------------------------------------------------
`))

func NewCheckConnStateCommand() *cli.Command {
	return &cli.Command{
		Name:      "check-conn-state",
		Usage:     "check network conn state",
		UsageText: "ncadm [global options] check-conn-state [options]",
		Action:    checkConnState,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "node-token",
				Usage:       "uid of node conn",
				Destination: &nodeToken,
			},
		},
	}
}

func checkConnState(ctx *cli.Context) error {
	client := NewClient()
	if client == nil {
		return errors.New("fail to init gRPC client")
	}

	req := &pb_gen.CheckConnStateRequest{
		Token: nodeToken,
	}
	resp, err := client.GClient.C.CheckConnState(context.Background(), req)
	if err != nil {
		return err
	} else if resp.Code != "200" {
		return errors.New(resp.Message)
	}

	var stateMsg string
	switch resp.ConnState {
	case contant.STATE_IDLE:
		stateMsg = "UnConnected"
		break
	case contant.STATE_INIT:
		stateMsg = "Init Connecting"
		break
	case contant.STATE_CONNECTED:
		stateMsg = "Connected"
		break
	default:
		stateMsg = "invalid state"
	}

	data := struct {
		NodeToken string
		BindIp    string
		StateMsg  string
	}{
		NodeToken: nodeToken,
		BindIp:    resp.BindIp,
		StateMsg:  stateMsg,
	}

	checkTemplate.Execute(os.Stdout, &data)

	return nil
}
