var deleteGroupCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing group",
	Run: func(cmd *cobra.Command, args []string) {
		if groupName == "" {
			fmt.Println("Error: --name flag is required to delete a group")
			return
		}

		err := core.DeleteGroup(groupName)
		if err != nil {
			fmt.Printf("Error deleting group '%s': %v\n", groupName, err)
			return
		}

		fmt.Printf("Group '%s' deleted successfully!\n", groupName)
	},
}