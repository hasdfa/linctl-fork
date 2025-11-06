# Story 1.3: Project Updates & Enhanced Display

**Status:** ready-for-dev

---

## User Story

As a project manager using linctl for project tracking,
I want to update project fields and see complete project information in list/detail views,
So that I can manage project state, priority, initiatives, labels, descriptions, and summary entirely from the CLI.

---

## Acceptance Criteria

**AC #1:** Given a project UUID and updated field, when I run `linctl project update PROJECT-UUID --name "New Name"`, then the project name is updated in Linear.

**AC #2:** Given a project UUID, when I run `linctl project update PROJECT-UUID --state started`, then the project state is updated and validated against allowed values (planned, started, paused, completed, canceled).

**AC #3:** Given multiple field updates, when I run `linctl project update PROJECT-UUID --state started --priority 1 --label "urgent,backend"`, then all specified fields are updated in a single API call.

**AC #3.1:** Given a project UUID, when I run `linctl project update PROJECT-UUID --description "Full project description"`, then the project description is updated.

**AC #3.2:** Given a project UUID, when I run `linctl project update PROJECT-UUID --summary "Short summary text"`, then the project shortSummary field is updated.

**AC #4:** Given no field flags provided, when I run `linctl project update PROJECT-UUID`, then the command fails with error: "At least one field to update is required".

**AC #5:** Given invalid state value, when I attempt update, then validation fails with error: "Invalid state. Must be one of: planned, started, paused, completed, canceled".

**AC #6:** Given invalid priority value, when I attempt update, then validation fails with error: "Priority must be between 0 and 4".

**AC #7:** Given a project UUID, when I run `linctl project get PROJECT-UUID`, then the output displays all fields including state, priority, initiative, and labels.

**AC #8:** When I run `linctl project list`, then the table output includes State and Priority columns for each project.

**AC #9:** All commands support --json and --plaintext output formats with complete field data.

---

## Implementation Details

### Tasks / Subtasks

