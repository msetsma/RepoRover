package group

import (
	"fmt"

	defaultCmd "github.com/msetsma/RepoRover/cmd/config/default"
	showCmd "github.com/msetsma/RepoRover/cmd/config/show"
	"github.com/msetsma/RepoRover/core"
	"github.com/spf13/cobra"
)

func NewCmdGroup () *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage config",
		Long:  `Make changes to the configuration of RepoRover`,
		Example: heredoc.Doc(`
			$ rr group list
			$ rr group delete -n <group name>
		`)
		GroupID: "group",
	},


}