// bb-cli is a command-line interface for Bitbucket Cloud pull requests.
//
// Install:
//
//	go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@latest
//
// Authentication:
//
//	App password (most common):
//	  export BITBUCKET_USERNAME=myuser
//	  export BITBUCKET_APP_PASSWORD=ATBBxxxxxxxx
//
//	OAuth2 access token:
//	  export BITBUCKET_TOKEN=<token>
package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/FabianSchurig/bitbucket-cli/internal/commands"
)

// Set via ldflags at build time (see goreleaser.yaml).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "bb-cli",
		Short:   "Bitbucket CLI",
		Version: version,
		Long: `bb-cli is a command-line interface for Bitbucket Cloud.

Set authentication environment variables before running:
  App passwords (most common):
    BITBUCKET_USERNAME    your Bitbucket username
    BITBUCKET_APP_PASSWORD  your app password

  OAuth2 access token:
    BITBUCKET_TOKEN       your OAuth2 access token`,
	}

	setColoredHelp(rootCmd)

	prCmd := commands.NewPRCommand()
	commands.AddOutputFlag(prCmd)

	hooksCmd := commands.NewHooksCommand()
	commands.AddOutputFlag(hooksCmd)

	searchCmd := commands.NewSearchCommand()
	commands.AddOutputFlag(searchCmd)

	refsCmd := commands.NewRefsCommand()
	commands.AddOutputFlag(refsCmd)

	commitsCmd := commands.NewCommitsCommand()
	commands.AddOutputFlag(commitsCmd)

	reportsCmd := commands.NewReportsCommand()
	commands.AddOutputFlag(reportsCmd)

	rootCmd.AddCommand(prCmd, hooksCmd, searchCmd, refsCmd, commitsCmd, reportsCmd)
	return rootCmd
}

func setColoredHelp(cmd *cobra.Command) {
	bold := color.New(color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	cobra.AddTemplateFunc("bold", bold)
	cobra.AddTemplateFunc("yellow", yellow)
	cobra.AddTemplateFunc("cyan", cyan)

	cmd.SetUsageTemplate(`{{bold "Usage:"}}{{if .Runnable}}
  {{yellow .UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{yellow .CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

{{bold "Aliases:"}} {{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

{{bold "Available Commands:"}}{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{cyan (rpad .Name .NamePadding)}} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{bold .Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{cyan (rpad .Name .NamePadding)}} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

{{bold "Additional Commands:"}}{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{cyan (rpad .Name .NamePadding)}} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{bold "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

{{bold "Global Flags:"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{yellow (print .CommandPath " [command] --help")}}" for more information about a command.{{end}}
`)
}
