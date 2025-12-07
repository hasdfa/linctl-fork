# Project CRUD Implementation Scratchpad

## Execution Status

| Task | Status | Commit |
|------|--------|--------|
| Analysis & Planning | ✅ Complete | - |
| Part 1: Project Create | ✅ Complete | e34cd26 |
| Part 2: Project Update | ✅ Complete | c0c6427 |
| Part 3: Project Delete | ✅ Complete | 9e8dacc |
| Final Validation | ✅ Complete | - |

## Validation Checklist (Run After Each Task)
- [x] `go build` succeeds
- [x] `go vet ./...` passes
- [x] `gofmt` formatting verified
- [x] All command help text verified
- [x] Code review passed (9.5/10 quality score)

---

## Testing Structure Analysis

### Current State: No Unit Tests

**Key Finding:** The linctl project currently has **zero Go unit test files** (`*_test.go`). All testing is done through:
1. **Smoke Tests** (bash script) - Integration-level tests against real Linear API
2. **Manual Testing** - Developer validation during development

### Existing Test Infrastructure

#### 1. Smoke Test Script (`smoke_test.sh`)
- **Location:** `/Users/m/workspace/linctl-feat-project-crud/smoke_test.sh`
- **Type:** Bash-based integration testing
- **Scope:** Read-only commands only (GET operations)
- **Coverage:** 185 lines covering:
  - Authentication validation
  - User commands (whoami, list)
  - Team commands (list, get, members)
  - **Project commands** (list, get with filters)
  - Issue commands (list, get, search)
  - Comment commands (list)
  - Help text validation
  - Error handling for unknown commands

**Project-Specific Tests in Smoke Script:**
- Line 103-117: Project list/get commands
  - `project list` (default, plaintext, json)
  - `project list --state started`
  - `project list --newer-than 1_month_ago`
  - `project get PROJECT_ID`

#### 2. Make Test Target
- **Command:** `make test`
- **Action:** Executes `./smoke_test.sh`
- **Requirements:**
  - Must be authenticated (`linctl auth`)
  - Real Linear API key required

#### 3. Test Environment Configuration
- **File:** `.env.test.example`
- **Purpose:** Template for integration test credentials
- **Variables:**
  - `LINEAR_TEST_API_KEY`: API key for testing
  - `LINEAR_TEST_TEAM_ID`: (Optional) specific team for tests
  - `TEST_DEBUG`: (Optional) debug output

### Documented But Not Implemented

From `README.md` (lines 580-584):
```
### Test Structure
- `tests/unit/` - Unit tests with mocked API responses
- `tests/integration/` - End-to-end tests with real Linear API
- `tests/testutils/` - Shared test utilities and helpers
```

**Reality:** These directories **do not exist**. This is aspirational documentation.

### Testing Approach for New Commands

Based on the existing codebase patterns, new CRUD commands should be tested using:

#### Option 1: Extend Smoke Tests (Recommended for MVP)
**Why:** Matches existing patterns, quick to implement, real-world validation

**How to extend for project CRUD:**
1. Add to `smoke_test.sh` after existing project tests (after line 117)
2. Create test project via API first
3. Test create/update/delete operations
4. Clean up test data

**Example additions:**
```bash
# Project CRUD tests section
echo -e "\n${YELLOW}Testing project CRUD commands...${NC}"

# Test project create
run_test "project create (basic)" \
  "go run main.go project create --name 'Smoke Test Project' --team $team_key --json" \
  '"name": "Smoke Test Project"'

# Capture created project ID from JSON output
created_project_id=$(go run main.go project create --name "Update Test" --team $team_key --json 2>/dev/null | jq -r '.id')

if [ -n "$created_project_id" ]; then
    # Test project update
    run_test "project update (name)" \
      "go run main.go project update $created_project_id --name 'Updated Name' --json" \
      '"name": "Updated Name"'

    # Test project delete (archive)
    run_test "project delete (archive)" \
      "echo 'y' | go run main.go project delete $created_project_id"
fi
```

#### Option 2: Create Go Unit Tests (Recommended for Future)
**Why:** Better code isolation, faster execution, no API dependencies

**Structure to create:**
```
tests/
├── unit/
│   ├── api/
│   │   └── project_test.go      # Test API client methods
│   └── cmd/
│       └── project_test.go      # Test command handlers
├── integration/
│   └── project_integration_test.go
└── testutils/
    ├── mock_api.go              # Mock Linear API responses
    └── test_helpers.go          # Shared test utilities
```

