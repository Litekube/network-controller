package version

var (
	NetworkController = "0.1.0"
	GitBranch         = "default-main" // branch of git
	GitVersion        = "v2.25.1"
	GitCommit         = "$HEAD"                 // sha1 from git, output of $(git rev-parse HEAD)
	BuildDate         = "2022-04-10 T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	Version           = "0.1.0"
)
