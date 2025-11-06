# linctl - Technical Specification

**Author:** John
**Date:** 2025-11-06
**Project Level:** 1
**Change Type:** Feature Addition
**Development Context:** Brownfield - Existing Go CLI Application

---

## Context

### Available Documents

**Brownfield Documentation (Comprehensive):**
- `docs-bmad/index.md` - Project documentation index and primary AI retrieval source
- `docs-bmad/architecture.md` - Detailed technical architecture (CLI layered application)
- `docs-bmad/source-tree-analysis.md` - Annotated directory tree and code structure
- `docs-bmad/development-guide.md` - Development workflow, build process, testing strategy
- `docs-bmad/project-overview.md` - High-level project information and capabilities

**Key Insights:**
- Project is a mature Go CLI tool for Linear API with established patterns
- Cobra-based command structure with consistent flag handling
- Existing project commands already in place (`project list`, `project get`)
- Issue commands support extensive operations but lack project assignment
- All commands support table, JSON, and plaintext output formats

### Project Stack

**Language & Runtime:**
- **Go 1.23.0** (toolchain go1.24.5)

**Core Dependencies (from go.mod):**
- **github.com/spf13/cobra v1.8.0** - CLI framework (command structure, flag parsing)
- **github.com/spf13/viper v1.18.2** - Configuration management (config files, environment variables)
- **github.com/olekukonko/tablewriter v0.0.5** - ASCII table output formatting
- **github.com/fatih/color v1.16.0** - Terminal colored output

