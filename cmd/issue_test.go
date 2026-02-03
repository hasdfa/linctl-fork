package cmd

import (
	"errors"
	"testing"
)

func TestIsValidUUID(t *testing.T) {
	valid := []string{
		"123e4567-e89b-12d3-a456-426614174000",
		"00000000-0000-0000-0000-000000000000",
		"ABCDEFAB-CDEF-ABCD-EFAB-CDEFABCDEFAB",
	}
	invalid := []string{
		"", "unassigned", "1234", "g23e4567-e89b-12d3-a456-426614174000",
		"123e4567e89b12d3a456426614174000",
		"123e4567-e89b-12d3-a456-426614174000-extra",
	}

	for _, v := range valid {
		if !isValidUUID(v) {
			t.Errorf("expected valid UUID: %s", v)
		}
	}
	for _, v := range invalid {
		if isValidUUID(v) {
			t.Errorf("expected invalid UUID: %s", v)
		}
	}
}

func TestBuildProjectInput(t *testing.T) {
	// Empty → ok=false, no input
	if val, ok, err := buildProjectInput(""); err != nil || ok || val != nil {
		t.Errorf("empty flag: want (nil,false,nil), got (%v,%v,%v)", val, ok, err)
	}

	// unassigned → ok=true, val=nil
	if val, ok, err := buildProjectInput("unassigned"); err != nil || !ok || val != nil {
		t.Errorf("unassigned: want (nil,true,nil), got (%v,%v,%v)", val, ok, err)
	}

	// valid uuid → ok=true, val=uuid
	uuid := "123e4567-e89b-12d3-a456-426614174000"
	if val, ok, err := buildProjectInput(uuid); err != nil || !ok || val != uuid {
		t.Errorf("uuid: want (%s,true,nil), got (%v,%v,%v)", uuid, val, ok, err)
	}

	// invalid uuid → error
	if _, _, err := buildProjectInput("not-a-uuid"); err == nil {
		t.Errorf("expected error for invalid uuid")
	}
}

func TestIsProjectNotFoundErr(t *testing.T) {
	cases := []struct {
		in   error
		want bool
	}{
		{errors.New("GraphQL errors: [{ message: 'Project not found' }]"), true},
		{errors.New("something about projectId not found"), true},
		{errors.New("issue not found"), false},
		{errors.New("unknown error"), false},
		{nil, false},
	}
	for _, c := range cases {
		got := isProjectNotFoundErr(c.in)
		if got != c.want {
			t.Errorf("isProjectNotFoundErr(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

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
