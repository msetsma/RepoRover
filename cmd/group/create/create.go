package cmd

import (
	"fmt"

	"github.com/msetsma/RepoRover/core"
	"github.com/spf13/cobra"
)

func NewCmdCreate() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new group",
		Run: func(cmd *cobra.Command, args []string) {
			if groupName == "" {
				fmt.Println("Error: --name flag is required to create a group")
				return
			}

			err := core.CreateGroup(groupName)
			if err != nil {
				fmt.Printf("Error creating group: %v\n", err)
				return
			}

			fmt.Printf("Group '%s' created successfully!\n", groupName)
			if defaultGroup {
				err := core.SetDefaultGroup(groupName)
				if err != nil {
					fmt.Printf("Error setting default group: %v\n", err)
					return
				}
				fmt.Printf("Group '%s' is now the default group.\n", groupName)
			}
		},
	}

	return cmd

}