**Project Type:** Command-Line Interface (CLI) Application
**API Integration:** Linear GraphQL API (https://api.linear.app/graphql)
**Authentication:** Personal API Key via Bearer token
**Build Tool:** Makefile (make build, make test, make install)
**Test Framework:** Bash-based smoke tests (tests/smoke_test.sh)

### Existing Codebase Structure

**Project Organization:**
```
linctl/
├── cmd/              # Cobra command definitions
│   ├── root.go       # Root command + global flags (3.3KB)
│   ├── auth.go       # Authentication commands (3.8KB)
│   ├── issue.go      # Issue management - LARGEST (35.6KB)
│   ├── project.go    # Project commands - EXTEND THIS (19.2KB)
│   ├── team.go       # Team commands (7.7KB)
│   ├── user.go       # User commands (8.4KB)
│   └── comment.go    # Comment commands (6.2KB)
├── pkg/              # Internal packages
│   ├── api/          # Linear GraphQL API client
│   │   ├── client.go # HTTP client setup
│   │   └── queries.go # GraphQL query/mutation builders
│   ├── auth/         # Authentication helpers
│   │   └── auth.go   # API key management
│   ├── output/       # Output formatters
│   │   └── output.go # Table, JSON, plaintext formatters
│   └── utils/        # Utilities
│       └── time.go   # Time parsing (newer-than expressions)
├── tests/            # Test suite
│   └── smoke_test.sh # Read-only command smoke tests
├── main.go           # Application entry point
├── Makefile          # Build automation
└── go.mod            # Go module definition
```

**Key Patterns Detected:**

1. **Command Structure Pattern:**
   - Each command in separate file (e.g., `cmd/project.go`)
   - Export command variable: `var projectCmd = &cobra.Command{...}`
   - Subcommands added in `init()` function
   - Consistent flag handling using `cmd.Flags().GetString()`

2. **Error Handling Pattern:**
   ```go
   if err != nil {
       output.Error(fmt.Sprintf("message: %v", err), plaintext, jsonOut)
       os.Exit(1)
   }
   ```

3. **API Client Pattern:**
   ```go
   authHeader, err := auth.GetAuthHeader()
   client := api.NewClient(authHeader)
   result, err := client.MethodName(context.Background(), params)
   ```

4. **Output Format Pattern:**
   - Check `plaintext` and `jsonOut` boolean flags
   - Call appropriate formatter from `output` package
   - All commands support: table (default), JSON (`--json`), plaintext (`--plaintext`)

5. **Naming Conventions:**
   - Files: lowercase.go (project.go, auth.go)
   - Commands: camelCase (projectListCmd, issueCreateCmd)
   - Functions: camelCase (constructProjectURL, GetAuthHeader)
   - Packages: lowercase (api, auth, output, utils)
   - Indentation: Tabs (Go standard)

6. **GraphQL Query Pattern (from existing code):**
   ```go
   query := fmt.Sprintf(`
       query {
           %s(filter: %s, first: %d) {
               nodes {
                   id
                   field1
                   field2
               }
           }
       }
   `, entityType, filterJSON, limit)
   ```

---

## The Change

### Problem Statement

Users of linctl cannot currently perform comprehensive project management operations. Specifically:

1. **Issue-Project Assignment Gap:** When creating or updating issues, users cannot assign them to Linear projects, forcing them to switch to the Linear web UI for this common workflow step.

2. **Incomplete Project Management:** The existing `project` commands (`list`, `get`) are read-only. Users cannot create, update, or archive projects from the CLI, making linctl insufficient for complete project lifecycle management.

3. **Limited Project Updates:** Even if a project exists, users cannot update critical fields like state (planned/started/paused/completed/canceled), priority, initiatives, or labels without leaving the terminal.

**User Impact:**
- Developers must context-switch to web UI during issue creation
- Project managers cannot manage project lifecycle from terminal
- CI/CD automation cannot create or update projects programmatically
- Incomplete CLI coverage reduces linctl's value proposition

### Proposed Solution

Extend linctl with comprehensive project management capabilities:

**1. Issue-Project Assignment (Extend Existing Commands):**
- Add `--project` flag to `linctl issue create`
- Add `--project` flag to `linctl issue update`
- Support `--project unassigned` to remove project assignment

**2. Project Creation:**
- Implement `linctl project create` with required fields (name, team)
- Support optional fields at creation (description, state, priority, etc.)

**3. Project Updates:**
- Implement `linctl project update PROJECT-ID` with multiple field flags
- Support updating: name, description, shortSummary, state, priority, initiative assignment, labels
- Enable multi-field updates in single command

**4. Project Archival:**
- Implement `linctl project archive PROJECT-ID`
- Follow Linear's archival semantics (soft delete, recoverable)

**5. Enhanced Display:**
- Update `linctl project get` to show all new fields
- Update `linctl project list` output to include state and priority
- Maintain consistent output format support (table/JSON/plaintext)

**Technical Approach:**
- Leverage Linear's GraphQL API mutations (issueCreate, issueUpdate, projectCreate, projectUpdate, projectArchive)
- Follow existing linctl patterns for command structure and error handling
- Reuse existing API client infrastructure
- Maintain backward compatibility (no breaking changes)

### Scope

**In Scope:**

✅ **Issue-Project Assignment:**
- `--project PROJECT-ID` flag on `issue create` command
- `--project PROJECT-ID` flag on `issue update` command
- `--project unassigned` to remove project assignment
- GraphQL mutations: issueCreate, issueUpdate with projectId field

✅ **Project CRUD Operations:**
- `project create --name "..." --team TEAM-KEY` command
- `project update PROJECT-ID --field value` command with multiple fields:
  - `--name` (string)
  - `--description` (string) - Full project description
  - `--summary` (string) - Short summary text (shortSummary field)
  - `--state` (planned|started|paused|completed|canceled)
  - `--priority` (0-4: None, Urgent, High, Normal, Low)
  - `--initiative INITIATIVE-ID` (parent initiative UUID)
  - `--label "tag1,tag2"` (comma-separated labels)
- `project archive PROJECT-ID` command
- GraphQL mutations: projectCreate, projectUpdate, projectArchive

✅ **Enhanced Display:**
- Update `project get` output to show: description, shortSummary, state, priority, initiative, labels
- Update `project list` output to include state and priority columns
- Maintain table/JSON/plaintext output format support

✅ **Testing:**
- Smoke tests for all new commands
- Manual testing documentation
- Error handling validation

**Out of Scope:**

❌ **Project Un-archival:**
- Not implementing restore/unarchive (rare use case, can be manual)

❌ **Advanced Project Features:**
- Project milestones management (future enhancement)
- Project roadmap operations (future enhancement)
- Project templates (future enhancement)
- Bulk project operations (future enhancement)

❌ **Issue Filtering by Project:**
- Not adding `--project` filter to `issue list` in this iteration
- Can be added in future enhancement based on user feedback

❌ **Project Members Management:**
- Not managing project lead/members assignments (future enhancement)

❌ **UI Changes:**
- No changes to linctl's terminal UI framework
- Maintaining existing table/JSON/plaintext formatters

---

## Implementation Details

### Source Tree Changes

**Files to MODIFY:**

1. **cmd/issue.go** (MODIFY - Add --project flag support)
   - Line ~200-250: Add `--project` flag to `issueCreateCmd`
   - Line ~300-350: Add `--project` flag to `issueUpdateCmd`
   - Add project ID validation and handling logic
   - Update GraphQL mutation input to include projectId field

2. **cmd/project.go** (MODIFY - Add create/update/archive commands)
   - Add new command: `projectCreateCmd` (after existing projectListCmd)
   - Add new command: `projectUpdateCmd` (after projectGetCmd)
   - Add new command: `projectArchiveCmd` (end of file)
   - Update `init()` function to register new subcommands
   - Enhance `projectGetCmd` output to show new fields

3. **pkg/api/queries.go** (MODIFY - Add GraphQL mutations)
   - Add function: `CreateProject(ctx, input) (*Project, error)`
   - Add function: `UpdateProject(ctx, id, input) (*Project, error)`
   - Add function: `ArchiveProject(ctx, id) (bool, error)`
   - Add function: `UpdateIssue(ctx, id, input) (*Issue, error)` (if not exists)
   - Define mutation query strings for each operation

4. **pkg/api/client.go** (MODIFY - Extend Client type if needed)
   - Add Project type definition if missing (struct with all fields)
   - Add ProjectInput type for create/update operations
   - Ensure Issue type has ProjectId field

**Files to CREATE:**

None - All functionality fits into existing file structure.

**Files to READ (for reference):**

- `cmd/issue.go` lines 200-400: Existing issue create/update patterns
- `cmd/project.go` lines 50-180: Existing project list/get implementations
- `pkg/api/queries.go`: Existing GraphQL query patterns
- `pkg/output/output.go`: Output formatter interfaces

### Technical Approach

**1. Issue-Project Assignment Implementation:**

Add `--project` flag to issue commands following existing flag pattern:

```go
// In cmd/issue.go - issueCreateCmd.Flags() section
issueCreateCmd.Flags().String("project", "", "Project ID to assign issue to")

// In Run function:
projectID, _ := cmd.Flags().GetString("project")
if projectID != "" {
    input["projectId"] = projectID
}
```

**2. Project Creation Implementation:**

Create new Cobra command following existing pattern in cmd/project.go:

```go
var projectCreateCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new project",
    Long:  `Create a new project in Linear workspace.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Get flags
        name, _ := cmd.Flags().GetString("name")
        teamKey, _ := cmd.Flags().GetString("team")

        // Validate required fields
        if name == "" || teamKey == "" {
            output.Error("Both --name and --team are required", plaintext, jsonOut)
            os.Exit(1)
        }

        // Get auth and create client
        authHeader, err := auth.GetAuthHeader()
        client := api.NewClient(authHeader)

        // Get team ID from key
        team, err := client.GetTeam(context.Background(), teamKey)

        // Build input
        input := map[string]interface{}{
            "name": name,
            "teamId": team.ID,
        }

        // Add optional fields if provided
        if state, _ := cmd.Flags().GetString("state"); state != "" {
            input["state"] = state
        }

        // Call API
        project, err := client.CreateProject(context.Background(), input)

        // Format and display output
        output.ProjectDetails(project, plaintext, jsonOut)
    },
}
```

**3. Project Update Implementation:**

Multi-field update following existing update patterns:

```go
var projectUpdateCmd = &cobra.Command{
    Use:   "update PROJECT-ID",
    Short: "Update project fields",
    Run: func(cmd *cobra.Command, args []string) {
        // Validate project ID provided
        if len(args) < 1 {
            output.Error("Project ID is required", plaintext, jsonOut)
            os.Exit(1)
        }
        projectID := args[0]

        // Build input map with only provided flags
        input := make(map[string]interface{})

        if name, _ := cmd.Flags().GetString("name"); name != "" {
            input["name"] = name
        }
        if state, _ := cmd.Flags().GetString("state"); state != "" {
            input["state"] = state
        }
        if priority, _ := cmd.Flags().GetInt("priority"); cmd.Flags().Changed("priority") {
            input["priority"] = priority
        }

        // Validate at least one field provided
        if len(input) == 0 {
            output.Error("At least one field to update is required", plaintext, jsonOut)
            os.Exit(1)
        }

        // Call API
        project, err := client.UpdateProject(context.Background(), projectID, input)

        // Display result
        output.ProjectDetails(project, plaintext, jsonOut)
    },
}
```

**4. GraphQL Mutation Implementation:**

In pkg/api/queries.go, add mutations following existing pattern:

```go
func (c *Client) CreateProject(ctx context.Context, input map[string]interface{}) (*Project, error) {
    inputJSON, _ := json.Marshal(input)

    query := fmt.Sprintf(`
        mutation {
            projectCreate(input: %s) {
                success
                project {
                    id
                    name
                    state
                    priority
                    url
                    team { id key name }
                    initiative { id name }
                    createdAt
                    updatedAt
                }
            }
        }
    `, string(inputJSON))

    result, err := c.query(ctx, query)
    // Parse result and return Project
}
```

**5. Use Existing Framework Versions:**

All operations use existing dependencies from go.mod:
- Cobra v1.8.0 for command handling
- Viper v1.18.2 for flag management
- Native Go encoding/json for GraphQL query building
- tablewriter v0.0.5 for table output

### Existing Patterns to Follow

**From cmd/project.go (existing code):**

1. **Command Registration Pattern:**
   ```go
   func init() {
       projectCmd.AddCommand(projectListCmd)
       projectCmd.AddCommand(projectGetCmd)
       // Add new commands here:
       projectCmd.AddCommand(projectCreateCmd)
       projectCmd.AddCommand(projectUpdateCmd)
       projectCmd.AddCommand(projectArchiveCmd)
   }
   ```

2. **Team ID Resolution Pattern:**
   ```go
   // Always resolve team key to team ID before API calls
   team, err := client.GetTeam(context.Background(), teamKey)
   if err != nil {
       output.Error(fmt.Sprintf("Failed to find team '%s': %v", teamKey, err), plaintext, jsonOut)
       os.Exit(1)
   }
   filter["team"] = map[string]interface{}{"id": team.ID}
   ```

3. **Flag Handling Pattern:**
   ```go
   // Check if flag was explicitly set vs using default
   if cmd.Flags().Changed("priority") {
       priority, _ := cmd.Flags().GetInt("priority")
       input["priority"] = priority
   }
   ```

4. **Error Handling Pattern:**
   ```go
   // Always use output.Error for consistency
   if err != nil {
       output.Error(fmt.Sprintf("Failed to X: %v", err), plaintext, jsonOut)
       os.Exit(1)
   }
   ```

**From cmd/issue.go (existing patterns):**

5. **Optional Field Handling:**
   ```go
   // Only add to input map if flag provided
   if description, _ := cmd.Flags().GetString("description"); description != "" {
       input["description"] = description
   }
   ```

6. **Priority Validation:**
   ```go
   // Priority must be 0-4
   if priority < 0 || priority > 4 {
       output.Error("Priority must be between 0 (None) and 4 (Low)", plaintext, jsonOut)
       os.Exit(1)
   }
   ```

7. **State Validation:**
   ```go
   // Validate state is one of allowed values
   allowedStates := []string{"planned", "started", "paused", "completed", "canceled"}
   if !contains(allowedStates, state) {
       output.Error(fmt.Sprintf("Invalid state. Must be one of: %v", allowedStates), plaintext, jsonOut)
       os.Exit(1)
   }
   ```

### Integration Points

**1. Linear GraphQL API Endpoints:**
- **Base URL:** https://api.linear.app/graphql
- **Authentication:** Bearer token from `~/.linctl-auth.json`
- **Timeout:** 30 seconds (from existing client)
- **Retry:** 3 attempts (from existing client)

**2. GraphQL Mutations to Implement:**

```graphql
# Issue mutations (extend existing)
mutation IssueCreate {
    issueCreate(input: {
        title: "..."
        teamId: "..."
        projectId: "..."  # NEW FIELD
    }) {
        success
        issue { id identifier projectId }
    }
}

