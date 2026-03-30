package commands

import (
	"github.com/spf13/cobra"

	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

const outputFlagName = "output"

// AddOutputFlag registers the --output flag on the given command.
func AddOutputFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&output.Format, outputFlagName, "table",
		"Output format: table, markdown, json, id")
}
