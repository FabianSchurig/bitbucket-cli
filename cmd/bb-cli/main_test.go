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

// TestDeprecatedCommandsAreMarked verifies that deprecation survives code
// generation, without pinning the assertion to one live endpoint. Naming a
// specific deprecated endpoint made this test fail on every schema sync where
// Atlassian retired that endpoint (see #130); the generator-level guarantee is
// covered by scripts/internal/spec.TestDeprecatedOperation, so here we only
// assert that the marking is applied consistently across the real command tree.
func TestDeprecatedCommandsAreMarked(t *testing.T) {
	cmd := newRootCmd()

	var deprecated int
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			if sub.Deprecated != "" {
				deprecated++
			}
			walk(sub)
		}
	}
	walk(cmd)

	// Bitbucket always publishes some deprecated endpoints, and operations the
	// spec drops are retained as deprecated by partition_spec.py, so losing all
	// of them signals a broken generator rather than upstream drift.
	if deprecated == 0 {
		t.Fatal("expected at least one command to be marked deprecated")
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
