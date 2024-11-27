package main

import (
	"os"

	"github.com/msetsma/RepoRover/cmd"
	"github.com/msetsma/RepoRover/core/util"
)

func main() {
	tool := util.NewCmdTool()
	code := cmd.Run(tool)
	os.Exit(int(code))
}
