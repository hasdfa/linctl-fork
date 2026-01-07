package cmd

import (
	"testing"
)

func TestCommentCommandExists(t *testing.T) {
	if commentCmd == nil {
		t.Fatal("commentCmd should not be nil")
	}

	if commentCmd.Use != "comment" {
		t.Errorf("Expected Use 'comment', got '%s'", commentCmd.Use)
	}

	if commentCmd.Short != "Manage issue comments" {
		t.Errorf("Expected Short description mismatch, got '%s'", commentCmd.Short)
	}
}

func TestCommentListCommandExists(t *testing.T) {
	if commentListCmd == nil {
		t.Fatal("commentListCmd should not be nil")
	}

	if commentListCmd.Use != "list ISSUE-ID" {
		t.Errorf("Expected Use 'list ISSUE-ID', got '%s'", commentListCmd.Use)
	}

	// Check aliases
	if len(commentListCmd.Aliases) != 1 || commentListCmd.Aliases[0] != "ls" {
		t.Errorf("Expected alias 'ls', got %v", commentListCmd.Aliases)
	}
}

func TestCommentListRequiresOneArg(t *testing.T) {
	err := commentListCmd.Args(commentListCmd, []string{})
	if err == nil {
		t.Error("Expected error with 0 args")
	}

	err = commentListCmd.Args(commentListCmd, []string{"ENG-123"})
	if err != nil {
		t.Errorf("Expected no error with 1 arg, got: %v", err)
	}

	err = commentListCmd.Args(commentListCmd, []string{"ENG-123", "extra"})
	if err == nil {
		t.Error("Expected error with 2 args")
	}
}

func TestCommentCreateCommandExists(t *testing.T) {
	if commentCreateCmd == nil {
		t.Fatal("commentCreateCmd should not be nil")
	}

	if commentCreateCmd.Use != "create ISSUE-ID" {
		t.Errorf("Expected Use 'create ISSUE-ID', got '%s'", commentCreateCmd.Use)
	}

	// Check aliases
	expectedAliases := []string{"add", "new"}
	if len(commentCreateCmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(commentCreateCmd.Aliases))
	}
}

func TestCommentCreateBodyFlag(t *testing.T) {
	flag := commentCreateCmd.Flags().Lookup("body")
	if flag == nil {
		t.Fatal("commentCreateCmd should have --body flag")
	}
	if flag.Shorthand != "b" {
		t.Errorf("Expected shorthand 'b', got '%s'", flag.Shorthand)
	}
}

func TestCommentDeleteCommandExists(t *testing.T) {
	if commentDeleteCmd == nil {
		t.Fatal("commentDeleteCmd should not be nil")
	}

	if commentDeleteCmd.Use != "delete COMMENT-ID" {
		t.Errorf("Expected Use 'delete COMMENT-ID', got '%s'", commentDeleteCmd.Use)
	}

	if commentDeleteCmd.Short != "Delete a comment" {
		t.Errorf("Expected Short 'Delete a comment', got '%s'", commentDeleteCmd.Short)
	}

	// Check aliases
	if len(commentDeleteCmd.Aliases) != 1 || commentDeleteCmd.Aliases[0] != "rm" {
		t.Errorf("Expected alias 'rm', got %v", commentDeleteCmd.Aliases)
	}
}

func TestCommentDeleteRequiresOneArg(t *testing.T) {
	err := commentDeleteCmd.Args(commentDeleteCmd, []string{})
	if err == nil {
		t.Error("Expected error with 0 args")
	}

	err = commentDeleteCmd.Args(commentDeleteCmd, []string{"comment-123"})
	if err != nil {
		t.Errorf("Expected no error with 1 arg, got: %v", err)
	}

	err = commentDeleteCmd.Args(commentDeleteCmd, []string{"comment-123", "extra"})
	if err == nil {
		t.Error("Expected error with 2 args")
	}
}

func TestCommentListFlags(t *testing.T) {
	// Check limit flag
	limitFlag := commentListCmd.Flags().Lookup("limit")
	if limitFlag == nil {
		t.Fatal("commentListCmd should have --limit flag")
	}
	if limitFlag.Shorthand != "l" {
		t.Errorf("Expected shorthand 'l', got '%s'", limitFlag.Shorthand)
	}

	// Check sort flag
	sortFlag := commentListCmd.Flags().Lookup("sort")
	if sortFlag == nil {
		t.Fatal("commentListCmd should have --sort flag")
	}
	if sortFlag.Shorthand != "o" {
		t.Errorf("Expected shorthand 'o', got '%s'", sortFlag.Shorthand)
	}
}
