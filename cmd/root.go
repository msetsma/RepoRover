package cmd

import (
	"fmt"
	"os"

	config "github.com/msetsma/RepoRover/cmd/config"
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

func cmdRoot() (*cobra.Command, error) int {
	// rootCmd represents the base command when called without any subcommands
	cmd := &cobra.Command{
		Use:   "rr <command> <subcommand> [flags]",
		Short: "RepoRover CLI",
		Long: "Explore groups of repos!",
		Example: heredoc.Doc(`
				$ rr groups list
				$ rr sync group -n myrepos
				$ rr config set-default
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
	// Child commands
	// cmd.AddCommand(versionCmd.NewCmdVersion(f, version, buildDate))
	// cmd.AddCommand(actionsCmd.NewCmdActions(f))
	// cmd.AddCommand(aliasCmd.NewCmdAlias(f))
	// cmd.AddCommand(authCmd.NewCmdAuth(f))
	// cmd.AddCommand(attestationCmd.NewCmdAttestation(f))
	// cmd.AddCommand(configCmd.NewCmdConfig(f))
	// cmd.AddCommand(creditsCmd.NewCmdCredits(f, nil))
	// cmd.AddCommand(gistCmd.NewCmdGist(f))
	// cmd.AddCommand(gpgKeyCmd.NewCmdGPGKey(f))
	// cmd.AddCommand(completionCmd.NewCmdCompletion(f.IOStreams))
	// cmd.AddCommand(extensionCmd.NewCmdExtension(f))
	// cmd.AddCommand(searchCmd.NewCmdSearch(f))
	// cmd.AddCommand(secretCmd.NewCmdSecret(f))
	// cmd.AddCommand(variableCmd.NewCmdVariable(f))
	// cmd.AddCommand(sshKeyCmd.NewCmdSSHKey(f))
	// cmd.AddCommand(statusCmd.NewCmdStatus(f, nil))
	// cmd.AddCommand(codespaceCmd.NewCmdCodespace(f))
	// cmd.AddCommand(projectCmd.NewCmdProject(f))


}


// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitError)
	}
}