mutation IssueUpdate {
    issueUpdate(
        id: "ISS-123"
        input: { projectId: "..." }  # NEW FIELD
    ) {
        success
        issue { id projectId }
    }
}

# Project mutations (new)
mutation ProjectCreate {
    projectCreate(input: {
        name: "..."
        teamId: "..."
        state: "planned"
        priority: 0
    }) {
        success
        project {
            id name state priority url
            team { id key name }
            createdAt
        }
    }
}

mutation ProjectUpdate {
    projectUpdate(
        id: "PROJECT-UUID"
        input: {
            name: "..."
            state: "started"
            priority: 1
            initiativeId: "..."
        }
    ) {
        success
        project { id name state priority }
    }
}

mutation ProjectArchive {
    projectArchive(id: "PROJECT-UUID") {
        success
    }
}
```

**3. Internal Module Dependencies:**

- **pkg/auth:** Use `GetAuthHeader()` for API authentication
- **pkg/api:** Extend `Client` type with new mutation methods
- **pkg/output:** Use existing formatters (no changes needed)
- **pkg/utils:** May need validation helper functions

**4. Data Flow:**

```
User Command → Cobra Handler → Flag Parsing →
  → Team ID Resolution (if needed) →
  → Input Map Construction →
  → API Client Method →
  → GraphQL Mutation →
  → Response Parsing →
  → Output Formatter →
  → Terminal Display
