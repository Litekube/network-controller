package version

import (
	"fmt"
	"runtime"
)

type Info struct {
	NcAdm             string `json:"ncadm"`
	NetworkController string `json:"network-controller"`
	GitVersion        string `json:"gitVersion"`
	GitBranch         string `json:"gitBranch"`
	GitCommit         string `json:"gitCommit"`
	BuildDate         string `json:"buildDate"`
	GoVersion         string `json:"goVersion"`
	Compiler          string `json:"compiler"`
	Platform          string `json:"platform"`
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in pkg/version/base.go
	return Info{
		NetworkController: NetworkController,
		NcAdm:             Version,
		GitVersion:        GitVersion,
		GitBranch:         GitBranch,
		GitCommit:         GitCommit,
		BuildDate:         BuildDate,
		GoVersion:         runtime.Version(),
		Compiler:          runtime.Compiler,
		Platform:          fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func PrintAndExitIfRequested() {
	fmt.Printf("%#v\n", Get())
}
