package group

import (
	initGroupCmd "github.com/msetsma/RepoRover/cmd/group/init"
	"github.com/MakeNowJust/heredoc"
	"github.com/msetsma/RepoRover/core/util"
	"github.com/spf13/cobra"
)

func NewCmdGroup(tool *util.CmdTool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group <command>",
		Short: "Manage groups",
		Long:  `Make changes, get information on groups.`,
		Example: heredoc.Doc(`
			$ rr group list
			$ rr group delete -n <group name>
		`),
		GroupID: "group",
	}

	cmd.AddCommand(initGroupCmd.CmdGroupInit(tool))

	return cmd
}
