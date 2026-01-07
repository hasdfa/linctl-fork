package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/pkg/auth"
	"github.com/dorkitude/linctl/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentCmd = &cobra.Command{
	Use:   "agent [issue-id]",
	Short: "View agent session for an issue",
	Long: `View the agent session status and activity stream for an issue.

Examples:
  linctl agent ENG-80           # View agent session
  linctl agent ENG-80 --json    # Output as JSON`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		issue, err := client.GetIssueAgentSession(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Find agent session from comments
		var session *api.AgentSession
		if issue.Comments != nil {
			for _, comment := range issue.Comments.Nodes {
				if comment.AgentSession != nil {
					session = comment.AgentSession
					break
				}
			}
		}

		// Check if there's a delegate but no session yet
		if session == nil && issue.Delegate == nil {
			output.Info(fmt.Sprintf("No agent session found for %s", issue.Identifier), plaintext, jsonOut)
			return
		}

		if jsonOut {
			result := map[string]interface{}{
				"issue": issue.Identifier,
				"title": issue.Title,
			}
			if issue.Delegate != nil {
				result["delegate"] = issue.Delegate
			}
			if session != nil {
				result["agentSession"] = session
			}
			output.JSON(result)
			return
		}

		if plaintext {
			fmt.Printf("# Agent Session for %s\n\n", issue.Identifier)
			fmt.Printf("**Title**: %s\n", issue.Title)
			if issue.Delegate != nil {
				fmt.Printf("**Delegate**: %s (%s)\n", issue.Delegate.Name, issue.Delegate.DisplayName)
			}
			if session != nil {
				fmt.Printf("**Status**: %s\n", session.Status)
				if session.AppUser != nil {
					fmt.Printf("**Agent**: %s (%s)\n", session.AppUser.Name, session.AppUser.DisplayName)
				}
				fmt.Printf("**Started**: %s\n", session.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("**Updated**: %s\n", session.UpdatedAt.Format("2006-01-02 15:04:05"))

				if session.Activities != nil && len(session.Activities.Nodes) > 0 {
					fmt.Printf("\n## Activity Stream\n\n")
					for _, activity := range session.Activities.Nodes {
						activityType := "unknown"
						body := ""
						if t, ok := activity.Content["type"].(string); ok {
							activityType = t
						}
						if b, ok := activity.Content["body"].(string); ok {
							body = b
						} else if action, ok := activity.Content["action"].(string); ok {
							param, _ := activity.Content["parameter"].(string)
							body = fmt.Sprintf("%s: %s", action, param)
						}
						fmt.Printf("### [%s] %s\n", activityType, activity.CreatedAt.Format("15:04:05"))
						if body != "" {
							fmt.Printf("%s\n\n", body)
						}
					}
				}
			} else {
				fmt.Printf("**Status**: delegated (no session yet)\n")
			}
			return
		}

		// Rich display
		fmt.Printf("%s %s\n",
			color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier),
			color.New(color.FgWhite, color.Bold).Sprint(issue.Title))

		// Delegate info
		if issue.Delegate != nil {
			fmt.Printf("\n%s %s\n",
				color.New(color.FgYellow).Sprint("Delegate:"),
				color.New(color.FgCyan).Sprint(issue.Delegate.DisplayName))
		}

		if session == nil {
			fmt.Printf("\n%s\n", color.New(color.FgWhite, color.Faint).Sprint("Delegated but no session started yet"))
			return
		}

		// Status with color
		statusColor := color.New(color.FgWhite)
		switch session.Status {
		case "active":
			statusColor = color.New(color.FgGreen)
		case "complete":
			statusColor = color.New(color.FgBlue)
		case "awaitingInput":
			statusColor = color.New(color.FgYellow)
		case "error":
			statusColor = color.New(color.FgRed)
		case "pending":
			statusColor = color.New(color.FgMagenta)
		}

		fmt.Printf("%s %s\n",
			color.New(color.FgYellow).Sprint("Status:"),
			statusColor.Sprint(session.Status))

		if session.AppUser != nil {
			fmt.Printf("%s %s\n",
				color.New(color.FgYellow).Sprint("Agent:"),
				color.New(color.FgCyan).Sprint(session.AppUser.DisplayName))
		}

		fmt.Printf("%s %s\n",
			color.New(color.FgYellow).Sprint("Started:"),
			session.CreatedAt.Format("2006-01-02 15:04:05"))

		// Activity stream
		if session.Activities != nil && len(session.Activities.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.FgYellow, color.Bold).Sprint("Activity Stream:"))

			for _, activity := range session.Activities.Nodes {
				activityType := "unknown"
				body := ""
				if t, ok := activity.Content["type"].(string); ok {
					activityType = t
				}
				// Get body or action+parameter depending on type
				if b, ok := activity.Content["body"].(string); ok {
					body = b
				} else if action, ok := activity.Content["action"].(string); ok {
					param, _ := activity.Content["parameter"].(string)
					body = fmt.Sprintf("%s: %s", action, param)
				}

				// Color by type
				typeColor := color.New(color.FgWhite)
				switch activityType {
				case "thought":
					typeColor = color.New(color.FgMagenta)
				case "response":
					typeColor = color.New(color.FgGreen)
				case "action":
					typeColor = color.New(color.FgBlue)
				case "error":
					typeColor = color.New(color.FgRed)
				}

				timestamp := color.New(color.FgWhite, color.Faint).Sprint(activity.CreatedAt.Format("15:04:05"))
				fmt.Printf("\n  %s [%s]\n", timestamp, typeColor.Sprint(activityType))

				if body != "" {
					// Indent body text
					lines := strings.Split(body, "\n")
					for _, line := range lines {
						if len(line) > 80 {
							line = line[:77] + "..."
						}
						fmt.Printf("    %s\n", line)
					}
				}
			}

			if session.Activities.PageInfo.HasNextPage {
				fmt.Printf("\n%s More activities available\n",
					color.New(color.FgYellow).Sprint("ℹ️"))
			}
		} else {
			fmt.Printf("\n%s\n", color.New(color.FgWhite, color.Faint).Sprint("No activities yet"))
		}
	},
}

