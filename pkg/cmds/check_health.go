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

var healthTemplate = template.Must(template.New("network-controller CheckHealth").Parse(`
------------------------------------------------
network-controller:
    control grpc client health: {{.CtrlHealth}}
    bootstrap grpc client health: {{.BootHealth}}
------------------------------------------------
`))

func NewCheckHealthCommand() *cli.Command {
	return &cli.Command{
		Name:      "check-health",
		Usage:     "check health of control and bootstrap grpc",
		UsageText: "ncadm [global options] check-health",
		Action:    checkHealth,
	}
}

func checkHealth(ctx *cli.Context) error {
	client := NewClient()
	if client == nil {
		return errors.New("fail to init gRPC client")
	}

	req := &pb_gen.HealthCheckRequest{}
	backCtx := context.Background()

	ctrlResp, err := client.GClient.C.HealthCheck(backCtx, req)
	if err != nil {
		return err
	}

	bootResp, err := client.BootstrapClient.BootstrapC.HealthCheck(backCtx, req)
	if err != nil {
		return err
	}

	checker := func(code string) string {
		if code == contant.STATUS_OK {
			return "Health"
		}
		return "UnHealth"
	}

	data := struct {
		CtrlHealth string
		BootHealth string
	}{
		CtrlHealth: checker(ctrlResp.Code),
		BootHealth: checker(bootResp.Code),
	}

	healthTemplate.Execute(os.Stdout, &data)

	return nil
}
