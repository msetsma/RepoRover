package activate

import (
	"github.com/msetsma/RepoRover/core/config"
	"github.com/msetsma/RepoRover/core/util"
	"github.com/spf13/cobra"
)

func CmdSetActiveGroup(tool *util.CmdTool) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "active",
		Short: "Set the active group",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd

}
