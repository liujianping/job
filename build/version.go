package build

import "fmt"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

//String version
func String() string {
	return fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
}

//Info build info
func Info(v, m, d string) {
	version = v
	commit = m
	date = d
}
