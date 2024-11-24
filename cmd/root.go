package cmd

import (
	"os"

	"github.com/MakeNowJust/heredoc"
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

func NewCmdRoot() (*cobra.Command, error) {
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
	// Child commands
	// cmd.AddCommand(versionCmd.NewCmdVersion(f, version, buildDate))

	return cmd, nil
}

func Execute() exitCode {
	root, err := NewCmdRoot()
	if err != nil {
		return exitError
	}
	err = root.Execute()
	if err != nil {
		return exitError
	}
    return exitOK
}
