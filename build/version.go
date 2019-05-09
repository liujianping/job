package build

import "fmt"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

//Print version
func Print() {
	fmt.Printf("%v, commit %v, built at %v\n", version, commit, date)
}

//Info build info
func Info(v, m, d string) {
	version = v
	commit = m
	date = d
}
