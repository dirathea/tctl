package main

import "github.com/dirathea/tctl/cmd"

var (
	version = "v0.0.0"
)

func main() {
	cmd.Execute(version)
}