```

**5. Error Propagation:**

- API errors bubble up from `client.query()` method
- GraphQL errors extracted from response JSON
- User-facing errors formatted via `output.Error()`
- Exit code 1 on all errors (existing pattern)

---

## Development Context

### Relevant Existing Code

**Key Files to Reference:**

1. **cmd/issue.go (lines 100-400):**
   - `issueCreateCmd` - Pattern for creating entities with multiple flags
   - `issueUpdateCmd` - Pattern for updating entities with optional fields
   - Flag registration and validation logic
   - Input map construction from flags

2. **cmd/project.go (lines 50-180):**
   - `projectListCmd` - Existing list implementation with filters
   - `projectGetCmd` - Existing detail view implementation
   - Team key to team ID resolution pattern
   - Project URL construction logic

3. **pkg/api/queries.go:**
   - `GetTeam()` method - Pattern for entity retrieval
   - `query()` method - Low-level GraphQL execution
   - Response parsing and error handling

4. **pkg/api/client.go:**
   - `Client` struct - API client configuration
   - HTTP client setup with timeout and retry
   - Authentication header handling

### Dependencies

**Framework/Libraries (from go.mod):**

- **Go 1.23.0** - Language and standard library
- **github.com/spf13/cobra v1.8.0** - CLI framework for command structure
- **github.com/spf13/viper v1.18.2** - Configuration and flag management
- **github.com/olekukonko/tablewriter v0.0.5** - ASCII table formatting
- **github.com/fatih/color v1.16.0** - Terminal color output
- **encoding/json** (stdlib) - JSON marshaling for GraphQL
- **context** (stdlib) - Context for API calls
- **fmt** (stdlib) - String formatting
- **os** (stdlib) - Exit codes and environment

**Internal Modules:**

- **pkg/auth** - Authentication token management
  - `GetAuthHeader()` - Returns Bearer token for API calls
- **pkg/api** - Linear GraphQL API client
  - `Client` type - API client with methods
  - `query()` - Execute GraphQL queries/mutations
- **pkg/output** - Output formatting
  - `Error()` - Consistent error display
  - `ProjectDetails()` - Project detail formatter (may need extension)
- **pkg/utils** - Utility functions
  - Time parsing utilities
  - May need validation helpers

### Configuration Changes

**No Configuration File Changes Required**

The feature uses existing configuration:
- `~/.linctl.yaml` - No new config keys needed
- `~/.linctl-auth.json` - Uses existing auth token
- Environment variables - No new variables

**Command Help Text Updates:**

Update command help strings in:
- `cmd/issue.go` - Add `--project` flag documentation
- `cmd/project.go` - Add new command documentation

**README.md Updates (Post-Implementation):**

Add documentation for:
- New `--project` flag on issue commands
- New `project create` command with examples
- New `project update` command with examples
- New `project archive` command with examples

### Existing Conventions (Brownfield)

**Code Style (Must Follow):**
- **Indentation:** Tabs (Go standard via `gofmt`)
- **Imports:** Group stdlib, third-party, then local packages
- **Comments:** Function-level doc comments for exported items
- **Error Handling:** Always check errors, use `output.Error()` + `os.Exit(1)`
- **Naming:** camelCase for functions/variables, PascalCase for types

**Command Patterns (Must Follow):**
- Each command group in separate file (`cmd/project.go`, `cmd/issue.go`)
- Export command variables (`var projectCreateCmd = &cobra.Command{...}`)
- Register in `init()` function
- Consistent flag handling with Viper
- Support `--json`, `--plaintext` on all list/get commands

**API Client Patterns (Must Follow):**
- Always use `context.Background()` for API calls
- Resolve team keys to UUIDs before API calls
- Build input maps dynamically from flags
- Parse response JSON and extract data/errors
- Return structured types, not raw JSON

**Testing Patterns:**
- Bash-based smoke tests in `tests/smoke_test.sh`
- Test read-only commands automatically
- Write commands tested manually
- No Go unit tests currently (future enhancement)

### Test Framework & Standards

**Current Test Infrastructure:**

**Framework:** Bash shell scripts (tests/smoke_test.sh)
- Tests all read-only commands (`list`, `get`)
- Verifies command exits cleanly (exit code 0)
- Checks help text is accessible
- Uses real Linear API (requires `LINEAR_API_KEY`)

**Test Organization:**
- Single file: `tests/smoke_test.sh`
- Functions per command group
- Sequential execution

**Test Coverage Requirements:**
- All read-only commands MUST have smoke tests
- Write commands tested manually (due to side effects)
- Error handling tested manually

**For This Feature:**
- Add smoke tests for: `project list`, `project get` (already exist)
- Manual test: `project create`, `project update`, `project archive`
- Manual test: `issue create --project`, `issue update --project`

---

## Implementation Stack

**Runtime Environment:**
- **Language:** Go 1.23.0
- **Toolchain:** go1.24.5
- **Minimum Go Version:** 1.23.0 (from go.mod)

**Core Framework:**
- **CLI Framework:** Cobra v1.8.0
- **Configuration:** Viper v1.18.2
- **Context:** Go stdlib context package

**API Integration:**
- **Protocol:** HTTPS + GraphQL
- **API Endpoint:** https://api.linear.app/graphql
- **Authentication:** Bearer token (Personal API Key)
- **Request Format:** JSON (GraphQL queries/mutations)
- **Response Format:** JSON

**Output Formatting:**
- **Table Output:** tablewriter v0.0.5
- **Color Output:** fatih/color v1.16.0
- **JSON Output:** encoding/json (stdlib)
- **Plaintext Output:** fmt package (stdlib)

**Build & Development:**
- **Build System:** Makefile
- **Package Manager:** Go modules (go.mod)
- **Formatter:** gofmt (built-in)
- **Linter:** golangci-lint (optional, via make lint)

**Testing:**
- **Smoke Tests:** Bash shell scripts
- **Test Runner:** make test
- **API Testing:** Real Linear API

**Distribution:**
- **Primary:** Homebrew tap (dorkitude/linctl)
- **Secondary:** Source installation via make install
- **CI/CD:** GitHub Actions (bump-tap.yml)

---

## Technical Details

### GraphQL Mutation Specifics

**1. Issue Project Assignment:**

```graphql
mutation IssueCreate {
    issueCreate(input: {
        title: String!
        teamId: String!
        projectId: String  # Optional UUID
        description: String
        priority: Int
        stateId: String
    }) {
        success
        issue {
            id
            identifier
            title
            project { id name }
        }
    }
}

