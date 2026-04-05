package commands

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

func TestAddOutputFlag(t *testing.T) {
	output.Format = ""
	cmd := &cobra.Command{Use: "test"}

	AddOutputFlag(cmd)

	if cmd.PersistentFlags().Lookup(outputFlagName) == nil {
		t.Fatal("expected output flag to be registered")
	}
	if output.Format != "table" {
		t.Fatalf("expected default output format to be table, got %q", output.Format)
	}

	if err := cmd.PersistentFlags().Parse([]string{"--output", "json"}); err != nil {
		t.Fatalf("parsing output flag: %v", err)
	}
	if output.Format != "json" {
		t.Fatalf("expected parsed output format json, got %q", output.Format)
	}
}
