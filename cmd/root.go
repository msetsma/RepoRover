package cmd

import (
	"github.com/MakeNowJust/heredoc"
	CmdConfig "github.com/msetsma/RepoRover/cmd/config"
	CmdGroup "github.com/msetsma/RepoRover/cmd/group"
	"github.com/msetsma/RepoRover/core/util"
	"github.com/spf13/cobra"
)

type exitCode int

const (
	exitOK      exitCode = 0
	exitError   exitCode = 1
	exitCancel  exitCode = 2
	exitAuth    exitCode = 4
	exitPending exitCode = 8
)

func CmdRoot(tool *util.CmdTool) (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	cmd := &cobra.Command{
		Use:   "rr <command> <subcommand> [flags]",
		Short: "RepoRover CLI",
		Long:  "Manage groups of repos!",
		Example: heredoc.Doc(`
				$ rr groups list
				$ rr sync group -n myrepos
				$ rr config set-default -n <group_name>
		`),
	}

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.AddGroup(&cobra.Group{
		ID:    "config",
		Title: "config commands",
	})
	cmd.AddGroup(&cobra.Group{
		ID:    "sync",
		Title: "sync repo/groups commands",
	})
	cmd.AddGroup(&cobra.Group{
		ID:    "group",
		Title: "group commands",
	})
	cmd.AddGroup(&cobra.Group{
		ID:    "status",
		Title: "status commands",
	})
	// Example adding commands
	cmd.AddCommand(CmdConfig.NewCmdConfig(tool))
	cmd.AddCommand(CmdGroup.NewCmdGroup(tool))

	//

	return cmd, nil
}

func Run(tool *util.CmdTool) exitCode {
	root, err := CmdRoot(tool)
	if err != nil {
		return exitError
	}
	err = root.Execute()
	if err != nil {
		return exitError
	}
	return exitOK
}