mutation IssueUpdate {
    issueUpdate(
        id: String!  # Can be UUID or identifier (ISS-123)
        input: {
            projectId: String  # UUID or null to unassign
            # Other fields...
        }
    ) {
        success
        issue {
            id
            identifier
            project { id name }
        }
    }
}
```

**2. Project Creation:**

```graphql
mutation ProjectCreate {
    projectCreate(input: {
        name: String!
        teamId: String!
        description: String
        state: String  # "planned" | "started" | "paused" | "completed" | "canceled"
        priority: Int  # 0-4
        leadId: String  # User UUID
        initiativeId: String  # Initiative UUID
        targetDate: String  # ISO date
        color: String  # Hex color
    }) {
        success
        project {
            id
            name
            state
            priority
            url
            team { id key name }
            lead { id name email }
            initiative { id name }
            createdAt
            updatedAt
            targetDate
        }
    }
}
```

**3. Project Update:**

```graphql
mutation ProjectUpdate {
    projectUpdate(
        id: String!  # Project UUID
        input: {
            name: String
            description: String
            shortSummary: String  # Short summary text
            state: String
            priority: Int
            leadId: String
            initiativeId: String
            targetDate: String
            # All fields optional, only provided fields updated
        }
    ) {
        success
        project {
            id
            name
            description
            shortSummary
            state
            priority
            updatedAt
        }
    }
}
```

**4. Project Archive:**

```graphql
mutation ProjectArchive {
    projectArchive(id: String!) {
        success
        entity { id archivedAt }
    }
}
```

### Data Validation Rules

**Project Name:**
- Required for creation
- Min length: 1 character
- Max length: 255 characters
- No special validation (Linear handles)

**Project State:**
- Must be one of: `planned`, `started`, `paused`, `completed`, `canceled`
- Default: `planned` (if not specified)
- Validate before API call to give better error messages

**Priority:**
- Must be integer 0-4
- 0 = None (default)
- 1 = Urgent
- 2 = High
- 3 = Normal
- 4 = Low
- Validate range before API call

**Team Key:**
- Required for project creation
- Must exist in workspace (resolve to UUID first)
- Use existing `GetTeam()` method for resolution

**Project ID (UUID):**
- Format: Standard UUID v4 (8-4-4-4-12 hex digits)
- Accept any string, let Linear API validate
- Error handling for invalid UUIDs

**Initiative ID:**
- Optional UUID
- No local validation (Linear validates)

**Labels:**
- Comma-separated string input from user
- Split and trim in code
- Send as array to API

### Error Scenarios & Handling

**1. Authentication Errors:**
- Missing/invalid API key → "Not authenticated: Run 'linctl auth' first"
- Expired token → "Authentication failed: Please re-authenticate"

**2. Validation Errors:**
- Missing required fields → "Field X is required"
- Invalid state value → "State must be one of: planned, started, paused, completed, canceled"
- Invalid priority → "Priority must be between 0 and 4"
- Invalid team key → "Team 'X' not found"

**3. API Errors:**
- Network timeout → "Request timed out: Check network connection"
- GraphQL errors → Extract error message from response
- Rate limiting → "Rate limit exceeded: Try again in X seconds"

**4. Not Found Errors:**
- Project not found → "Project 'X' not found"
- Team not found → "Team 'X' not found"
- Initiative not found → "Initiative 'X' not found"

**5. Permission Errors:**
- Insufficient permissions → "Insufficient permissions: Contact workspace admin"

**Error Handling Pattern:**
```go
if err != nil {
    // Check for specific error types
    if strings.Contains(err.Error(), "not found") {
        output.Error(fmt.Sprintf("Resource not found: %v", err), plaintext, jsonOut)
    } else if strings.Contains(err.Error(), "unauthorized") {
        output.Error("Not authenticated: Run 'linctl auth' first", plaintext, jsonOut)
    } else {
        output.Error(fmt.Sprintf("Operation failed: %v", err), plaintext, jsonOut)
    }
    os.Exit(1)
}
```

### Performance Considerations

**API Call Optimization:**
- Single API call per command (no N+1 queries)
- Minimal field selection in GraphQL queries
- Team ID caching not needed (single resolution per command)

**Response Time:**
- Target: < 2 seconds per command
- Linear API typically responds in 200-500ms
- Network latency is primary factor

**Memory Usage:**
- Minimal (CLI tool, short-lived)
- No caching or state persistence
- Garbage collected after command completion

**Concurrency:**
- Single-threaded execution (CLI tool)
- No concurrent API calls needed
- Context timeout: 30 seconds (existing)

---

## Development Setup

**Prerequisites:**
- Go 1.23.0 or later installed
- Git for version control
- Linear workspace with admin access (for testing)
- Linear Personal API Key

**Local Development Setup:**

```bash
# 1. Clone repository (if not already)
git clone https://github.com/dorkitude/linctl.git
cd linctl

# 2. Install dependencies
make deps

# 3. Verify build
make build

# 4. Run locally without installing
go run main.go --help

# 5. Authenticate with Linear
./linctl auth
# Enter your Personal API Key when prompted

# 6. Verify authentication
./linctl whoami

# 7. Create feature branch
git checkout -b feature/project-management

# 8. Run smoke tests (requires LINEAR_API_KEY)
export LINEAR_API_KEY="your-key"
make test
```

**Development Workflow:**

```bash
# Make changes to code
vim cmd/project.go

# Format code (required before commit)
make fmt

# Build and test locally
make build
./linctl project create --name "Test" --team ENG

# Run smoke tests
make test

