package config

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/msetsma/RepoRover/core/util"
	"github.com/spf13/cobra"

	setConfigValueCmd "github.com/msetsma/RepoRover/cmd/config/set"
	showConfigCmd "github.com/msetsma/RepoRover/cmd/config/show"
)

func NewCmdConfig(tool *util.CmdTool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage config",
		Long:  `Make changes to the configuration of RepoRover`,
		Example: heredoc.Doc(`
			$ rr config show
			$ rr config set -n <group name>
		`),
	}

	cmd.AddCommand(showConfigCmd.CmdShowConfig(tool))
	cmd.AddCommand(setConfigValueCmd.CmdSetConfig(tool))

	return cmd
}