- [ ] **Task 1:** Implement UpdateProject() method in pkg/api/queries.go (AC: #1, #2, #3)
  - [ ] Define GraphQL projectUpdate mutation string
  - [ ] Implement function signature: `UpdateProject(ctx context.Context, id string, input map[string]interface{}) (*Project, error)`
  - [ ] Support partial updates (only provided fields in input map)
  - [ ] Marshal input to JSON for GraphQL query
  - [ ] Call `c.query(ctx, mutation)` with project UUID and input
  - [ ] Parse response and return updated Project
  - [ ] Handle errors appropriately

- [ ] **Task 2:** Create projectUpdateCmd in cmd/project.go (AC: #1-#6, #9)
  - [ ] Define `projectUpdateCmd` variable following Cobra pattern
  - [ ] Accept project UUID as required positional argument
  - [ ] Add optional flags: --name, --state, --priority, --initiative, --label, --description, --summary
  - [ ] Implement Run function with multi-field support
  - [ ] Validate project UUID argument provided
  - [ ] Build input map with ONLY changed flags (use `cmd.Flags().Changed()`)
  - [ ] Validate at least one field provided for update
  - [ ] Validate state against allowed values
  - [ ] Validate priority range (0-4)
  - [ ] Get auth header and create API client
  - [ ] Call client.UpdateProject() with project UUID and input
  - [ ] Format and display output (show updated fields)
  - [ ] Add comprehensive help text with examples

- [ ] **Task 3:** Enhance projectGetCmd GraphQL query (AC: #7, #9)
  - [ ] Locate existing projectGetCmd in cmd/project.go
  - [ ] Update GraphQL query to include additional fields:
    - state
    - priority
    - initiative { id name }
    - labels (array of strings)
  - [ ] Update Project type in pkg/api/client.go if needed
  - [ ] Ensure output formatter displays new fields

- [ ] **Task 4:** Enhance projectListCmd GraphQL query and output (AC: #8, #9)
  - [ ] Locate existing projectListCmd in cmd/project.go
  - [ ] Update GraphQL query to include state and priority fields
  - [ ] Modify table output to add "State" column
  - [ ] Modify table output to add "Priority" column
  - [ ] Keep table width reasonable (may need to abbreviate state values)
  - [ ] Ensure JSON/plaintext outputs include new fields

- [ ] **Task 5:** Update output formatters if needed (AC: #7, #8, #9)
  - [ ] Check pkg/output/output.go for project formatters
  - [ ] Update ProjectDetails formatter to show state, priority, initiative, labels
  - [ ] Update ProjectList formatter to include State and Priority columns
  - [ ] Maintain backward compatibility with existing output

- [ ] **Task 6:** Register projectUpdateCmd in init() function (AC: all)
  - [ ] Add `projectCmd.AddCommand(projectUpdateCmd)` to init()
  - [ ] Ensure command appears in help text

- [ ] **Task 7:** Test all acceptance criteria (AC: #1-#9)
  - [ ] Manual test: Update single field (name)
  - [ ] Manual test: Update single field (state)
  - [ ] Manual test: Update single field (priority)
  - [ ] Manual test: Multi-field update (state + priority)
  - [ ] Manual test: Update with initiative ID
  - [ ] Manual test: Update with labels (comma-separated)
  - [ ] Manual test: Update description field
  - [ ] Manual test: Update summary (shortSummary) field
  - [ ] Manual test: Multi-field update including description and summary
  - [ ] Test error: No fields provided
  - [ ] Test error: Invalid state value
  - [ ] Test error: Invalid priority value (< 0 or > 4)
  - [ ] Test error: Project not found (invalid UUID)
  - [ ] Verify enhanced `project get` shows all new fields
  - [ ] Verify enhanced `project list` includes state and priority
  - [ ] Test all output formats: table, --json, --plaintext

### Technical Summary

This story completes the project management feature by enabling field updates and enhancing display capabilities. The implementation adds a multi-field update command and extends existing list/get commands to show comprehensive project information including state, priority, initiatives, and labels.

**Key Implementation Points:**
- Create projectUpdateCmd with multiple optional field flags
- Use `cmd.Flags().Changed()` to detect explicitly set flags (vs defaults)
- Support partial updates (only send changed fields to API)
- Validate inputs before API call (state, priority ranges)
- Update GraphQL queries in existing commands to fetch additional fields
- Enhance output formatters to display new fields
- Maintain table width constraints (may abbreviate state names)

**GraphQL Changes:**
- `projectUpdate` mutation: Accept project UUID + input map with optional fields (name, description, shortSummary, state, priority, initiativeId, labels)
- `projectGet` query: Include description, shortSummary, state, priority, initiative, labels in response
- `projectList` query: Include state, priority in list response

**Output Enhancements:**
- Table format: Add "State" and "Priority" columns to list view
- Detail view: Show state, priority, initiative name, labels
- JSON format: Include all fields for automation
- Plaintext: Show all fields in readable format

### Project Structure Notes

- **Files to modify:**
  - `pkg/api/queries.go` (add UpdateProject method)
  - `cmd/project.go` (add projectUpdateCmd, enhance projectGetCmd and projectListCmd GraphQL queries)
  - `pkg/output/output.go` (update formatters if needed)

- **Expected test locations:**
  - Manual testing procedures in `tests/manual_project_tests.sh`
  - Smoke tests already cover enhanced display (project list, project get)

- **Estimated effort:** 5 story points (4 hours)

- **Prerequisites:** Story 1.2 complete (requires project CRUD foundation)

### Key Code References

**Existing Patterns to Follow:**

1. **Multi-Field Update Command Pattern** (from cmd/issue.go:300-350):
   ```go
   var projectUpdateCmd = &cobra.Command{
       Use:   "update PROJECT-ID",
       Short: "Update project fields",
       Long:  `Update one or more project fields. At least one field must be provided.`,
       Run: func(cmd *cobra.Command, args []string) {
           // Validate project ID provided
           if len(args) < 1 {
               output.Error("Project ID is required", plaintext, jsonOut)
               os.Exit(1)
           }
           projectID := args[0]

           // Build input map with only changed fields
           input := make(map[string]interface{})

           if name, _ := cmd.Flags().GetString("name"); name != "" {
               input["name"] = name
           }
           if cmd.Flags().Changed("state") {
               state, _ := cmd.Flags().GetString("state")
               input["state"] = state
           }
           if cmd.Flags().Changed("priority") {
               priority, _ := cmd.Flags().GetInt("priority")
               input["priority"] = priority
           }

           // Validate at least one field
           if len(input) == 0 {
               output.Error("At least one field to update is required", plaintext, jsonOut)
               os.Exit(1)
           }

           // Validate field values
           if state, ok := input["state"]; ok {
               // Validate state
           }

           // Call API
           project, err := client.UpdateProject(context.Background(), projectID, input)
           // Display result
       },
   }
   ```

2. **Flag Changed Detection Pattern**:
   ```go
   // Use cmd.Flags().Changed() to detect explicitly set flags
   if cmd.Flags().Changed("priority") {
       priority, _ := cmd.Flags().GetInt("priority")
       input["priority"] = priority
   }
   // This prevents sending default values (e.g., priority 0) when user didn't specify
   ```

3. **Enhanced GraphQL Query Pattern** (update existing commands):
   ```go
   // In projectGetCmd - enhance query
   query := fmt.Sprintf(`
       query {
           project(id: "%s") {
               id
               name
               state          # ADD
               priority       # ADD
               url
               team { id key name }
               initiative { id name }  # ADD
               labels         # ADD (array)
               createdAt
               updatedAt
           }
       }
   `, projectID)
   ```

4. **Table Output Enhancement Pattern**:
   ```go
   // In projectListCmd - add columns
   table := tablewriter.NewWriter(os.Stdout)
   table.SetHeader([]string{"ID", "Name", "Team", "State", "Priority", "URL"})  # Add State, Priority

   for _, project := range projects {
       table.Append([]string{
           project.ID,
           project.Name,
           project.Team.Key,
           project.State,          # Add state
           strconv.Itoa(project.Priority),  # Add priority
           project.URL,
       })
   }
   ```

5. **Validation Patterns** (from tech-spec section 5.2):
   ```go
   // State validation
   allowedStates := []string{"planned", "started", "paused", "completed", "canceled"}
   if state, ok := input["state"].(string); ok {
       if !contains(allowedStates, state) {
           output.Error(fmt.Sprintf("Invalid state. Must be one of: %v", allowedStates), plaintext, jsonOut)
           os.Exit(1)
       }
   }

   // Priority validation
   if priority, ok := input["priority"].(int); ok {
       if priority < 0 || priority > 4 {
           output.Error("Priority must be between 0 and 4", plaintext, jsonOut)
           os.Exit(1)
       }
   }
   ```

**Relevant Code Locations:**
- `cmd/project.go:50-100` - projectListCmd (enhance GraphQL query)
- `cmd/project.go:150-200` - projectGetCmd (enhance GraphQL query)
- `pkg/api/queries.go` - Add UpdateProject() method
- `pkg/output/output.go` - Update formatters if needed

---

## Context References

**Tech-Spec:** [tech-spec.md](../tech-spec.md) - Primary context document containing:

- **Section 2.3 "Project Update Implementation"** - Complete code example with multi-field support
- **Section 2.5 "Use Existing Framework Versions"** - Framework details
- **Section 4.3 "Project Update (GraphQL)"** - Complete mutation schema
- **Section 5.2 "Data Validation Rules"** - All validation requirements (state, priority, etc.)
- **Section 5.3 "Error Scenarios & Handling"** - Error handling patterns
- **Section 6.1 "Files to Modify"** - Complete file modification list
- **Section 7.3 "Flag Handling Pattern"** - Flag.Changed() detection
- **Section 7.6 "Priority Validation"** - Priority range validation
- **Section 7.7 "State Validation"** - State value validation
- **Section 9.5 "Story 3 Implementation Steps"** - Step-by-step guide
- **Section 9.2 "Testing Strategy"** - Story 3 test cases
- **Section 11 "UX/UI Considerations"** - Output format design

**Architecture:** See tech-spec.md sections:
- "Existing Codebase Structure" - File organization
- "Integration Points" - GraphQL API details
- "Technical Approach" - Implementation strategy

---

## Dev Agent Record

### Context Reference

- [Story Context XML](./1-3-project-updates-enhanced-display.context.xml) - Generated 2025-11-06

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
