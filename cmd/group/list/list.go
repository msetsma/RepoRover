var listGroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups",
	Run: func(cmd *cobra.Command, args []string) {
		groups, err := core.ListGroups()
		if err != nil {
			fmt.Printf("Error listing groups: %v\n", err)
			return
		}

		if len(groups) == 0 {
			fmt.Println("No groups found.")
			return
		}

		fmt.Println("Groups:")
		for _, group := range groups {
			fmt.Printf("- %s\n", group)
		}
	},
}