package config

import (
	"fmt"

	defaultCmd "github.com/msetsma/RepoRover/cmd/config/default"
	showCmd "github.com/msetsma/RepoRover/cmd/config/show"
	"github.com/msetsma/RepoRover/utils"
	"github.com/msetsma/RepoRover/core"
	"github.com/spf13/cobra"
)

func NewCmdConfig () (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage config",
		Long:  `Make changes to the configuration of RepoRover`,
		Example: heredoc.Doc(`
			$ rr config show
			$ rr config default -n <group name>
		`)
	}

	cmd.AddCommand(defaultCmd.NewCmdConfigDefault())
	cmd.AddCommand(showCmd.NewCmdConfigShow())

	return cmd, nil
}