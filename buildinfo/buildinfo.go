package buildinfo

import "fmt"

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func Info() string {
	return fmt.Sprintf("%s (%s, %s)", Version, Commit, BuildDate)
}

func Short() string {
	return Version
}
