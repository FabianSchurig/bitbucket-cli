package main

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestFullVersion(t *testing.T) {
	oldVersion, oldCommit, oldDate := version, commit, date
	version, commit, date = "v1.0.0", "abc123", "2026-04-04"
	t.Cleanup(func() {
		version, commit, date = oldVersion, oldCommit, oldDate
	})

	got := fullVersion()
	if !strings.Contains(got, "v1.0.0") || !strings.Contains(got, "abc123") || !strings.Contains(got, "2026-04-04") {
		t.Fatalf("unexpected fullVersion output: %q", got)
	}
}

func TestNewRootCmd(t *testing.T) {
	cmd := newRootCmd()
	if cmd.Use != "bb-cli" {
		t.Fatalf("unexpected command use %q", cmd.Use)
	}
	if cmd.Version == "" {
		t.Fatal("expected version on root command")
	}
	if len(cmd.Commands()) != 20 {
		t.Fatalf("expected 20 subcommands, got %d", len(cmd.Commands()))
	}

	for _, sub := range cmd.Commands() {
		if sub.PersistentFlags().Lookup("output") == nil {
			t.Fatalf("expected --output flag on %s", sub.Name())
		}
	}
}

func TestSetColoredHelp(t *testing.T) {
	cmd := &cobra.Command{Use: "demo"}
	setColoredHelp(cmd)

	usage := cmd.UsageTemplate()
	if !strings.Contains(usage, "{{bold \"Usage:\"}}") || !strings.Contains(usage, "{{yellow .UseLine}}") {
		t.Fatalf("expected colored help template, got %q", usage)
	}
}
