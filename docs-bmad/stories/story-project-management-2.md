# Story 1.2: Project Creation & Archival

**Status:** ready-for-dev

---

## User Story

As a project manager or team lead using linctl,
I want to create new projects and archive completed ones from the CLI,
So that I can manage project lifecycle without leaving my terminal workflow.

---

## Acceptance Criteria

**AC #1:** Given a project name and team key, when I run `linctl project create --name "Q1 Backend" --team ENG`, then a new project is created in Linear with default values (state: planned, priority: 0).

**AC #2:** Given project creation with optional fields, when I run `linctl project create --name "Test" --team ENG --state started --priority 1 --description "Test project"`, then the project is created with all specified field values.

**AC #3:** Given a valid project UUID, when I run `linctl project archive PROJECT-UUID`, then the project is archived in Linear and a success message is displayed.

**AC #4:** Given missing required fields, when I attempt project creation without --name or --team, then the command fails with error: "Both --name and --team are required".

**AC #5:** Given an invalid team key, when I attempt project creation, then the command fails with error: "Team 'INVALID' not found. Use 'linctl team list' to see available teams."

**AC #6:** Given invalid state or priority values, when I attempt project creation, then validation fails before API call with clear error messages.

**AC #7:** All commands support --json, --plaintext output formats and follow existing linctl formatting conventions.

---

## Implementation Details

### Tasks / Subtasks

