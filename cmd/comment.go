package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/pkg/auth"
	"github.com/dorkitude/linctl/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage issue comments",
	Long: `Manage comments on Linear issues including listing, creating, deleting, and resolving comments.

Examples:
  linctl comment list LIN-123                           # List comments for an issue
  linctl comment ls LIN-123                             # List comments (alias)
  linctl comment create LIN-123 --body "This is fixed"  # Add a comment
  linctl comment delete <comment-id>                    # Delete a comment
  linctl comment rm <comment-id>                        # Delete a comment (alias)
  linctl comment remove <comment-id>                    # Delete a comment (alias)
  linctl comment resolve <comment-id>                   # Resolve a comment thread
  linctl comment unresolve <comment-id>                 # Unresolve a comment thread`,
}

var commentListCmd = &cobra.Command{
	Use:     "list ISSUE-ID",
	Aliases: []string{"ls"},
	Short:   "List comments for an issue",
	Long:    `List all comments for a specific issue.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		issueID := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get limit
		limit, _ := cmd.Flags().GetInt("limit")

		// Get sort option
		sortBy, _ := cmd.Flags().GetString("sort")
		orderBy := ""
		if sortBy != "" {
			switch sortBy {
			case "created", "createdAt":
				orderBy = "createdAt"
			case "updated", "updatedAt":
				orderBy = "updatedAt"
			case "linear":
				// Use empty string for Linear's default sort
				orderBy = ""
			default:
				output.Error(fmt.Sprintf("Invalid sort option: %s. Valid options are: linear, created, updated", sortBy), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		// Get comments
		comments, err := client.GetIssueComments(context.Background(), issueID, limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list comments: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Filter out child comments if --no-children flag is set
		noChildren, _ := cmd.Flags().GetBool("no-children")
		if noChildren {
			var rootComments []api.Comment
			for _, comment := range comments.Nodes {
				if comment.Parent == nil || comment.Parent.ID == "" {
					rootComments = append(rootComments, comment)
				}
			}
			comments.Nodes = rootComments
		}

		// Handle output
		if jsonOut {
			output.JSON(comments.Nodes)
		} else if plaintext {
			for i, comment := range comments.Nodes {
				if i > 0 {
					fmt.Println("---")
				}
				fmt.Printf("Author: %s\n", commentAuthorName(&comment))
				fmt.Printf("Date: %s\n", comment.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("Comment:\n%s\n", comment.Body)
			}
		} else {
			// Rich display
			if len(comments.Nodes) == 0 {
				fmt.Printf("\n%s No comments on issue %s\n",
					color.New(color.FgYellow).Sprint("ℹ️"),
					color.New(color.FgCyan).Sprint(issueID))
				return
			}

			fmt.Printf("\n%s Comments on %s (%d)\n\n",
				color.New(color.FgCyan, color.Bold).Sprint("💬"),
				color.New(color.FgCyan).Sprint(issueID),
				len(comments.Nodes))

			for i, comment := range comments.Nodes {
				if i > 0 {
					fmt.Println(strings.Repeat("─", 50))
				}

				// Header with author and time
				timeAgo := formatTimeAgo(comment.CreatedAt)
				fmt.Printf("%s %s %s\n",
					color.New(color.FgCyan, color.Bold).Sprint(commentAuthorName(&comment)),
					color.New(color.FgWhite, color.Faint).Sprint("•"),
					color.New(color.FgWhite, color.Faint).Sprint(timeAgo))

				// Comment body
				fmt.Printf("\n%s\n\n", comment.Body)
			}
		}
	},
}

var commentCreateCmd = &cobra.Command{
	Use:     "create ISSUE-ID",
	Aliases: []string{"add", "new"},
	Short:   "Create a comment on an issue",
	Long: `Add a new comment to a specific issue.

Use the --parent flag to create a threaded reply under an existing comment.

Examples:
  linctl comment create LIN-123 --body "This is a top-level comment"
  linctl comment create LIN-123 --body "This is a reply" --parent COMMENT-ID`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		issueID := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get comment body
		body, _ := cmd.Flags().GetString("body")
		if body == "" {
			output.Error("Comment body is required (--body)", plaintext, jsonOut)
			os.Exit(1)
		}

		// Get optional parent ID for threaded replies
		parentID, _ := cmd.Flags().GetString("parent")

		// Create comment
		comment, err := client.CreateComment(context.Background(), issueID, body, parentID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(comment)
		} else if plaintext {
			fmt.Printf("Created comment on %s\n", issueID)
			fmt.Printf("ID: %s\n", comment.ID)
			fmt.Printf("Author: %s\n", comment.User.Name)
			fmt.Printf("Date: %s\n", comment.CreatedAt.Format("2006-01-02 15:04:05"))
			if comment.Parent != nil && comment.Parent.ID != "" {
				fmt.Printf("Parent: %s\n", comment.Parent.ID)
			}
		} else {
			if comment.Parent != nil && comment.Parent.ID != "" {
				fmt.Printf("%s Added reply to comment on %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					color.New(color.FgCyan, color.Bold).Sprint(issueID))
			} else {
				fmt.Printf("%s Added comment to %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					color.New(color.FgCyan, color.Bold).Sprint(issueID))
			}
			fmt.Printf("ID: %s\n", color.New(color.FgWhite, color.Faint).Sprint(comment.ID))
			fmt.Printf("\n%s\n", comment.Body)
		}
	},
}

var commentDeleteCmd = &cobra.Command{
	Use:     "delete COMMENT-ID",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a comment",
	Long:    `Delete a comment by its ID.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		commentID := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		err = client.DeleteComment(context.Background(), commentID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(map[string]interface{}{
				"status":    "success",
				"commentId": commentID,
				"message":   "Comment deleted successfully",
			})
		} else if plaintext {
			fmt.Printf("Deleted comment %s\n", commentID)
		} else {
			fmt.Printf("%s Deleted comment %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan).Sprint(commentID))
		}
	},
}