# Commit changes
git add .
git commit -m "Add project create command"
```

---

## Implementation Guide

### Setup Steps

**Pre-Implementation Checklist:**

1. ✅ **Review Existing Code:**
   - Read `cmd/project.go` (existing project commands)
   - Read `cmd/issue.go` (existing issue commands with flags)
   - Read `pkg/api/queries.go` (existing GraphQL patterns)

2. ✅ **Create Feature Branch:**
   ```bash
   git checkout -b feature/project-management
   ```

3. ✅ **Verify Development Environment:**
   ```bash
   make build
   make test
   ./linctl --version
   ```

4. ✅ **Set Up Test Data:**
   - Have test Linear workspace ready
   - Know team key for testing (e.g., "ENG")
   - Have a test project ID ready

5. ✅ **Review Tech-Spec:**
   - Understand all three stories
   - Review GraphQL mutations
   - Understand validation requirements

### Implementation Steps

This feature is broken into **3 user stories** for Level 1 implementation:

#### **Story 1: Issue-Project Assignment**
*Enable assigning projects to issues during creation and updates*

**Implementation Steps:**

1. **Modify cmd/issue.go - Add --project flag to issue create** (~30 minutes)
   - Locate `issueCreateCmd` variable definition
   - Add flag in `init()` or inline: `issueCreateCmd.Flags().String("project", "", "Project ID to assign issue to")`
   - In `Run` function, add project ID to input map:
     ```go
     if projectID, _ := cmd.Flags().GetString("project"); projectID != "" {
         input["projectId"] = projectID
     }
     ```
   - Update help text to document --project flag

2. **Modify cmd/issue.go - Add --project flag to issue update** (~30 minutes)
   - Locate `issueUpdateCmd` variable definition
   - Add same flag registration pattern
   - Add to input map in Run function
   - Handle "unassigned" special value:
     ```go
     if projectID := cmd.Flags().GetString("project"); projectID == "unassigned" {
         input["projectId"] = nil
     } else if projectID != "" {
         input["projectId"] = projectID
     }
     ```

3. **Test Story 1** (~30 minutes)
   - Manual test: `linctl issue create --title "Test" --team ENG --project PROJECT-UUID`
   - Manual test: `linctl issue update ISS-123 --project PROJECT-UUID`
   - Manual test: `linctl issue update ISS-123 --project unassigned`
   - Verify JSON output includes project field
   - Test error cases (invalid project ID)

**Story 1 Estimate:** 1.5 hours

---

#### **Story 2: Project Creation & Archival**
*Enable creating new projects and archiving existing ones*

**Implementation Steps:**

1. **Extend pkg/api/queries.go - Add CreateProject method** (~45 minutes)
   - Define GraphQL mutation string
   - Implement `CreateProject(ctx, input) (*Project, error)` method
   - Marshal input to JSON
   - Call `c.query(ctx, mutation)`
   - Parse response and extract project data
   - Handle errors

2. **Extend pkg/api/queries.go - Add ArchiveProject method** (~30 minutes)
   - Define GraphQL mutation string
   - Implement `ArchiveProject(ctx, id) (bool, error)` method
   - Call API and handle response

3. **Modify cmd/project.go - Add project create command** (~60 minutes)
   - Create `projectCreateCmd` variable following existing pattern
   - Add flags: `--name` (required), `--team` (required), `--description`, `--state`, `--priority`
   - Validate required fields
   - Resolve team key to team ID
   - Build input map from flags
   - Call `client.CreateProject()`
   - Format and display output
   - Register command in `init()`

4. **Modify cmd/project.go - Add project archive command** (~30 minutes)
   - Create `projectArchiveCmd` variable
   - Accept project ID as argument
   - Call `client.ArchiveProject()`
   - Display success message
   - Register command in `init()`

5. **Test Story 2** (~45 minutes)
   - Manual test: `linctl project create --name "Q1 Backend" --team ENG`
   - Manual test: `linctl project create --name "Test" --team ENG --state started --priority 1`
   - Manual test: `linctl project archive PROJECT-UUID`
   - Test all output formats (table, JSON, plaintext)
   - Test error cases (missing fields, invalid team)

**Story 2 Estimate:** 3.5 hours

---

#### **Story 3: Project Updates & Enhanced Display**
*Enable updating project fields and showing all project details*

**Implementation Steps:**

1. **Extend pkg/api/queries.go - Add UpdateProject method** (~45 minutes)
   - Define GraphQL mutation string
   - Implement `UpdateProject(ctx, id, input) (*Project, error)` method
   - Handle partial updates (only provided fields)
   - Parse response and return updated project

2. **Modify cmd/project.go - Add project update command** (~90 minutes)
   - Create `projectUpdateCmd` variable
   - Accept project ID as argument
   - Add flags: `--name`, `--state`, `--priority`, `--initiative`, `--label`
   - Validate at least one field provided
   - Validate state and priority values
   - Build input map with only changed fields (use `cmd.Flags().Changed()`)
   - Call `client.UpdateProject()`
   - Format and display output
   - Register command in `init()`

3. **Modify cmd/project.go - Enhance project get output** (~30 minutes)
   - Update GraphQL query in `projectGetCmd` to include: state, priority, initiative, labels
   - Update output formatter to display new fields
   - Maintain table/JSON/plaintext compatibility

4. **Modify cmd/project.go - Enhance project list output** (~30 minutes)
   - Update GraphQL query to include state and priority
   - Add state and priority columns to table output
   - Update JSON output to include new fields

5. **Test Story 3** (~60 minutes)
   - Manual test: `linctl project update PROJECT-UUID --name "New Name"`
   - Manual test: `linctl project update PROJECT-UUID --state started`
   - Manual test: `linctl project update PROJECT-UUID --priority 1`
   - Manual test: `linctl project update PROJECT-UUID --state started --priority 1` (multi-field)
   - Manual test: `linctl project get PROJECT-UUID` (verify new fields shown)
   - Manual test: `linctl project list` (verify state/priority columns)
   - Test all output formats
   - Test validation errors

**Story 3 Estimate:** 4.5 hours (added description and shortSummary fields)

---

**Total Implementation Time:** ~9 hours (spread across 3 stories)

### Testing Strategy

**Testing Approach:**

**1. Smoke Tests (Automated):**
- Existing smoke tests for `project list` and `project get` already cover display changes
- No new smoke tests needed (write commands have side effects)

**2. Manual Testing (Required):**

**Story 1 Test Cases:**
- ✅ Create issue with project assignment
- ✅ Create issue without project (existing behavior)
- ✅ Update issue to assign project
- ✅ Update issue to change project
- ✅ Update issue to remove project (unassigned)
- ✅ Error: Invalid project UUID
- ✅ JSON output includes project field

**Story 2 Test Cases:**
- ✅ Create project with required fields only
- ✅ Create project with all optional fields
- ✅ Archive project successfully
- ✅ Error: Missing required fields (name, team)
- ✅ Error: Invalid team key
- ✅ Error: Invalid state value
- ✅ Error: Invalid priority value
- ✅ All output formats work (table, JSON, plaintext)

**Story 3 Test Cases:**
- ✅ Update single field (name, state, priority)
- ✅ Update multiple fields at once
- ✅ Update with initiative ID
- ✅ Update with labels (comma-separated)
- ✅ Enhanced `project get` shows all fields
- ✅ Enhanced `project list` includes state and priority
- ✅ Error: No fields provided
- ✅ Error: Invalid state value
- ✅ Error: Invalid priority value
- ✅ Error: Project not found

**3. Integration Testing:**
- Test complete workflow: Create project → Create issue with project → Update project → Archive project
- Verify changes reflected in Linear web UI
- Test with multiple teams
- Test with existing vs new projects

**4. Error Handling Testing:**
- Test without authentication
- Test with expired token
- Test with insufficient permissions
- Test with network timeout (disconnect during command)
- Test with invalid UUIDs
- Test with non-existent resources

**Test Data Requirements:**
- Test Linear workspace
- At least 2 teams (for multi-team testing)
- Test project IDs (create and delete after testing)
- Test issue IDs

### Acceptance Criteria

**Story 1: Issue-Project Assignment**

Given a user wants to assign a project to an issue,
When they run `linctl issue create --title "Test" --team ENG --project PROJECT-UUID`,
Then the issue is created and assigned to the specified project.

Given a user wants to update an issue's project,
When they run `linctl issue update ISS-123 --project NEW-PROJECT-UUID`,
Then the issue's project is updated to the new project.

Given a user wants to remove a project from an issue,
When they run `linctl issue update ISS-123 --project unassigned`,
Then the issue's project assignment is removed.

**Story 2: Project Creation & Archival**

Given a user wants to create a new project,
When they run `linctl project create --name "Q1 Backend" --team ENG`,
Then a new project is created in Linear with default values.

Given a user wants to create a project with specific initial values,
When they run `linctl project create --name "Test" --team ENG --state started --priority 1`,
Then the project is created with the specified field values.

Given a user wants to archive a project,
When they run `linctl project archive PROJECT-UUID`,
Then the project is archived in Linear and success message is displayed.

**Story 3: Project Updates & Enhanced Display**

Given a user wants to update a project's name,
When they run `linctl project update PROJECT-UUID --name "New Name"`,
Then the project's name is updated in Linear.

Given a user wants to update multiple project fields,
When they run `linctl project update PROJECT-UUID --state started --priority 1 --label "urgent"`,
Then all specified fields are updated in a single API call.

Given a user wants to see full project details,
When they run `linctl project get PROJECT-UUID`,
Then all project fields are displayed including state, priority, initiative, and labels.

Given a user wants to see project states in the list view,
When they run `linctl project list`,
Then the output includes state and priority columns for each project.

**Cross-Story Acceptance Criteria:**

- ✅ All commands follow existing linctl conventions
- ✅ All commands support `--json` and `--plaintext` output
- ✅ Error messages are clear and actionable
- ✅ Help text is comprehensive for all new flags
- ✅ Code follows Go conventions and passes `gofmt`
- ✅ No breaking changes to existing commands
- ✅ README.md updated with examples

---

## Developer Resources

### File Paths Reference

**Complete list of all files involved:**

**Files to Modify:**
- `/Users/johnpates/Documents/github_clones/linctl/cmd/issue.go` (add --project flag)
- `/Users/johnpates/Documents/github_clones/linctl/cmd/project.go` (add create/update/archive commands)
- `/Users/johnpates/Documents/github_clones/linctl/pkg/api/queries.go` (add GraphQL mutations)
- `/Users/johnpates/Documents/github_clones/linctl/pkg/api/client.go` (add types if needed)
- `/Users/johnpates/Documents/github_clones/linctl/README.md` (add documentation)

**Files to Reference:**
- `/Users/johnpates/Documents/github_clones/linctl/cmd/root.go` (root command structure)
- `/Users/johnpates/Documents/github_clones/linctl/pkg/auth/auth.go` (auth patterns)
- `/Users/johnpates/Documents/github_clones/linctl/pkg/output/output.go` (output formatters)
- `/Users/johnpates/Documents/github_clones/linctl/go.mod` (dependencies)
- `/Users/johnpates/Documents/github_clones/linctl/Makefile` (build commands)

### Key Code Locations

**Important functions, structures, and patterns:**

1. **cmd/issue.go:**
   - Line ~100: `issueCreateCmd` definition
   - Line ~200: `issueUpdateCmd` definition
   - Line ~250: Flag handling patterns
   - Line ~300: Input map construction

2. **cmd/project.go:**
   - Line ~37: `projectCmd` root definition
   - Line ~50: `projectListCmd` implementation
   - Line ~150: `projectGetCmd` implementation
   - Line ~18: `constructProjectURL()` helper function
   - Bottom of file: `init()` function for command registration

3. **pkg/api/queries.go:**
   - GraphQL query/mutation functions
   - `query(ctx, graphqlQuery)` low-level method
   - Response parsing patterns

4. **pkg/api/client.go:**
   - `Client` struct definition
   - `NewClient(authHeader)` constructor
   - HTTP client configuration

5. **pkg/auth/auth.go:**
   - `GetAuthHeader()` function
   - Auth file management

### Testing Locations

**Test Structure:**

- `/Users/johnpates/Documents/github_clones/linctl/tests/smoke_test.sh` - Main smoke test file
  - Contains test functions for each command group
  - Tests read-only commands automatically
  - Write commands tested manually

**Manual Test Script Location:**
- Create manual test script at: `/Users/johnpates/Documents/github_clones/linctl/tests/manual_project_tests.sh`
- Document manual test procedures

**Test Commands to Add:**
```bash
# In smoke_test.sh (no changes needed - existing tests cover display)
# Manual testing documented in tests/manual_project_tests.sh
```

### Documentation to Update

**1. README.md** - Add new sections and examples:

Location: `/Users/johnpates/Documents/github_clones/linctl/README.md`

**Sections to Add/Update:**

- **Issue Management section** (line ~98):
  - Add `--project` flag documentation
  - Example: `linctl issue create --title "..." --team ENG --project PROJECT-ID`
  - Example: `linctl issue update ISS-123 --project PROJECT-ID`
  - Example: `linctl issue update ISS-123 --project unassigned`

- **Project Tracking section** (line ~149):
  - Add project create documentation
  - Add project update documentation
  - Add project archive documentation
  - Examples for each command with all flags

- **Command Reference section** (line ~203):
  - Add project create command details
  - Add project update command details
  - Add project archive command details
  - Document all flags and their values

- **Scripting & Automation section** (line ~587):
  - Add example for creating projects via API
  - Add example for updating projects programmatically

**2. master_api_ref.md** (if exists):

Update API reference with new mutations if this file exists.

**3. CONTRIBUTING.md:**

No changes needed - existing guidelines cover this feature.

**4. Command Help Text:**

Update in-code help strings:
- `cmd/issue.go` - issueCreateCmd and issueUpdateCmd Long descriptions
- `cmd/project.go` - Add Long descriptions for new commands

---

## UX/UI Considerations

**CLI-Specific User Experience:**

**Terminal Output Design:**

1. **Table Format Enhancements:**
   - Add "State" column to `project list` output
   - Add "Priority" column to `project list` output
   - Keep table width reasonable (max 120 characters)
   - Use abbreviated state names if needed (e.g., "Strtd" for "started")

2. **Success Messages:**
   - Project created: "✓ Project created: [PROJECT-NAME] (ID: [UUID])"
   - Project updated: "✓ Project updated: [PROJECT-NAME]"
   - Project archived: "✓ Project archived: [PROJECT-NAME]"
   - Issue project assigned: "✓ Issue [ISS-123] assigned to project [PROJECT-NAME]"

3. **Error Messages:**
   - Clear and actionable error messages
   - Include hints for resolution
   - Example: "Team 'INVALID' not found. Use 'linctl team list' to see available teams."

4. **Progress Indicators:**
   - Not needed (CLI operations are fast, < 2 seconds)
   - Only display final result

**Command Design Principles:**

1. **Consistency:**
   - Follow existing command patterns
   - Use same flag names across commands (e.g., `--project` everywhere)
   - Consistent output format selection (`--json`, `--plaintext`)

2. **Discoverability:**
   - Comprehensive help text via `--help`
   - Examples in Long description
   - Mention in README.md

3. **Error Prevention:**
   - Validate inputs before API calls
   - Clear required vs optional flags
   - Helpful error messages

4. **Efficiency:**
   - Support multi-field updates in single command
   - Minimize required typing (use sensible defaults)
   - Enable scripting with JSON output

**No Visual UI Changes:**
- This is a CLI tool, no graphical interface
- All output via terminal (STDOUT/STDERR)
- Colors via fatih/color library (existing)
- Tables via tablewriter library (existing)

**Accessibility Considerations (CLI):**
- Plain text output works with screen readers
- Color output respects terminal capabilities
- JSON output parseable by assistive tools
- Help text comprehensive and structured

---

## Testing Approach

**Conforming to Existing Test Standards:**

- **Test Framework:** Bash shell scripts (tests/smoke_test.sh)
- **Test File Naming:** Single file for all smoke tests
- **Test Organization:** Functions per command group
- **Assertion Style:** Exit code checking (0 = success, 1 = failure)
- **Coverage Requirements:** Read-only commands only (write commands have side effects)
- **Mocking/Stubbing:** No mocking - tests use real Linear API

**Test Strategy for This Feature:**

**1. Automated Smoke Tests:**
- No changes needed to `tests/smoke_test.sh`
- Existing tests for `project list` and `project get` cover display changes
- Write commands (create, update, archive) tested manually

**2. Manual Test Suite:**

Create `tests/manual_project_tests.sh` with test procedures:

```bash
#!/bin/bash
# Manual test procedures for project management feature

