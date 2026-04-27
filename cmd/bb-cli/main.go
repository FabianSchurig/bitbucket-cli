// bb-cli is a command-line interface for Bitbucket Cloud pull requests.
//
// Install:
//
//	go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-cli@latest
//
// Authentication:
//
//	API token (recommended):
//	  export BITBUCKET_USERNAME=myuser
//	  export BITBUCKET_TOKEN=<token>
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

func fullVersion() string {
	return fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

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
		Version: fullVersion(),
		Long: `bb-cli is a command-line interface for Bitbucket Cloud.

Set authentication environment variables before running:
  API token (recommended):
    BITBUCKET_USERNAME    your Bitbucket username
    BITBUCKET_TOKEN       your API token

  OAuth2 access token:
    BITBUCKET_TOKEN       your OAuth2 access token

  Internal API (e.g. project-branch-restrictions, /!api/internal/...):
    BITBUCKET_CSRF_TOKEN           value of the csrftoken cookie
    BITBUCKET_CLOUD_SESSION_TOKEN  value of the cloud.session.token cookie
  Internal endpoints reject Basic Auth — both cookies above are required.`,
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

	reposCmd := commands.NewReposCommand()
	commands.AddOutputFlag(reposCmd)

	workspacesCmd := commands.NewWorkspacesCommand()
	commands.AddOutputFlag(workspacesCmd)

	projectsCmd := commands.NewProjectsCommand()
	commands.AddOutputFlag(projectsCmd)

	pipelinesCmd := commands.NewPipelinesCommand()
	commands.AddOutputFlag(pipelinesCmd)

	issuesCmd := commands.NewIssuesCommand()
	commands.AddOutputFlag(issuesCmd)

	snippetsCmd := commands.NewSnippetsCommand()
	commands.AddOutputFlag(snippetsCmd)

	deploymentsCmd := commands.NewDeploymentsCommand()
	commands.AddOutputFlag(deploymentsCmd)

	branchRestrictionsCmd := commands.NewBranchRestrictionsCommand()
	commands.AddOutputFlag(branchRestrictionsCmd)

	projectBranchRestrictionsCmd := commands.NewProjectBranchRestrictionsCommand()
	commands.AddOutputFlag(projectBranchRestrictionsCmd)

	branchingModelCmd := commands.NewBranchingModelCommand()
	commands.AddOutputFlag(branchingModelCmd)

	commitStatusesCmd := commands.NewCommitStatusesCommand()
	commands.AddOutputFlag(commitStatusesCmd)

	downloadsCmd := commands.NewDownloadsCommand()
	commands.AddOutputFlag(downloadsCmd)

	usersCmd := commands.NewUsersCommand()
	commands.AddOutputFlag(usersCmd)

	propertiesCmd := commands.NewPropertiesCommand()
	commands.AddOutputFlag(propertiesCmd)

	addonCmd := commands.NewAddonCommand()
	commands.AddOutputFlag(addonCmd)

	rootCmd.AddCommand(
		prCmd, hooksCmd, searchCmd, refsCmd, commitsCmd, reportsCmd,
		reposCmd, workspacesCmd, projectsCmd, pipelinesCmd,
		issuesCmd, snippetsCmd, deploymentsCmd,
		branchRestrictionsCmd, projectBranchRestrictionsCmd,
		branchingModelCmd, commitStatusesCmd,
		downloadsCmd, usersCmd, propertiesCmd, addonCmd,
	)
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