**Example Unit Test Pattern (to create):**
```go
// tests/unit/api/project_test.go
package api_test

import (
    "context"
    "testing"
    "github.com/dorkitude/linctl/pkg/api"
)

func TestCreateProject(t *testing.T) {
    // Setup mock HTTP server
    mockServer := setupMockLinearAPI(t)
    defer mockServer.Close()

    client := api.NewClientWithURL(mockServer.URL, "test-token")

    input := map[string]interface{}{
        "name": "Test Project",
        "teamIds": []string{"team-123"},
    }

    project, err := client.CreateProject(context.Background(), input)

    if err != nil {
        t.Fatalf("CreateProject failed: %v", err)
    }

    if project.Name != "Test Project" {
        t.Errorf("Expected name 'Test Project', got '%s'", project.Name)
    }
}
```

### Testing Anti-Patterns to Avoid

1. **Don't test against production data** - Use dedicated test workspace or mocks
2. **Don't leave orphaned test resources** - Always clean up created projects
3. **Don't hardcode IDs** - Extract from API responses dynamically
4. **Don't skip error cases** - Test both success and failure paths

### Recommended Testing Strategy for Project CRUD

#### Phase 1: Immediate (For PR approval)
1. **Extend smoke_test.sh** with project CRUD operations
2. Add manual test checklist to PR description
3. Document test commands in implementation plan

#### Phase 2: Near-term (Follow-up PR)
1. Create `tests/` directory structure
2. Add unit tests for API methods
3. Mock GraphQL responses for consistency
4. Add GitHub Actions workflow for automated tests

#### Phase 3: Long-term (Future enhancement)
1. Add integration tests with dedicated test workspace
2. Create test data factories for common scenarios
3. Add coverage reporting
4. Implement table-driven tests for edge cases

### Test Coverage Gaps

**Current gaps that affect project CRUD:**
- No tests for GraphQL mutation operations (create, update, delete)
- No tests for input validation
- No tests for error handling (API failures, network issues)
- No tests for flag parsing and validation
- No tests for output formatting (table, JSON, plaintext)

### Key Testing Utilities Needed

**For project CRUD specifically:**
1. **Project Factory:** Helper to create test projects with common defaults
2. **Team Lookup Mock:** Mock team resolution to avoid API calls
3. **User Lookup Mock:** Mock user resolution for lead assignment
4. **State Validator:** Test valid/invalid state transitions
5. **Cleanup Helper:** Ensure test projects are archived/deleted

### Success Criteria for Testing

**Minimum viable testing for project CRUD PR:**
- [ ] All project CRUD commands added to smoke_test.sh
- [ ] Smoke tests pass on clean workspace
- [ ] Manual testing checklist completed
- [ ] Error cases documented and tested manually

**Complete testing (future work):**
- [ ] Unit tests for all API methods (CreateProject, UpdateProject, DeleteProject)
- [ ] Unit tests for command handlers
- [ ] Integration tests with mocked GraphQL
- [ ] CI/CD pipeline running tests on PRs
- [ ] Code coverage >70%

### Test Execution Commands

```bash
# Current approach
make test                          # Runs smoke_test.sh
bash -x smoke_test.sh              # Verbose smoke test output

# Future approach (once implemented)
go test ./...                      # All Go tests
go test -v ./tests/unit/...        # Unit tests only
go test -v ./tests/integration/... # Integration tests only
go test -cover ./...               # With coverage report
```

### Dependencies and Tools

**Already available:**
- Go 1.23+ (supports table-driven tests, subtests)
- `jq` for JSON parsing in smoke tests
- `make` for test orchestration

**Needed for comprehensive testing:**
- `testify` package for assertions (optional but recommended)
- `httptest` package for mocking HTTP (stdlib)
- `github.com/stretchr/testify/mock` for interface mocking (optional)
- Coverage tools: `go test -cover`, `go tool cover`

### Related Files to Review for Testing Patterns

- `/Users/m/workspace/linctl-feat-project-crud/smoke_test.sh` - Current test patterns
- `/Users/m/workspace/linctl-feat-project-crud/pkg/api/client.go` - HTTP client suitable for mocking
- `/Users/m/workspace/linctl-feat-project-crud/pkg/api/queries.go` - GraphQL queries to mock
- `/Users/m/workspace/linctl-feat-project-crud/.env.test.example` - Test configuration

### Conclusion

**For the project CRUD implementation:**

