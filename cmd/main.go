// bb-cli: A command-line interface for Bitbucket pull requests.
//
// Authentication:
//
//	App password (most common):
//	  export BITBUCKET_USERNAME=myuser
//	  export BITBUCKET_APP_PASSWORD=ATBBxxxxxxxx
//
//	OAuth2 access token:
//	  export BITBUCKET_TOKEN=<token>
//
// Usage examples:
//
//	bb-cli pr list --workspace myorg --repo myrepo --state OPEN
//	bb-cli pr list --workspace myorg --repo myrepo --all
//	bb-cli pr get  --workspace myorg --repo myrepo --id 42
//	bb-cli pr create --workspace myorg --repo myrepo --title "My feature" \
//	  --source-branch feature/x --destination-branch main
//	bb-cli pr merge --workspace myorg --repo myrepo --id 42 --strategy squash
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/FabianSchurig/bitbucket-cli/internal/commands"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bb-cli",
		Short: "Bitbucket CLI",
		Long: `bb-cli is a command-line interface for Bitbucket Cloud.

Set authentication environment variables before running:
  App passwords (most common):
    BITBUCKET_USERNAME    your Bitbucket username
    BITBUCKET_APP_PASSWORD  your app password

  OAuth2 access token:
    BITBUCKET_TOKEN       your OAuth2 access token`,
	}

	prCmd := commands.NewPRCommand()
	commands.AddOutputFlag(prCmd)

	rootCmd.AddCommand(prCmd)
	return rootCmd
}
