package set

import (
	//"github.com/msetsma/RepoRover/core/config"
	"github.com/msetsma/RepoRover/core/util"
	"github.com/spf13/cobra"
)

func CmdSetConfig(cfg *util.CmdTool) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set values in config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd

}