var commentResolveCmd = &cobra.Command{
	Use:   "resolve COMMENT-ID",
	Short: "Resolve a comment thread",
	Long:  `Mark a comment thread as resolved. The current user will be set as the resolver.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		commentID := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		comment, err := client.ResolveComment(context.Background(), commentID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to resolve comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(comment)
		} else if plaintext {
			fmt.Printf("Resolved comment %s\n", commentID)
			if comment.ResolvingUser != nil {
				fmt.Printf("Resolved by: %s\n", comment.ResolvingUser.Name)
			}
		} else {
			fmt.Printf("%s Resolved comment %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan).Sprint(commentID))
			if comment.ResolvingUser != nil {
				fmt.Printf("Resolved by: %s\n",
					color.New(color.FgWhite, color.Faint).Sprint(comment.ResolvingUser.Name))
			}
		}
	},
}

var commentUnresolveCmd = &cobra.Command{
	Use:   "unresolve COMMENT-ID",
	Short: "Unresolve a comment thread",
	Long:  `Remove the resolved status from a comment thread.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		commentID := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		comment, err := client.UnresolveComment(context.Background(), commentID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to unresolve comment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(comment)
		} else if plaintext {
			fmt.Printf("Unresolved comment %s\n", commentID)
		} else {
			fmt.Printf("%s Unresolved comment %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan).Sprint(commentID))
		}
	},
}

// formatTimeAgo formats a time as a human-readable "time ago" string
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if duration < 30*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if duration < 365*24*time.Hour {
		months := int(duration.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := int(duration.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

func init() {
	rootCmd.AddCommand(commentCmd)
	commentCmd.AddCommand(commentListCmd)
	commentCmd.AddCommand(commentCreateCmd)
	commentCmd.AddCommand(commentDeleteCmd)
	commentCmd.AddCommand(commentResolveCmd)
	commentCmd.AddCommand(commentUnresolveCmd)

	// List command flags
	commentListCmd.Flags().IntP("limit", "l", 50, "Maximum number of comments to return")
	commentListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
	commentListCmd.Flags().Bool("no-children", false, "Only show root comments (skip comments that have a parent)")

	// Create command flags
	commentCreateCmd.Flags().StringP("body", "b", "", "Comment body (required)")
	commentCreateCmd.Flags().String("parent", "", "Parent comment ID for threaded replies")
	_ = commentCreateCmd.MarkFlagRequired("body")
}
