package {{CommandName}}

import (
	"fmt"

	"github.com/msetsma/RepoRover/core"
	"github.com/spf13/cobra"
)

// NewCmd{{CommandName}} creates a new cobra command for {{CommandName}}
func NewCmd{{CommandName}}() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "{{CommandUse}}",
		Short: "{{command description}}",
		Run: func(cmd *cobra.Command, args []string) {
			// add logic here
		},
	}

	return cmd
}
