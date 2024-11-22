package show

import (
	"fmt"

	"github.com/msetsma/RepoRover/core"
	"github.com/spf13/cobra"
)

func NewCmdConfigShow() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Display the entire config.toml file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return core.PrintConfig()
		},
	}

	return cmd

}