var agentMentionCmd = &cobra.Command{
	Use:   "mention [issue-id] [message]",
	Short: "@mention an agent with a message",
	Long: `@mention an agent on an issue to trigger them with a message.

Examples:
  linctl agent mention ENG-80 "Fix this bug"
  linctl agent mention ENG-80 "Please update the authentication flow to use JWT tokens"`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linctl auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		issueID := args[0]
		message := args[1]

		// Get the issue to find the delegated agent
		issue, err := client.GetIssueAgentSession(context.Background(), issueID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Find agent display name from delegate or existing session
		var agentDisplayName string
		if issue.Delegate != nil && issue.Delegate.DisplayName != "" {
			agentDisplayName = issue.Delegate.DisplayName
		} else if issue.Comments != nil {
			for _, comment := range issue.Comments.Nodes {
				if comment.AgentSession != nil && comment.AgentSession.AppUser != nil {
					agentDisplayName = comment.AgentSession.AppUser.DisplayName
					break
				}
			}
		}

		if agentDisplayName == "" {
			output.Error(fmt.Sprintf("No agent found for %s", issueID), plaintext, jsonOut)
			os.Exit(1)
		}

		// @mention the agent to trigger them
		commentID, err := client.MentionAgent(context.Background(), issue.ID, agentDisplayName, message)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to mention agent: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{
				"success":   true,
				"commentId": commentID,
				"issue":     issueID,
				"agent":     agentDisplayName,
				"message":   message,
			})
			return
		}

		fmt.Printf("✓ @%s mentioned on %s\n", agentDisplayName, issue.Identifier)
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentMentionCmd)
}
