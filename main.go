package main

import (
	"os"

	"github.com/msetsma/RepoRover/cmd"
)

func main() {
	code := cmd.Execute()
	os.Exit(int(code))
}
