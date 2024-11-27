package init

import (
	"github.com/msetsma/RepoRover/core/config"
	"github.com/msetsma/RepoRover/core/util"
	"github.com/spf13/cobra"
)

func CmdGroupInit(tool *util.CmdTool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Display the group values.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := tool.Config()
			if err != nil {
				return err
			}
			cfg.ActiveGroup = args[0]
			return config.Update(cfg)
		},
	}

	return cmd
}