1. **Test via smoke_test.sh extension** (quickest path to PR)
2. **Add comprehensive manual test cases** in PR description
3. **Create unit test framework** as follow-up work
4. **Document test approach** in commit messages

The project intentionally keeps testing lightweight (smoke tests only) to maintain velocity. This is acceptable for a CLI tool with:
- Read-heavy operations (safe to test against real API)
- Simple GraphQL mutations (linear API is well-tested)
- Human-in-the-loop validation (users will catch issues quickly)

However, for production-grade reliability, unit tests should be added in a follow-up PR after the CRUD commands are proven functional via smoke testing.

## Codebase Pattern Analysis - pkg/api/queries.go

### Method Structure Patterns

#### Query Methods (Read Operations)
**Pattern**: `Get<Entity>` or `Get<Entities>` naming convention

**Single Entity Retrieval** (Lines 554-820, 928-1082):
```go
func (c *Client) Get<Entity>(ctx context.Context, id string) (*<Entity>, error)
```
- Takes `context.Context` as first parameter (always)
- Takes identifying parameter (`id string`, `email string`, or `key string`)
- Returns pointer to entity type and error
- GraphQL query embedded as multiline string constant
- Query uses variables: `$id: String!` or similar
- Response struct wraps entity: `struct { <Entity> <Entity> \`json:"entity"\` }`
- Calls `c.Execute(ctx, query, variables, &response)`
- Returns `&response.<Entity>, nil` on success

**Multiple Entity Retrieval** (Lines 400-470, 823-863, 866-925):
```go
func (c *Client) Get<Entities>(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string) (*<Entities>, error)
```
- Takes optional `filter map[string]interface{}` for filtering
- Takes pagination params: `first int`, `after string`, `orderBy string`
- Returns pointer to paginated collection type (e.g., `*Issues`, `*Projects`)
- Conditionally adds variables only if they're non-empty/non-nil
- Response includes `PageInfo` for cursor-based pagination

#### Mutation Methods (Write Operations)
**Pattern**: `<Action><Entity>` naming convention (CreateIssue, UpdateIssue, NOT IssueCreate/IssueUpdate)

**Create Mutations** (Lines 1147-1205):
```go
func (c *Client) Create<Entity>(ctx context.Context, input map[string]interface{}) (*<Entity>, error)
```
- Takes `context.Context` first
- Takes `input map[string]interface{}` containing creation data
- Returns pointer to created entity and error
- GraphQL mutation keyword used instead of query
- Mutation name follows Linear's convention: `<entity>Create` (e.g., `issueCreate`)
- Response struct nests entity inside mutation result:
  ```go
  var response struct {
      <Entity>Create struct {
          Issue <Entity> `json:"<entity>"`
      } `json:"<entity>Create"`
  }
  ```
- Returns `&response.<Entity>Create.<Entity>, nil`

**Update Mutations** (Lines 1085-1144):
```go
func (c *Client) Update<Entity>(ctx context.Context, id string, input map[string]interface{}) (*<Entity>, error)
```
- Takes `context.Context` first
- Takes `id string` to identify entity to update
- Takes `input map[string]interface{}` containing update data
- Returns pointer to updated entity and error
- Mutation name: `<entity>Update` with both `$id` and `$input` variables
- Response structure identical to Create pattern

### GraphQL Query/Mutation Patterns

#### Query Structure
```go
query := `
    <query|mutation> <OperationName>($var1: Type!, $var2: Type) {
        <operation>(var1: $var1, var2: $var2) {
            <field1>
            <field2>
            <nestedField> {
                <subField>
            }
        }
    }
`
```

#### Variables Construction
```go
variables := map[string]interface{}{
    "requiredParam": value,
}
if optionalParam != "" {
    variables["optionalParam"] = optionalParam
}
```
**Pattern**: Required params always included, optional params only added if non-empty/non-nil

#### Field Selection
- Queries request specific fields needed (not all available fields)
- Related entities use nested field selection (e.g., `assignee { id name email }`)
- Paginated collections include both `nodes` array and `pageInfo` object
- Mutation responses request same fields as corresponding Get query for consistency

### Response Parsing Patterns

#### Anonymous Struct Wrapping
**All methods** use inline anonymous structs to match GraphQL response shape:
```go
var response struct {
    <TopLevelField> <Type> `json:"<jsonFieldName>"`
}
```

For mutations with nested results:
```go
var response struct {
    <MutationName> struct {
        Success bool     `json:"success"` // Optional
        <Entity> <Type>  `json:"<entity>"`
    } `json:"<mutationName>"`
}
```

