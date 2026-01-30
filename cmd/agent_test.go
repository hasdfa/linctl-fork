package cmd

import (
	"testing"
)

func TestAgentCommandExists(t *testing.T) {
	if agentCmd == nil {
		t.Fatal("agentCmd should not be nil")
	}

	if agentCmd.Use != "agent [issue-id]" {
		t.Errorf("Expected Use 'agent [issue-id]', got '%s'", agentCmd.Use)
	}

	if agentCmd.Short != "View agent session for an issue" {
		t.Errorf("Expected Short description mismatch, got '%s'", agentCmd.Short)
	}
}

func TestAgentMentionCommandExists(t *testing.T) {
	if agentMentionCmd == nil {
		t.Fatal("agentMentionCmd should not be nil")
	}

	if agentMentionCmd.Use != "mention [issue-id] [message]" {
		t.Errorf("Expected Use 'mention [issue-id] [message]', got '%s'", agentMentionCmd.Use)
	}

	if agentMentionCmd.Short != "@mention an agent with a message" {
		t.Errorf("Expected Short '@mention an agent with a message', got '%s'", agentMentionCmd.Short)
	}
}

func TestAgentMentionRequiresTwoArgs(t *testing.T) {
	err := agentMentionCmd.Args(agentMentionCmd, []string{"ENG-80"})
	if err == nil {
		t.Error("Expected error with only 1 arg")
	}

	err = agentMentionCmd.Args(agentMentionCmd, []string{"ENG-80", "message"})
	if err != nil {
		t.Errorf("Expected no error with 2 args, got: %v", err)
	}

	err = agentMentionCmd.Args(agentMentionCmd, []string{"ENG-80", "message", "extra"})
	if err == nil {
		t.Error("Expected error with 3 args")
	}
}
