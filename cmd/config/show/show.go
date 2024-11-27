package show

import (
	"fmt"

	"github.com/msetsma/RepoRover/core/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func CmdShowConfig(tool *util.CmdTool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Display the config values.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configData, err := tool.Config()
			if err != nil {
				fmt.Fprintln(tool.IOStreams.ErrOut, "Error loading configuration:", err)
				return err
			}
			yamlData, err := yaml.Marshal(configData)
			if err != nil {
				fmt.Fprintln(tool.IOStreams.ErrOut, "Failed to marshal configuration to YAML:", err)
				return err
			}
			fmt.Fprintln(tool.IOStreams.Out, string(yamlData))
			return nil
		},
	}

	return cmd
}