#### Execute Method Call
```go
err := c.Execute(ctx, query, variables, &response)
if err != nil {
    return nil, err
}
```
**Pattern**: Always pass pointer to response struct, check error before accessing data

#### Return Pattern
- Queries return: `&response.<Field>, nil`
- Mutations return: `&response.<MutationName>.<Entity>, nil`
- Delete operations could return `error` only (no entity return needed)

### Error Handling Patterns

**Consistent pattern across all methods**:
1. Check error from `c.Execute()` immediately
2. Return `nil, err` or just `err` for void operations
3. No custom error wrapping at API layer (done in `client.go:61-112`)
4. GraphQL errors handled by `client.Execute()` automatically (Lines 101-103 in client.go)

### Type Patterns

#### Project Struct (Lines 95-128 in queries.go)
- Uses pointer types for optional fields: `*string`, `*User`, `*time.Time`
- Uses concrete types for required fields: `string`, `float64`, `time.Time`
- Nested collections use custom paginated types: `*Teams`, `*Users`, `*Issues`
- JSON tags match Linear's GraphQL field names exactly (camelCase)

#### Paginated Collection Pattern (Lines 130-149)
```go
type <Entities> struct {
    Nodes    []<Entity> `json:"nodes"`
    PageInfo PageInfo   `json:"pageInfo"`
}

type PageInfo struct {
    HasNextPage bool   `json:"hasNextPage"`
    EndCursor   string `json:"endCursor"`
}
```

### Key Implementation Requirements for New Project Mutations

1. **CreateProject** must follow:
   - Method signature: `func (c *Client) CreateProject(ctx context.Context, input map[string]interface{}) (*Project, error)`
   - Mutation name: `projectCreate` (not `createProject`)
   - Response: `ProjectCreate.Project` (Linear may include `.success` field)
   - Place after `GetProject` method (~line 1082)

2. **UpdateProject** must follow:
   - Method signature: `func (c *Client) UpdateProject(ctx context.Context, id string, input map[string]interface{}) (*Project, error)`
   - Mutation name: `projectUpdate`
   - Variables: both `$id` and `$input`
   - Response: `ProjectUpdate.Project`
   - Place after `CreateProject` method

3. **DeleteProject** / **ArchiveProject** considerations:
   - Linear typically uses `<entity>Archive` for soft delete
   - Hard delete may be `<entity>Delete` (verify against Linear API docs)
   - Archive should return `*Project` with `archivedAt` timestamp
   - Delete could return only `error` or success boolean
   - Place after `UpdateProject` method

4. **Field Selection for Mutations**:
   - Should match fields returned by `GetProject` for consistency
   - Core fields: id, name, description, state, progress, startDate, targetDate, url, icon, color, timestamps
   - Related entities: lead, teams, creator (with their essential fields)
   - Avoid requesting heavy nested collections (issues, documents) in mutation responses

5. **Variable Handling**:
   - Use `map[string]interface{}` for flexibility (matches Issue pattern)
   - Required inputs: validate at CLI layer, not API layer
   - Optional inputs: only include if provided
   - Date fields: strings in `YYYY-MM-DD` format (already used in Project struct)

### Files to Check for Additional Context
- `/Users/m/workspace/linctl-feat-project-crud/pkg/api/client.go` - Core Execute method
- `/Users/m/workspace/linctl-feat-project-crud/cmd/project.go` - Existing commands pattern

