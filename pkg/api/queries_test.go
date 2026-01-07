package api

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAgentSessionUnmarshal(t *testing.T) {
	jsonData := `{
		"id": "test-123",
		"status": "active",
		"createdAt": "2026-01-06T20:17:02.878Z",
		"updatedAt": "2026-01-06T20:21:28.402Z",
		"appUser": {
			"id": "user-456",
			"name": "Test Agent",
			"displayName": "testagent"
		},
		"activities": {
			"nodes": [
				{
					"id": "activity-1",
					"createdAt": "2026-01-06T20:17:05.000Z",
					"ephemeral": false,
					"content": {"type": "thought", "body": "Thinking..."}
				}
			],
			"pageInfo": {"hasNextPage": false, "endCursor": ""}
		}
	}`

	var session AgentSession
	err := json.Unmarshal([]byte(jsonData), &session)
	if err != nil {
		t.Fatalf("Failed to unmarshal AgentSession: %v", err)
	}

	if session.ID != "test-123" {
		t.Errorf("Expected ID 'test-123', got '%s'", session.ID)
	}
	if session.Status != "active" {
		t.Errorf("Expected Status 'active', got '%s'", session.Status)
	}
	if session.AppUser == nil {
		t.Fatal("Expected AppUser to be non-nil")
	}
	if session.AppUser.DisplayName != "testagent" {
		t.Errorf("Expected DisplayName 'testagent', got '%s'", session.AppUser.DisplayName)
	}
	if session.Activities == nil || len(session.Activities.Nodes) != 1 {
		t.Fatal("Expected 1 activity")
	}
	if session.Activities.Nodes[0].Content["type"] != "thought" {
		t.Errorf("Expected activity type 'thought', got '%v'", session.Activities.Nodes[0].Content["type"])
	}
}

func TestAgentActivityContent(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantType string
		wantBody string
	}{
		{
			name:     "thought content",
			json:     `{"type": "thought", "body": "Processing request"}`,
			wantType: "thought",
			wantBody: "Processing request",
		},
		{
			name:     "response content",
			json:     `{"type": "response", "body": "Task completed"}`,
			wantType: "response",
			wantBody: "Task completed",
		},
		{
			name:     "action content",
			json:     `{"type": "action", "action": "Bash", "parameter": "git status"}`,
			wantType: "action",
			wantBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var content map[string]interface{}
			err := json.Unmarshal([]byte(tt.json), &content)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			if content["type"] != tt.wantType {
				t.Errorf("Expected type '%s', got '%v'", tt.wantType, content["type"])
			}
			if body, ok := content["body"].(string); ok && body != tt.wantBody {
				t.Errorf("Expected body '%s', got '%s'", tt.wantBody, body)
			}
		})
	}
}

func TestIssueDelegateField(t *testing.T) {
	jsonData := `{
		"id": "issue-123",
		"identifier": "ENG-80",
		"title": "Test Issue",
		"createdAt": "2026-01-06T10:00:00.000Z",
		"updatedAt": "2026-01-06T10:00:00.000Z",
		"delegate": {
			"id": "user-789",
			"name": "Test Agent",
			"displayName": "testagent",
			"email": "agent@example.com"
		}
	}`

	var issue Issue
	err := json.Unmarshal([]byte(jsonData), &issue)
	if err != nil {
		t.Fatalf("Failed to unmarshal Issue: %v", err)
	}

	if issue.Delegate == nil {
		t.Fatal("Expected Delegate to be non-nil")
	}
	if issue.Delegate.DisplayName != "testagent" {
		t.Errorf("Expected Delegate DisplayName 'testagent', got '%s'", issue.Delegate.DisplayName)
	}
}

func TestCommentAgentSessionField(t *testing.T) {
	jsonData := `{
		"id": "comment-123",
		"body": "Test comment",
		"createdAt": "2026-01-06T10:00:00.000Z",
		"updatedAt": "2026-01-06T10:00:00.000Z",
		"agentSession": {
			"id": "session-456",
			"status": "complete"
		}
	}`

	var comment Comment
	err := json.Unmarshal([]byte(jsonData), &comment)
	if err != nil {
		t.Fatalf("Failed to unmarshal Comment: %v", err)
	}

	if comment.AgentSession == nil {
		t.Fatal("Expected AgentSession to be non-nil")
	}
	if comment.AgentSession.Status != "complete" {
		t.Errorf("Expected AgentSession Status 'complete', got '%s'", comment.AgentSession.Status)
	}
}

func TestAgentSessionTimeFields(t *testing.T) {
	jsonData := `{
		"id": "test-123",
		"status": "active",
		"createdAt": "2026-01-06T20:17:02.878Z",
		"updatedAt": "2026-01-06T20:21:28.402Z"
	}`

	var session AgentSession
	err := json.Unmarshal([]byte(jsonData), &session)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	expectedCreated := time.Date(2026, 1, 6, 20, 17, 2, 878000000, time.UTC)
	if !session.CreatedAt.Equal(expectedCreated) {
		t.Errorf("CreatedAt mismatch: got %v, want %v", session.CreatedAt, expectedCreated)
	}
}

func TestOrganizationUnmarshal(t *testing.T) {
	jsonData := `{
		"id": "org-123",
		"name": "Acme Corp",
		"urlKey": "acme"
	}`

	var org Organization
	err := json.Unmarshal([]byte(jsonData), &org)
	if err != nil {
		t.Fatalf("Failed to unmarshal Organization: %v", err)
	}

	if org.ID != "org-123" {
		t.Errorf("Expected ID 'org-123', got '%s'", org.ID)
	}
	if org.Name != "Acme Corp" {
		t.Errorf("Expected Name 'Acme Corp', got '%s'", org.Name)
	}
	if org.URLKey != "acme" {
		t.Errorf("Expected URLKey 'acme', got '%s'", org.URLKey)
	}
}

func TestLabelUnmarshal(t *testing.T) {
	jsonData := `{
		"id": "label-123",
		"name": "backend",
		"color": "#ff0000"
	}`

	var label Label
	err := json.Unmarshal([]byte(jsonData), &label)
	if err != nil {
		t.Fatalf("Failed to unmarshal Label: %v", err)
	}

	if label.ID != "label-123" {
		t.Errorf("Expected ID 'label-123', got '%s'", label.ID)
	}
	if label.Name != "backend" {
		t.Errorf("Expected Name 'backend', got '%s'", label.Name)
	}
}