- [ ] **Task 1:** Implement CreateProject() method in pkg/api/queries.go (AC: #1, #2)
  - [ ] Define GraphQL projectCreate mutation string
  - [ ] Implement function signature: `CreateProject(ctx context.Context, input map[string]interface{}) (*Project, error)`
  - [ ] Marshal input to JSON for GraphQL query
  - [ ] Call `c.query(ctx, mutation)` and handle response
  - [ ] Parse response and extract project data
  - [ ] Return structured Project type with error handling

- [ ] **Task 2:** Implement ArchiveProject() method in pkg/api/queries.go (AC: #3)
  - [ ] Define GraphQL projectArchive mutation string
  - [ ] Implement function signature: `ArchiveProject(ctx context.Context, id string) (bool, error)`
  - [ ] Call API with project UUID
  - [ ] Return success boolean and error

- [ ] **Task 3:** Create projectCreateCmd in cmd/project.go (AC: #1, #2, #4, #5, #6, #7)
  - [ ] Define `projectCreateCmd` variable following Cobra pattern
  - [ ] Add required flags: --name, --team
  - [ ] Add optional flags: --description, --state, --priority, --target-date
  - [ ] Implement Run function with validation
  - [ ] Validate required fields (name, team)
  - [ ] Validate optional fields (state, priority ranges)
  - [ ] Get auth header and create API client
  - [ ] Resolve team key to team UUID using GetTeam()
  - [ ] Build input map from flags
  - [ ] Call client.CreateProject()
  - [ ] Format and display output (table/JSON/plaintext)
  - [ ] Add comprehensive help text with examples

- [ ] **Task 4:** Create projectArchiveCmd in cmd/project.go (AC: #3, #7)
  - [ ] Define `projectArchiveCmd` variable
  - [ ] Accept project UUID as required argument
  - [ ] Validate argument provided
  - [ ] Get auth header and create API client
  - [ ] Call client.ArchiveProject()
  - [ ] Display success message with project name
  - [ ] Handle errors with clear messages

- [ ] **Task 5:** Register new commands in init() function (AC: all)
  - [ ] Add `projectCmd.AddCommand(projectCreateCmd)` to init()
  - [ ] Add `projectCmd.AddCommand(projectArchiveCmd)` to init()
  - [ ] Ensure commands appear in help text

- [ ] **Task 6:** Test all acceptance criteria (AC: #1-#7)
  - [ ] Manual test: Create project with required fields only
  - [ ] Manual test: Create project with all optional fields
  - [ ] Manual test: Archive project successfully
  - [ ] Test error: Missing --name flag
  - [ ] Test error: Missing --team flag
  - [ ] Test error: Invalid team key
  - [ ] Test error: Invalid state value
  - [ ] Test error: Invalid priority value (< 0 or > 4)
  - [ ] Test all output formats: table, --json, --plaintext

### Technical Summary

This story implements complete project creation and archival capabilities, enabling users to manage the full project lifecycle from linctl. The implementation adds two new commands to the project command group and extends the API client with GraphQL mutation support.

**Key Implementation Points:**
- Create new Cobra commands following existing project.go patterns
- Implement GraphQL mutations in API client (projectCreate, projectArchive)
- Validate all inputs before making API calls
- Resolve team key to UUID using existing GetTeam() method
- Support optional field configuration at creation time
- Maintain consistent error handling and output formatting

**GraphQL Mutations:**
- `projectCreate`: Accept name (required), teamId (required), plus optional fields (description, state, priority, targetDate, color)
- `projectArchive`: Accept project UUID, return success status
- Return complete project data including team and timestamps

**Validation Rules:**
- Name: Required, 1-255 characters
- Team: Required, must exist in workspace
- State: Optional, must be one of: planned, started, paused, completed, canceled
- Priority: Optional, must be 0-4 (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)
- Description: Optional, any string

### Project Structure Notes

- **Files to modify:**
  - `pkg/api/queries.go` (add CreateProject and ArchiveProject methods)
  - `pkg/api/client.go` (add Project type if missing)
  - `cmd/project.go` (add projectCreateCmd and projectArchiveCmd, update init())

- **Expected test locations:**
  - Manual testing procedures in `tests/manual_project_tests.sh`
  - No smoke tests added (write commands have side effects)

- **Estimated effort:** 3 story points (3.5 hours)

- **Prerequisites:** Story 1.1 complete (enables full issue-project workflow)

### Key Code References

**Existing Patterns to Follow:**

1. **Command Definition Pattern** (from cmd/project.go:50-100):
   ```go
   var projectCreateCmd = &cobra.Command{
       Use:   "create",
       Short: "Create a new project",
       Long:  `Create a new project in Linear workspace with required and optional configuration.`,
       Run: func(cmd *cobra.Command, args []string) {
           // Implementation here
       },
   }

   func init() {
       projectCreateCmd.Flags().String("name", "", "Project name (required)")
       projectCreateCmd.Flags().String("team", "", "Team key (required)")
       projectCreateCmd.Flags().String("state", "", "Project state (planned|started|paused|completed|canceled)")
       projectCreateCmd.Flags().Int("priority", 0, "Priority (0-4: None, Urgent, High, Normal, Low)")
       // ... more flags
   }
   ```

2. **Team Resolution Pattern** (from cmd/project.go:150-180):
   ```go
   teamKey, _ := cmd.Flags().GetString("team")
   team, err := client.GetTeam(context.Background(), teamKey)
   if err != nil {
       output.Error(fmt.Sprintf("Failed to find team '%s': %v", teamKey, err), plaintext, jsonOut)
       os.Exit(1)
   }
   input["teamId"] = team.ID
   ```

3. **Validation Pattern** (from cmd/issue.go):
   ```go
   // Required field validation
   if name == "" || teamKey == "" {
       output.Error("Both --name and --team are required", plaintext, jsonOut)
       os.Exit(1)
   }

   // State validation
   allowedStates := []string{"planned", "started", "paused", "completed", "canceled"}
   if state != "" && !contains(allowedStates, state) {
       output.Error(fmt.Sprintf("Invalid state. Must be one of: %v", allowedStates), plaintext, jsonOut)
       os.Exit(1)
   }

   // Priority validation
   if cmd.Flags().Changed("priority") {
       priority, _ := cmd.Flags().GetInt("priority")
       if priority < 0 || priority > 4 {
           output.Error("Priority must be between 0 (None) and 4 (Low)", plaintext, jsonOut)
           os.Exit(1)
       }
   }
   ```

4. **GraphQL Mutation Pattern** (from pkg/api/queries.go):
   ```go
   func (c *Client) CreateProject(ctx context.Context, input map[string]interface{}) (*Project, error) {
       inputJSON, _ := json.Marshal(input)

       query := fmt.Sprintf(`
           mutation {
               projectCreate(input: %s) {
                   success
                   project {
                       id name state priority url
                       team { id key name }
                       createdAt updatedAt
                   }
               }
           }
       `, string(inputJSON))

       result, err := c.query(ctx, query)
       if err != nil {
           return nil, err
       }

       // Parse and return Project
   }
   ```

**Relevant Code Locations:**
- `cmd/project.go:50-150` - Existing project commands (list, get) for pattern reference
- `cmd/project.go:18-30` - Helper functions (constructProjectURL)
- `pkg/api/queries.go` - API client methods and GraphQL patterns
- `pkg/api/client.go` - Client struct and type definitions

---

## Context References

**Tech-Spec:** [tech-spec.md](../tech-spec.md) - Primary context document containing:

- **Section 2.2 "Project Creation Implementation"** - Complete code examples
- **Section 2.4 "GraphQL Mutation Implementation"** - Mutation builder patterns
- **Section 4.2 "Project Creation (GraphQL)"** - Complete mutation schema
- **Section 4.4 "Project Archive (GraphQL)"** - Archive mutation details
- **Section 5.2 "Data Validation Rules"** - All validation requirements
- **Section 6.1 "Files to Modify"** - Complete file modification list
- **Section 7.1-7.4 "Existing Patterns to Follow"** - All patterns to replicate
- **Section 9.3 "Story 2 Implementation Steps"** - Step-by-step guide
- **Section 9.2 "Testing Strategy"** - Story 2 test cases

**Architecture:** See tech-spec.md sections:
- "Existing Codebase Structure" - File organization and patterns
- "Integration Points" - Linear GraphQL API specifications
- "Technical Approach" - Implementation strategy

---

## Dev Agent Record

### Context Reference

- [Story Context XML](./1-2-project-creation-archival.context.xml) - Generated 2025-11-06

### Agent Model Used

<!-- Will be populated during dev-story execution -->

### Debug Log References

<!-- Will be populated during dev-story execution -->

### Completion Notes

<!-- Will be populated during dev-story execution -->

### Files Modified

<!-- Will be populated during dev-story execution -->

### Test Results

<!-- Will be populated during dev-story execution -->

---

## Review Notes

<!-- Will be populated during code review -->
