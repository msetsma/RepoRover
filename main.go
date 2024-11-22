package main

import (
	"github.com/msetsma/RepoRover/cmd"
)

func main() {
	code := cmd.Execute()
	os.Exit(int(code))
}