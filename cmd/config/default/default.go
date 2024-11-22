package default

import (
	"fmt"

	"github.com/msetsma/RepoRover/core"
	"github.com/spf13/cobra"
)

func NewCmdConfigDefault() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "default",
		Short: "Set the default group",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("Error: Group name is required to set as default")
				return
			}

			groupName := args[0]
			err := core.SetDefaultGroup(groupName)
			if err != nil {
				fmt.Printf("Error setting default group: %v\n", err)
				return
			}

			fmt.Printf("Group '%s' is now the default group.\n", groupName)
		},
	}

	return cmd

}