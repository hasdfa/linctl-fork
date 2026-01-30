package cmd

import (
	"strings"
	"testing"
)

func TestIssueCreateCmd_HelpIncludesProjectFlag(t *testing.T) {
	usage := issueCreateCmd.UsageString()
	if !(strings.Contains(usage, "--project") &&
		strings.Contains(usage, "Project ID to assign issue to")) {
		t.Fatalf("create usage missing project flag/help text. got:\n%s", usage)
	}
}

func TestIssueUpdateCmd_HelpIncludesProjectFlag(t *testing.T) {
	usage := issueUpdateCmd.UsageString()
	if !(strings.Contains(usage, "--project") &&
		strings.Contains(usage, "Project ID to assign issue to") &&
		strings.Contains(usage, "unassigned")) {
		t.Fatalf("update usage missing project flag/help text. got:\n%s", usage)
	}
}