### Anti-Patterns to Avoid
- Do NOT use positional arguments in GraphQL strings (always use variables)
- Do NOT construct GraphQL strings dynamically (security risk)
- Do NOT add custom error wrapping in API methods (handled by Execute)
- Do NOT fetch all available fields (request only what's needed)
- Do NOT break naming conventions (stick to Get/Create/Update/Delete prefix)

## Codebase Pattern Analysis - cmd/project.go

### Command Structure Patterns

#### 1. Command Definition Structure
```go
var commandNameCmd = &cobra.Command{
    Use:     "action [ARGS]",
    Aliases: []string{"alias1", "alias2"},
    Short:   "Brief description",
    Long:    `Detailed description with examples`,
    Args:    cobra.ExactArgs(n), // or cobra.NoArgs for list commands
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}
```

**Key observations:**
- Commands use `var` declarations for cobra.Command structs
- Naming convention: `{resource}{Action}Cmd` (e.g., `projectListCmd`, `projectGetCmd`)
- Common aliases: `"ls"` for list, `"show"` for get, `"new"` for create
- Args validation using `cobra.ExactArgs(1)` for commands requiring IDs

#### 2. Authentication Pattern
**Standard flow used in all commands:**
```go
plaintext := viper.GetBool("plaintext")
jsonOut := viper.GetBool("json")

// Get auth header
authHeader, err := auth.GetAuthHeader()
if err != nil {
    output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
    os.Exit(1)
}

// Create API client
client := api.NewClient(authHeader)
```

**Pattern rules:**
- Always retrieve `plaintext` and `jsonOut` flags from viper at the start
- Authentication happens before any API operations
- Use `output.Error()` helper with plaintext/jsonOut flags for consistent error formatting
- Always `os.Exit(1)` after authentication failure

#### 3. Flag Patterns

**Naming conventions:**
- Short flags use single letters: `-t`, `-s`, `-l`, `-c`, `-o`, `-n`, `-j`, `-p`
- Long flags use kebab-case: `--team`, `--state`, `--limit`, `--include-completed`
- Boolean flags for toggles: `--include-completed`, `--assign-me`
- String flags for filters and values: `--team`, `--state`, `--newer-than`
- Int flags for limits: `--limit`

**Common flag patterns:**
```go
// List command flags (projectListCmd)
projectListCmd.Flags().StringP("team", "t", "", "Filter by team key")
projectListCmd.Flags().StringP("state", "s", "", "Filter by state")
projectListCmd.Flags().IntP("limit", "l", 50, "Maximum number to return")
projectListCmd.Flags().BoolP("include-completed", "c", false, "Include completed")
projectListCmd.Flags().StringP("sort", "o", "linear", "Sort order")
projectListCmd.Flags().StringP("newer-than", "n", "", "Time filter")

// Get command typically has no additional flags (uses positional args)

// Create command flags (from issueCreateCmd)
Flags().StringP("title", "t", "", "Title")
Flags().StringP("description", "d", "", "Description")
Flags().StringP("team", "", "", "Team key")
Flags().IntP("priority", "p", 0, "Priority (0-4)")
Flags().BoolP("assign-me", "m", false, "Assign to current user")

// Update command flags (from issueUpdateCmd)
Flags().StringP("title", "t", "", "New title")
Flags().StringP("description", "d", "", "New description")
Flags().StringP("assignee", "a", "", "Assignee (email or 'me')")
Flags().StringP("state", "s", "", "State name")
```

**Flag retrieval patterns:**
```go
// Using GetString, GetInt, GetBool
teamKey, _ := cmd.Flags().GetString("team")
limit, _ := cmd.Flags().GetInt("limit")
includeCompleted, _ := cmd.Flags().GetBool("include-completed")

// Checking if flag was explicitly set (update commands)
if cmd.Flags().Changed("title") {
    title, _ := cmd.Flags().GetString("title")
    input["title"] = title
}
```

#### 4. Output Formatting Patterns

**Three output modes (always in this order):**
1. JSON output (`--json` flag)
2. Plaintext output (`--plaintext` flag)
3. Formatted/colored output (default)

**List command output pattern:**
```go
if jsonOut {
    output.JSON(projects.Nodes)
    return
} else if plaintext {
    // Markdown-style plaintext with headers and bullet points
    fmt.Println("# Projects")
    for _, project := range projects.Nodes {
        fmt.Printf("## %s\n", project.Name)
        fmt.Printf("- **ID**: %s\n", project.ID)
        // ... more fields
        fmt.Println()
    }
    fmt.Printf("\nTotal: %d projects\n", len(projects.Nodes))
    return
} else {
    // Table output with colors
    headers := []string{"Name", "State", "Lead", "Teams", "Created", "Updated", "URL"}
    rows := [][]string{}

    for _, item := range items {
        // Build rows with colored state indicators
        stateColor := color.New(color.FgGreen)
        switch item.State {
        case "planned":
            stateColor = color.New(color.FgCyan)
        // ... other states
        }

        rows = append(rows, []string{
            truncateString(item.Name, 25),
            stateColor.Sprint(item.State),
            // ... other fields
        })
    }

    output.Table(output.TableData{
        Headers: headers,
        Rows:    rows,
    }, plaintext, jsonOut)

    // Success message with count
    fmt.Printf("\n%s %d projects\n",
        color.New(color.FgGreen).Sprint("✓"),
        len(projects.Nodes))
}
```

**Create/Update command output pattern:**
```go
if jsonOut {
    output.JSON(issue)
} else if plaintext {
    fmt.Printf("Created issue %s: %s\n", issue.Identifier, issue.Title)
} else {
    fmt.Printf("%s Created issue %s: %s\n",
        color.New(color.FgGreen).Sprint("✓"),
        color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier),
        issue.Title)
    if issue.Assignee != nil {
        fmt.Printf("  Assigned to: %s\n",
            color.New(color.FgCyan).Sprint(issue.Assignee.Name))
    }
}
```

#### 5. Color Coding Conventions

**State colors (consistent across commands):**
```go
stateColor := color.New(color.FgGreen)
switch project.State {
case "planned":
    stateColor = color.New(color.FgCyan)
case "started":
    stateColor = color.New(color.FgBlue)
case "paused":
    stateColor = color.New(color.FgYellow)
case "completed":
    stateColor = color.New(color.FgGreen)
case "canceled":
    stateColor = color.New(color.FgRed)
}
```

**Other color usage:**
- Identifiers/keys: `color.FgCyan`
- Success messages: `color.FgGreen` with "✓" or "✅"
- Error messages: `color.FgRed` with "❌"
- Bold for labels: `color.Bold`
- Unassigned warnings: `color.FgYellow`

#### 6. Error Handling Pattern

**Consistent error handling across all operations:**
```go
result, err := client.SomeOperation(context.Background(), params)
if err != nil {
    output.Error(fmt.Sprintf("Failed to operation: %v", err), plaintext, jsonOut)
    os.Exit(1)
}
```

**Validation errors (before API calls):**
```go
if title == "" {
    output.Error("Title is required (--title)", plaintext, jsonOut)
    os.Exit(1)
}
```

#### 7. Helper Functions

**Used in project.go:**
- `truncateString(s string, maxLen int)` - Located in `cmd/issue.go:774`
- `priorityToString(priority int)` - Located in `cmd/issue.go:757`
- `constructProjectURL(projectID, originalURL string)` - Located in `cmd/project.go:18`

#### 8. init() Function Structure

```go
func init() {
    rootCmd.AddCommand(projectCmd)
    projectCmd.AddCommand(projectListCmd)
    projectCmd.AddCommand(projectGetCmd)

    // List command flags
    projectListCmd.Flags().StringP("team", "t", "", "Filter by team key")
    projectListCmd.Flags().StringP("state", "s", "", "Filter by state")
    projectListCmd.Flags().IntP("limit", "l", 50, "Maximum number to return")
    projectListCmd.Flags().BoolP("include-completed", "c", false, "Include completed")
    projectListCmd.Flags().StringP("sort", "o", "linear", "Sort order")
    projectListCmd.Flags().StringP("newer-than", "n", "", "Time filter")
}
```

#### 9. API Client Patterns

**Filter building for list operations:**
```go
filter := make(map[string]interface{})

if teamKey != "" {
    team, err := client.GetTeam(context.Background(), teamKey)
    if err != nil {
        output.Error(fmt.Sprintf("Failed to find team '%s': %v", teamKey, err), plaintext, jsonOut)
        os.Exit(1)
    }
    filter["team"] = map[string]interface{}{"id": team.ID}
}
```

**Input building for create operations:**
```go
input := map[string]interface{}{
    "title":  title,
    "teamId": team.ID,
}

if description != "" {
    input["description"] = description
}
```

**Input building for update operations (only changed fields):**
```go
input := make(map[string]interface{})

if cmd.Flags().Changed("title") {
    title, _ := cmd.Flags().GetString("title")
    input["title"] = title
}
```

#### 10. Implementation Checklist for New Commands

When implementing `project create` and `project update`:

- [ ] Define command variable: `projectCreateCmd`, `projectUpdateCmd`
- [ ] Retrieve plaintext/jsonOut flags from viper at start
- [ ] Implement authentication pattern with auth.GetAuthHeader()
- [ ] Create API client with api.NewClient(authHeader)
- [ ] Define and retrieve command-specific flags
- [ ] Validate required fields before API calls
- [ ] Build input map with only provided/changed fields
- [ ] Handle team/user lookups (convert keys/emails to IDs)
- [ ] Make API call with proper error handling
- [ ] Implement three output modes: JSON, plaintext, formatted
- [ ] Use appropriate color coding
- [ ] Register commands in init() function
- [ ] Define all command flags in init()
- [ ] Add examples to Long description
