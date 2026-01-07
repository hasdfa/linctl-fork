package cmd

import (
	"testing"
)

func TestIssueCommandExists(t *testing.T) {
	if issueCmd == nil {
		t.Fatal("issueCmd should not be nil")
	}

	if issueCmd.Use != "issue" {
		t.Errorf("Expected Use 'issue', got '%s'", issueCmd.Use)
	}
}

func TestIssueCreateCommandExists(t *testing.T) {
	if issueCreateCmd == nil {
		t.Fatal("issueCreateCmd should not be nil")
	}

	if issueCreateCmd.Use != "create" {
		t.Errorf("Expected Use 'create', got '%s'", issueCreateCmd.Use)
	}
}

func TestIssueUpdateCommandExists(t *testing.T) {
	if issueUpdateCmd == nil {
		t.Fatal("issueUpdateCmd should not be nil")
	}

	if issueUpdateCmd.Use != "update [issue-id]" {
		t.Errorf("Expected Use 'update [issue-id]', got '%s'", issueUpdateCmd.Use)
	}
}

func TestIssueDelegateFlagOnUpdate(t *testing.T) {
	flag := issueUpdateCmd.Flags().Lookup("delegate")
	if flag == nil {
		t.Fatal("issueUpdateCmd should have --delegate flag")
	}
	if flag.Usage != "Delegate to agent (email, name, displayName, or 'none' to remove)" {
		t.Errorf("Unexpected delegate flag usage: %s", flag.Usage)
	}
}

func TestIssueDelegateFlagOnCreate(t *testing.T) {
	flag := issueCreateCmd.Flags().Lookup("delegate")
	if flag == nil {
		t.Fatal("issueCreateCmd should have --delegate flag")
	}
}

func TestIssueLabelFlagOnCreate(t *testing.T) {
	flag := issueCreateCmd.Flags().Lookup("label")
	if flag == nil {
		t.Fatal("issueCreateCmd should have --label flag")
	}
}

func TestIssueCreateRequiredFlags(t *testing.T) {
	// Title flag should exist
	titleFlag := issueCreateCmd.Flags().Lookup("title")
	if titleFlag == nil {
		t.Fatal("issueCreateCmd should have --title flag")
	}

	// Team flag should exist
	teamFlag := issueCreateCmd.Flags().Lookup("team")
	if teamFlag == nil {
		t.Fatal("issueCreateCmd should have --team flag")
	}
	if teamFlag.Shorthand != "t" {
		t.Errorf("Expected shorthand 't' for team, got '%s'", teamFlag.Shorthand)
	}
}

func TestIssueUpdateFlags(t *testing.T) {
	// Title flag
	titleFlag := issueUpdateCmd.Flags().Lookup("title")
	if titleFlag == nil {
		t.Fatal("issueUpdateCmd should have --title flag")
	}

	// State flag
	stateFlag := issueUpdateCmd.Flags().Lookup("state")
	if stateFlag == nil {
		t.Fatal("issueUpdateCmd should have --state flag")
	}
	if stateFlag.Shorthand != "s" {
		t.Errorf("Expected shorthand 's' for state, got '%s'", stateFlag.Shorthand)
	}

	// Assignee flag
	assigneeFlag := issueUpdateCmd.Flags().Lookup("assignee")
	if assigneeFlag == nil {
		t.Fatal("issueUpdateCmd should have --assignee flag")
	}

	// Priority flag
	priorityFlag := issueUpdateCmd.Flags().Lookup("priority")
	if priorityFlag == nil {
		t.Fatal("issueUpdateCmd should have --priority flag")
	}
}
