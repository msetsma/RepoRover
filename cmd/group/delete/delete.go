package delete

import (
	"github.com/spf13/cobra"
)

var deleteGroupCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing group",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