echo "=== Story 1: Issue-Project Assignment ==="
echo "1. Create issue with project:"
echo "   linctl issue create --title 'Test Issue' --team ENG --project PROJECT-UUID"
echo "   Expected: Issue created and assigned to project"
echo ""
echo "2. Update issue project:"
echo "   linctl issue update ISS-123 --project PROJECT-UUID"
echo "   Expected: Issue project updated"
echo ""
echo "3. Remove issue project:"
echo "   linctl issue update ISS-123 --project unassigned"
echo "   Expected: Project assignment removed"

echo ""
echo "=== Story 2: Project Creation & Archival ==="
echo "4. Create project:"
echo "   linctl project create --name 'Test Project' --team ENG"
echo "   Expected: Project created with default values"
echo ""
echo "5. Archive project:"
echo "   linctl project archive PROJECT-UUID"
echo "   Expected: Project archived successfully"

echo ""
echo "=== Story 3: Project Updates ==="
echo "6. Update project name:"
echo "   linctl project update PROJECT-UUID --name 'New Name'"
echo "   Expected: Name updated"
echo ""
echo "7. Update project state:"
echo "   linctl project update PROJECT-UUID --state started"
echo "   Expected: State updated"
echo ""
echo "8. Multi-field update:"
echo "   linctl project update PROJECT-UUID --state started --priority 1"
echo "   Expected: Both fields updated"

