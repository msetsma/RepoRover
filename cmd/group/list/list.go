package delete

import (
	"github.com/spf13/cobra"
)

var listGroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