echo ""
echo "=== Enhanced Display ==="
echo "9. View project details:"
echo "   linctl project get PROJECT-UUID"
echo "   Expected: All fields displayed (state, priority, initiative, labels)"
echo ""
echo "10. List projects:"
echo "   linctl project list"
echo "   Expected: State and priority columns shown"
```

**3. Integration Testing Checklist:**
- [ ] Complete workflow: Create project → Assign to issue → Update project → Archive project
- [ ] Multi-team testing
- [ ] All output formats (table, JSON, plaintext)
- [ ] Verify changes in Linear web UI
- [ ] Test with existing projects
- [ ] Test with new projects

**4. Error Scenario Testing:**
- [ ] No authentication
- [ ] Invalid API key
- [ ] Invalid team key
- [ ] Invalid project UUID
- [ ] Missing required fields
- [ ] Invalid state values
- [ ] Invalid priority values
- [ ] Project not found
- [ ] Insufficient permissions
- [ ] Network timeout

---

## Deployment Strategy

### Deployment Steps

**Development → Production Process:**

1. **Development:**
   - Implement feature in feature branch
   - Test locally with `go run main.go`
   - Run `make fmt` and `make test`
   - Manual testing against test workspace

2. **Code Review:**
   - Create pull request to main branch
   - Review code changes
   - Verify all acceptance criteria met
   - Check README.md updates

3. **Merge to Main:**
   - Merge PR after approval
   - Ensure CI/CD pipeline passes
   - Main branch updated

4. **Create Release:**
   ```bash
   # Tag the release
   git tag v0.X.0 -a -m "v0.X.0: Add project management features"
   git push origin v0.X.0

   # Create GitHub release
   gh release create v0.X.0 \
     --title "linctl v0.X.0 - Project Management" \
     --notes "## New Features
   - Add project assignment to issues
   - Add project create command
   - Add project update command
   - Add project archive command
   - Enhanced project display

   ## Usage
   See README.md for examples.

   ## Breaking Changes
   None - fully backward compatible."
   ```

5. **Homebrew Tap Update (Automatic):**
   - GitHub Action (bump-tap.yml) triggers on release
   - Computes SHA256 of release tarball
   - Opens PR to homebrew-linctl tap
   - Merges after validation

6. **Verify Distribution:**
   ```bash
   brew update
   brew upgrade linctl
   linctl --version  # Should show new version
   linctl project create --help  # Verify new command available
   ```

### Rollback Plan

**If issues discovered post-release:**

1. **Identify Issue:**
   - Check issue tracker
   - Reproduce problem
   - Assess severity

2. **Quick Fix vs Rollback:**
   - **Minor issue:** Create hotfix branch, fix, release patch version
   - **Major issue:** Rollback to previous version

3. **Rollback Procedure:**
   ```bash
   # Revert to previous version in Homebrew tap
   cd homebrew-linctl
   git revert HEAD  # Revert the bump commit
   git push origin master

   # Users can downgrade
   brew uninstall linctl
   brew install linctl@0.X.Y  # Previous version
   ```

4. **Fix and Re-release:**
   - Fix issue in feature branch
   - Test thoroughly
   - Release new patched version
   - Update Homebrew tap

### Monitoring

**Post-Deployment Monitoring:**

1. **GitHub Issues:**
   - Monitor for bug reports
   - Respond to user questions
   - Track feature requests

2. **Homebrew Analytics (if enabled):**
   - Track install/upgrade counts
   - Monitor adoption rate

3. **API Usage:**
   - Linear API rate limits (5,000/hour per key)
   - No server-side monitoring needed (CLI tool)

4. **User Feedback:**
   - Monitor GitHub Discussions
   - Track issues mentioning project management
   - Collect feature enhancement requests

**Success Metrics:**

- Zero critical bugs reported
- Positive user feedback on GitHub
- Successful Homebrew tap update
- No rollback required
- Documentation complete and clear

**Support Plan:**

- Monitor GitHub issues daily for first week
- Respond to questions within 24 hours
- Create FAQ if common questions emerge
- Update documentation based on feedback

---

## Summary

**Tech-Spec Complete!**

**Feature:** Comprehensive Project Management for linctl

**Deliverables:**
- 3 User Stories (1.5h + 3.5h + 4h = 9h total)
- 4 files modified (cmd/issue.go, cmd/project.go, pkg/api/queries.go, pkg/api/client.go)
- README.md documentation updates
- Manual test procedures

**Key Capabilities:**
1. ✅ Issue-project assignment during create/update
2. ✅ Project CRUD operations (create, update, archive)
3. ✅ Multi-field project updates
4. ✅ Enhanced project display (state, priority, initiative, labels)

**Implementation Path:**
- Story 1: Issue-project assignment → 1.5 hours
- Story 2: Project create & archive → 3.5 hours
- Story 3: Project updates & display → 4 hours

**Next Steps:**
1. Run `create-story` workflow to generate user story markdown files
2. Begin implementation with Story 1
3. Test each story independently
4. Integrate and test complete workflow
5. Update documentation
6. Create PR and release

---

**Tech-Spec saved to:** `docs-bmad/tech-spec.md`
**Ready for story generation!**
