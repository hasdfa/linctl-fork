# Epic Technical Specification: Comprehensive Project Management

## Post-Review Follow-ups

- Standardize invalid project ID error surfaced by CLI to a clear message: "Project '<value>' not found". Consider mapping known GraphQL error cases in the API client (pkg/api/client.go) or at command error sites for create/update.
- Add basic UUID format validation for the `--project` flag (accept special `unassigned`).

- Include project name in archive success output (CLI). Options: return Project from `ArchiveProject` or fetch after archive. [Resolved 2025-11-06]
- Add tests for project create/archive flows and validation error paths (JSON/plaintext). [Resolved 2025-11-06]
- Verify/fix `GetTeam` GraphQL selector (key vs id) for team resolution. [Resolved 2025-11-06]
- Validate `--target-date` (YYYY-MM-DD) before API call. [Resolved 2025-11-06]

Date: 2025-11-06
Author: John
Epic ID: 1
Status: Draft

---

## Overview

This technical specification covers Epic 1: Comprehensive Project Management, which enables complete project lifecycle management from the linctl CLI. The epic extends linctl's existing Linear API integration to support project creation, updates, archival, and issue-project assignment operations. This eliminates the need for users to context-switch to Linear's web UI for project management tasks, supporting developers, project managers, and automation tools with comprehensive CLI access to project operations.

The epic builds upon linctl's established architecture—a layered CLI application using Cobra framework, Viper configuration, and GraphQL API integration—while adding new command groups and extending existing issue commands with project assignment capabilities.

## Objectives and Scope

**In Scope:**
- Issue-project assignment during issue creation and updates (`--project` flag)
- Project creation with required fields (name, team) and optional fields (state, priority, description, summary, initiative, labels)
- Project updates supporting multi-field modifications in single API calls
- Project archival operations
- Enhanced project display showing all key fields (state, priority, initiative, labels, description, summary)
- All operations support table/JSON/plaintext output formats consistent with existing commands
- Comprehensive validation and error handling with actionable messages
- README.md documentation with examples

**Out of Scope:**
- Project un-archival (rare use case, can be added later if needed)
- Advanced project features (milestones, roadmaps, custom templates)
- Bulk project operations
- Project members/permissions management
- Issue filtering by project (future enhancement)
- New output formats beyond existing table/JSON/plaintext

**Success Criteria:**
All 10 acceptance criteria from epics.md must be met, including backward compatibility, code quality standards (gofmt), and comprehensive documentation.

## System Architecture Alignment

This epic aligns with linctl's existing layered CLI architecture:

**UI Layer (cmd/):**
- Extends `cmd/issue.go` with `--project` flag for issueCreateCmd and issueUpdateCmd
- Extends `cmd/project.go` with new commands: `projectCreateCmd`, `projectUpdateCmd`, `projectArchiveCmd`
- Enhanced `projectGetCmd` and `projectListCmd` to display additional fields

**Business Logic Layer (cmd/):**
- Command handlers process flags and validate inputs before API calls
- Team key resolution to UUID (existing pattern)
- Project ID validation and "unassigned" special case handling
- Multi-field update logic using `cmd.Flags().Changed()` pattern

**API Client Layer (internal/api/):**
- New GraphQL mutations in `pkg/api/queries.go` (or equivalent location):
  - `issueCreate` and `issueUpdate` mutations with projectId field
  - `projectCreate` mutation with all required/optional fields
  - `projectUpdate` mutation with partial update support
  - `projectArchive` mutation
- Enhanced GraphQL queries to include new fields (description, shortSummary, state, priority, initiative, labels)

**Data Layer:**
- No changes to `~/.linctl.yaml` or `~/.linctl-auth.json` structure
- Uses existing authentication and configuration patterns

**Constraints:**
- Must maintain backward compatibility (no breaking changes)
- Must follow existing code conventions and patterns
- Must use Linear's GraphQL API schema (no custom extensions)
- Must respect Linear's rate limits (5,000 requests/hour)
- Must maintain stateless operation model (no local caching)

## Detailed Design

### Services and Modules

| Module | File Path | Responsibilities | Inputs | Outputs | Owner |
|--------|-----------|------------------|--------|---------|-------|
| **Issue Commands** | `cmd/issue.go` | Extended with `--project` flag for create/update commands | CLI flags, team key, project UUID | Issue data, success/error messages | UI Layer |
| **Project Commands** | `cmd/project.go` | New commands: create, update, archive; enhanced get/list display | CLI flags, project/team identifiers | Project data, success/error messages | UI Layer |
| **API Client - Mutations** | `internal/api/queries.go` (or similar) | GraphQL mutation execution for issue/project operations | Mutation input objects | GraphQL responses, errors | API Layer |
| **API Client - Queries** | `internal/api/queries.go` | Enhanced GraphQL queries to fetch project fields | Query filters, field selections | Project/Issue objects | API Layer |
| **Output Formatters** | `internal/output/` | Table/JSON/plaintext rendering with new fields | Structured data objects | Formatted output strings | Output Layer |
| **Validation Logic** | `cmd/project.go` | Input validation (state, priority, required fields) | User input strings | Validated values or error messages | Business Logic |
| **Team Resolution** | `cmd/project.go` (reuses existing pattern) | Resolve team key to UUID | Team key string | Team UUID or error | Business Logic |

### Data Models and Contracts

**Project Entity (Enhanced):**
```go
type Project struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    Description   string    `json:"description"`      // NEW
    ShortSummary  string    `json:"shortSummary"`     // NEW (called 'summary' in CLI)
    State         string    `json:"state"`            // NEW (planned|started|paused|completed|canceled)
    Priority      int       `json:"priority"`         // NEW (0-4)
    Progress      float64   `json:"progress"`         // EXISTING
    Team          Team      `json:"team"`             // EXISTING
    Initiative    *Initiative `json:"initiative"`     // NEW (nullable)
    Labels        []Label   `json:"labels"`           // NEW (array)
    CreatedAt     time.Time `json:"createdAt"`        // EXISTING
    UpdatedAt     time.Time `json:"updatedAt"`        // EXISTING
    ArchivedAt    *time.Time `json:"archivedAt"`      // NEW (nullable)
}

type Initiative struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type Label struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

**Issue Entity (Extended):**
```go
type Issue struct {
    ID         string   `json:"id"`
    Identifier string   `json:"identifier"`
    Title      string   `json:"title"`
    Project    *Project `json:"project"`  // NEW (nullable)
    // ... existing fields ...
}
```

**GraphQL Input Objects:**

```graphql
# Issue Create/Update Input (Extended)
input IssueUpdateInput {
  title: String
  description: String
  teamId: String
  assigneeId: String
  projectId: String      # NEW - can be null to unassign
  stateId: String
  priority: Int
  # ... other existing fields
}

# Project Create Input
input ProjectCreateInput {
  name: String!          # Required
  teamId: String!        # Required
  description: String
  shortSummary: String   # Called 'summary' in CLI
  state: String          # Enum: planned, started, paused, completed, canceled
  priority: Int          # 0-4
  initiativeId: String
  labelIds: [String!]
}

# Project Update Input
input ProjectUpdateInput {
  name: String
  description: String
  shortSummary: String
  state: String
  priority: Int
  initiativeId: String
  labelIds: [String!]
}

# Project Archive Input
input ProjectArchiveInput {
  # Linear API may only require project ID in mutation args
}
```

**Validation Rules:**
- **State**: Must be one of: `planned`, `started`, `paused`, `completed`, `canceled`
- **Priority**: Integer 0-4 (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)
- **Team Key**: Must exist in workspace (validated via team lookup)
- **Project UUID**: Must be valid UUID format and exist in workspace
- **Name**: Required, non-empty string for project creation
- **Labels**: Comma-separated string converted to array of label IDs (requires label lookup)

### APIs and Interfaces

**1. Issue Create/Update with Project Assignment**

```graphql
mutation IssueCreate($input: IssueCreateInput!) {
  issueCreate(input: $input) {
    success
    issue {
      id
      identifier
      title
      project {
        id
        name
      }
      team { id key name }
      assignee { id name email }
      state { id name }
      priority
      createdAt
      updatedAt
    }
  }
}

mutation IssueUpdate($id: String!, $input: IssueUpdateInput!) {
  issueUpdate(id: $id, input: $input) {
    success
    issue {
      id
      identifier
      title
      project {
        id
        name
      }
      # ... same fields as create
    }
  }
}
```

**CLI Interface:**
```bash
# Issue create with project
linctl issue create --title "Bug fix" --team ENG --project PROJECT_UUID

# Issue update with project
linctl issue update ISS-123 --project PROJECT_UUID

# Issue unassign from project
linctl issue update ISS-123 --project unassigned
```

**2. Project Creation**

```graphql
mutation ProjectCreate($input: ProjectCreateInput!) {
  projectCreate(input: $input) {
    success
    project {
      id
      name
      description
      shortSummary
      state
      priority
      team { id key name }
      initiative { id name }
      labels { nodes { id name } }
      progress
      createdAt
      updatedAt
    }
  }
}
```

**CLI Interface:**
```bash
# Minimal create (required fields only)
linctl project create --name "Q1 Backend" --team ENG

# Full create with all optional fields
linctl project create \
  --name "Q1 Backend" \
  --team ENG \
  --state started \
  --priority 1 \
  --description "Full project description" \
  --summary "Short summary" \
  --initiative INIT_UUID \
  --label "urgent,backend"
```

**3. Project Update**

```graphql
mutation ProjectUpdate($id: String!, $input: ProjectUpdateInput!) {
  projectUpdate(id: $id, input: $input) {
    success
    project {
      id
      name
      description
      shortSummary
      state
      priority
      initiative { id name }
      labels { nodes { id name } }
      updatedAt
    }
  }
}
```

**CLI Interface:**
```bash
# Single field update
linctl project update PROJECT_UUID --name "New Name"

# Multi-field update
linctl project update PROJECT_UUID \
  --state started \
  --priority 2 \
  --description "Updated description" \
  --summary "Updated summary" \
  --label "backend,api"
```

**4. Project Archive**

```graphql
mutation ProjectArchive($id: String!) {
  projectArchive(id: $id) {
    success
    project {
      id
      name
      archivedAt
    }
  }
}
```

**CLI Interface:**
```bash
linctl project archive PROJECT_UUID
```

**5. Enhanced Project Get/List**

```graphql
query Project($id: String!) {
  project(id: $id) {
    id
    name
    description
    shortSummary
    state
    priority
    progress
    team { id key name }
    initiative { id name }
    labels { nodes { id name } }
    createdAt
    updatedAt
    archivedAt
  }
}

query Projects($filter: ProjectFilter, $first: Int) {
  projects(filter: $filter, first: $first) {
    nodes {
      id
      name
      state
      priority
      progress
      team { key name }
      initiative { name }
      labels { nodes { name } }
      updatedAt
    }
  }
}
```

**Error Responses:**
- `400`: Invalid input (validation failure) → "Invalid state. Must be one of: planned, started, paused, completed, canceled"
- `404`: Resource not found → "Project 'UUID' not found" or "Team 'KEY' not found"
- `401`: Unauthorized → "Authentication required. Run 'linctl auth' to login"
- `429`: Rate limit → "Rate limit exceeded. Try again later"
- `500`: Server error → "Linear API error: <message>"

### Workflows and Sequencing

**Workflow 1: Issue Creation with Project Assignment**

```
Actor: User
1. User executes: linctl issue create --title "Bug" --team ENG --project PROJ_UUID
2. CLI parses flags and validates inputs
3. CLI resolves team key "ENG" to team UUID (existing pattern)
4. CLI validates project UUID format (basic check)
5. CLI constructs IssueCreateInput with projectId field
6. CLI sends GraphQL mutation to Linear API
7. API validates project exists and user has access
8. API creates issue and assigns to project
9. API returns issue object with project relationship
10. CLI formats output (table/JSON/plaintext) and displays
```

**Workflow 2: Project Creation**

```
Actor: Project Manager
1. User executes: linctl project create --name "Q1" --team ENG --state started --priority 1
2. CLI validates required fields (name, team) present
3. CLI validates optional fields (state, priority) if provided
4. CLI resolves team key to team UUID
5. CLI constructs ProjectCreateInput object
6. CLI sends GraphQL mutation to Linear API
7. API validates all inputs and creates project
8. API returns project object with all fields
9. CLI formats and displays success message with project details
```

**Workflow 3: Multi-Field Project Update**

```
Actor: Project Manager
1. User executes: linctl project update PROJ_UUID --state started --priority 2 --label "urgent"
2. CLI uses cmd.Flags().Changed() to detect which flags were explicitly set
3. CLI validates changed fields (state enum, priority range)
4. CLI constructs ProjectUpdateInput with only changed fields
5. If labels provided, CLI resolves label names to label IDs (requires label lookup)
6. CLI sends GraphQL mutation to Linear API
7. API applies partial update (only modified fields)
8. API returns updated project object
9. CLI displays updated fields to user
```

**Workflow 4: Project Archival**

```
Actor: Project Manager
1. User executes: linctl project archive PROJ_UUID
2. CLI validates project UUID format
3. CLI sends GraphQL mutation to Linear API
4. API archives project (sets archivedAt timestamp)
5. API returns success response
6. CLI displays success message: "Project 'Name' archived successfully"
```

**Sequence Diagram: Issue-Project Assignment Flow**

```
User -> CLI: linctl issue create --project PROJ_UUID
CLI -> Validator: Validate inputs
Validator -> CLI: OK
CLI -> TeamResolver: Resolve team key to UUID
TeamResolver -> LinearAPI: Query teams
LinearAPI -> TeamResolver: Team UUID
TeamResolver -> CLI: Team UUID
CLI -> LinearAPI: IssueCreate mutation (with projectId)
LinearAPI -> LinearAPI: Validate project exists
LinearAPI -> CLI: Issue object (with project)
CLI -> OutputFormatter: Format response
OutputFormatter -> User: Display issue details
```

**Error Handling Flows:**

- **Invalid Project UUID**: Validate before API call, fail with "Project 'UUID' not found" if API returns 404
- **Missing Required Fields**: Fail locally with clear message before API call
- **Invalid State/Priority**: Validate locally against allowed values before API call
- **Team Not Found**: Fail during team resolution step with "Team 'KEY' not found. Use 'linctl team list' to see available teams."
- **"unassigned" Special Case**: Set projectId to nil when user passes `--project unassigned`

## Non-Functional Requirements

### Performance

**Target Metrics:**
- **Command Execution Latency**: < 2 seconds for create/update operations (excluding network time)
- **Network Request Time**: < 1 second for GraphQL API calls under normal conditions
- **Startup Time**: < 100ms for CLI initialization (no change from existing)
- **Memory Usage**: < 50MB per command execution (typical Go CLI pattern)

**Performance Requirements:**
- Reuse existing HTTP client connection pooling patterns (already implemented in linctl)
- Single GraphQL request per command operation (no N+1 query problems)
- Efficient team/label resolution with minimal API round-trips
- No local caching required (stateless operation model)

**Optimization Strategies:**
- Batch label resolution if multiple labels provided (single query)
- Validate inputs locally before API calls to avoid unnecessary network requests
- Use GraphQL field selection to fetch only required fields
- Follow existing pagination patterns for list commands (default 50 items)

**Performance Baselines (from architecture.md):**
- Linear API rate limit: 5,000 requests/hour for Personal API Keys
- Default timeout: 30 seconds (existing configuration)
- Connection pooling: Enabled via Go's http.Client (existing)

### Security

**Authentication & Authorization:**
- Reuse existing Personal API Key storage in `~/.linctl-auth.json`
- File permissions: 0600 (read/write for owner only)
- API key transmitted via Bearer token in Authorization header over HTTPS
- No password storage or OAuth flows required
- User must have appropriate Linear workspace permissions for project operations

**Data Handling:**
- All API communication over HTTPS (TLS 1.2+)
- No sensitive data logged to stdout/stderr
- Project/issue data displayed according to user's Linear permissions
- API keys never echoed or included in error messages

**Input Validation:**
- Sanitize all user inputs before API transmission
- Validate UUID formats using regex pattern matching
- Validate enum values (state, priority) against allowed lists
- Protect against command injection in shell contexts (use existing patterns)

**Threat Considerations:**
- **Man-in-the-Middle**: Mitigated by HTTPS enforcement
- **API Key Exposure**: Mitigated by file permissions and no logging
- **Injection Attacks**: Mitigated by input validation and GraphQL parameterization
- **Rate Limit DoS**: Mitigated by Linear's built-in rate limiting

**Security Requirements from Architecture:**
- Follow existing security patterns in linctl codebase
- No changes to authentication mechanism
- Maintain stateless operation (no session management)
- Use GitHub Secrets for CI/CD automation tokens

### Reliability/Availability

**Availability Targets:**
- **CLI Availability**: 100% (local tool, no downtime)
- **API Dependency**: Dependent on Linear API availability (SLA managed by Linear)
- **Graceful Degradation**: Provide clear error messages when API unavailable

**Error Recovery:**
- Retry logic for transient network failures (existing pattern in linctl)
- Timeout handling with configurable limits (30s default)
- Clear error messages for user remediation
- Exit codes for scripting: 0 (success), 1 (error)

**Data Consistency:**
- Stateless operations ensure no local state corruption
- All operations idempotent where possible (update operations)
- No local caching means data always fresh from API
- Transaction management handled by Linear API

**Failure Scenarios & Handling:**
- **Network Timeout**: Retry with exponential backoff, then fail with message
- **Invalid API Key**: Prompt user to run `linctl auth`
- **Rate Limit Exceeded**: Display clear message with retry guidance
- **Partial Update Failure**: Display which fields failed, suggest retry
- **GraphQL Errors**: Parse and display specific error messages from API

**Backward Compatibility:**
- New `--project` flag on issue commands is optional (no breaking change)
- New project commands don't affect existing commands
- Output format flags work consistently across all commands
- Existing configuration files remain compatible

### Observability

**Logging Requirements:**
- Standard output for successful operations (formatted according to output flag)
- Standard error for error messages and warnings
- No verbose/debug logging in production builds (keep CLI simple)
- Error messages include actionable remediation steps

**Required Log Signals:**
- **Command Execution**: Log command invocation (optional verbose mode)
- **API Requests**: Log GraphQL operation name and status (optional verbose mode)
- **Validation Failures**: Log which validation failed and why (always)
- **API Errors**: Log HTTP status code and GraphQL error details (always)

**Metrics (for future consideration):**
- Command execution count by type
- API response time distribution
- Error rate by error type
- Most-used commands

**Tracing:**
- No distributed tracing required (single-process CLI)
- Request ID from Linear API included in error messages if available
- Stack traces for unexpected errors (Go panic recovery)

**Monitoring & Alerting:**
- Not applicable for client-side CLI tool
- Users report issues via GitHub Issues
- Smoke tests provide basic regression detection

**Debugging Support:**
- Clear error messages with context
- Suggest next steps for common errors
- Include Linear API error details when available
- Consider adding `--verbose` flag for debugging (future enhancement)

## Dependencies and Integrations

### External Dependencies

**Linear GraphQL API**
- **Version**: Current (no versioning, rolling API)
- **Endpoint**: `https://api.linear.app/graphql`
- **Purpose**: All project and issue data operations
- **Critical Mutations Required**:
  - `issueCreate` - Must support `projectId` field in input
  - `issueUpdate` - Must support `projectId` field in input
  - `projectCreate` - Full project creation with all fields
  - `projectUpdate` - Partial project updates
  - `projectArchive` - Archive projects
- **Critical Queries Required**:
  - `project(id)` - Fetch single project with all fields
  - `projects(filter)` - List projects with filtering
  - `teams` - Team lookup for key-to-UUID resolution
  - `labels` - Label lookup for name-to-ID resolution
- **Rate Limits**: 5,000 requests/hour (Personal API Keys)
- **Documentation**: https://developers.linear.app/

### Go Module Dependencies (from go.mod)

**Direct Dependencies (No Changes):**
- `github.com/spf13/cobra` v1.8.0 - CLI command framework
- `github.com/spf13/viper` v1.18.2 - Configuration management
- `github.com/olekukonko/tablewriter` v0.0.5 - Table output formatting
- `github.com/fatih/color` v1.16.0 - Terminal color support

**Indirect Dependencies (No Changes):**
- Standard Go library dependencies remain unchanged
- No new third-party dependencies required

**Go Version:**
- `go 1.23.0` (minimum)
- `toolchain go1.24.5` (current)

### Internal Dependencies

**Existing Modules (Reused):**
- `cmd/root.go` - Root command setup and global flags
- `internal/api/` - GraphQL client, HTTP configuration, request handling
- `internal/config/` - Config file and auth token management
- `internal/output/` - Table, JSON, plaintext formatters

**Existing Patterns (Reused):**
- Team key to UUID resolution (used in issue commands)
- Flag parsing with Cobra
- Output format selection via `--json`/`--plaintext` flags
- Error handling and display
- GraphQL query construction

### Integration Points

**1. Linear Workspace Integration**
- **Authentication**: Personal API Key (existing)
- **Permissions Required**:
  - Project read/write access
  - Issue read/write access
  - Team read access (for team resolution)
  - Label read access (for label resolution)
- **Workspace Configuration**: Must have teams and optionally labels/initiatives configured

**2. Cobra CLI Framework Integration**
- New command registration under `projectCmd` and existing `issueCmd`
- Flag definitions using `cmd.Flags().StringP()` pattern
- Help text generation (automatic via Cobra)
- Command execution flow (existing pattern)

**3. Output Formatter Integration**
- Extend table formatter to include new project fields (state, priority, description, summary)
- JSON output automatically includes all fields (no changes needed)
- Plaintext formatter follows existing patterns

**4. Configuration File Integration**
- No changes to `~/.linctl.yaml` schema
- No changes to `~/.linctl-auth.json` schema
- Reuse existing Viper configuration loading

### Integration Testing Requirements

**API Integration Tests:**
- Test against real Linear API (development workspace recommended)
- Validate GraphQL mutations work as expected
- Verify error handling for invalid inputs
- Test rate limiting behavior

**CLI Integration Tests:**
- Smoke tests for all new commands (extend `tests/smoke_test.sh`)
- Test flag combinations and validation
- Test output format consistency across commands
- Test backward compatibility (existing commands unchanged)

**End-to-End Workflows:**
- Create project → Create issue with project → Update project → Archive project
- Test "unassigned" special case for issue-project relationship
- Test multi-field project updates
- Verify table/JSON/plaintext output for all commands

### Version Constraints

**Minimum Versions:**
- Go: 1.23.0
- Linear API: Current (no minimum version, API is rolling)
- macOS: 10.15+ (for Homebrew distribution)
- Linux: Any modern distribution with glibc 2.28+

**Known Incompatibilities:**
- None (all dependencies are stable and widely used)

### Third-Party Service Dependencies

**GitHub (for CI/CD):**
- GitHub Actions for automated Homebrew tap updates
- GitHub Releases for version distribution
- Secrets: `HOMEBREW_TAP_TOKEN` (existing)

**Homebrew (for Distribution):**
- Custom tap: `dorkitude/homebrew-linctl` (existing)
- Formula auto-update workflow (existing)
- No changes required to distribution mechanism

## Acceptance Criteria (Authoritative)

### Epic-Level Success Criteria

**AC-E1**: Issue-Project Integration - Users can assign projects to issues during `issue create` and `issue update` operations using `--project` flag

**AC-E2**: Project Creation - Users can create new projects via `project create` command with team assignment and optional configuration

**AC-E3**: Project Updates - Users can update multiple project fields in a single command using `project update`

**AC-E4**: Project Archival - Users can archive projects via `project archive` command

**AC-E5**: Enhanced Display - `project get` and `project list` commands show state, priority, initiative, and labels

**AC-E6**: Output Format Support - All commands support `--json` and `--plaintext` flags for automation

**AC-E7**: Error Handling - Clear, actionable error messages for validation failures and API errors

**AC-E8**: Documentation - README.md includes examples for all new commands and flags

**AC-E9**: Backward Compatibility - No breaking changes to existing commands

**AC-E10**: Code Quality - All code follows existing linctl conventions and passes `gofmt`

### Story 1.1: Issue-Project Assignment

**AC-1.1.1**: Given a valid project UUID, when I run `linctl issue create --title "Test" --team ENG --project PROJECT-UUID`, then the issue is created and assigned to the specified project, and the output shows the project assignment.

**AC-1.1.2**: Given an existing issue and valid project UUID, when I run `linctl issue update ISS-123 --project PROJECT-UUID`, then the issue's project is updated to the new project.

**AC-1.1.3**: Given an issue with an existing project assignment, when I run `linctl issue update ISS-123 --project unassigned`, then the project assignment is removed from the issue.

**AC-1.1.4**: Given an invalid project UUID, when I attempt to assign it to an issue, then the command fails with a clear error message: "Project 'INVALID-UUID' not found".

**AC-1.1.5**: JSON output includes the project field showing project ID and name when `--json` flag is used.

### Story 1.2: Project Creation & Archival

**AC-1.2.1**: Given a project name and team key, when I run `linctl project create --name "Q1 Backend" --team ENG`, then a new project is created in Linear with default values (state: planned, priority: 0).

**AC-1.2.2**: Given project creation with optional fields, when I run `linctl project create --name "Test" --team ENG --state started --priority 1 --description "Test project"`, then the project is created with all specified field values.

**AC-1.2.3**: Given a valid project UUID, when I run `linctl project archive PROJECT-UUID`, then the project is archived in Linear and a success message is displayed.

**AC-1.2.4**: Given missing required fields, when I attempt project creation without --name or --team, then the command fails with error: "Both --name and --team are required".

**AC-1.2.5**: Given an invalid team key, when I attempt project creation, then the command fails with error: "Team 'INVALID' not found. Use 'linctl team list' to see available teams."

**AC-1.2.6**: Given invalid state or priority values, when I attempt project creation, then validation fails before API call with clear error messages.

**AC-1.2.7**: All commands support --json, --plaintext output formats and follow existing linctl formatting conventions.

### Story 1.3: Project Updates & Enhanced Display

**AC-1.3.1**: Given a project UUID and updated field, when I run `linctl project update PROJECT-UUID --name "New Name"`, then the project name is updated in Linear.

**AC-1.3.2**: Given a project UUID, when I run `linctl project update PROJECT-UUID --state started`, then the project state is updated and validated against allowed values (planned, started, paused, completed, canceled).

**AC-1.3.3**: Given multiple field updates, when I run `linctl project update PROJECT-UUID --state started --priority 1 --label "urgent,backend"`, then all specified fields are updated in a single API call.

**AC-1.3.3.1**: Given a project UUID, when I run `linctl project update PROJECT-UUID --description "Full description"`, then the project description is updated.

**AC-1.3.3.2**: Given a project UUID, when I run `linctl project update PROJECT-UUID --summary "Short summary"`, then the project shortSummary field is updated.

**AC-1.3.4**: Given no field flags provided, when I run `linctl project update PROJECT-UUID`, then the command fails with error: "At least one field to update is required".

**AC-1.3.5**: Given invalid state value, when I attempt update, then validation fails with error: "Invalid state. Must be one of: planned, started, paused, completed, canceled".

**AC-1.3.6**: Given invalid priority value, when I attempt update, then validation fails with error: "Priority must be between 0 and 4".

**AC-1.3.7**: Given a project UUID, when I run `linctl project get PROJECT-UUID`, then the output displays all fields including description, shortSummary, state, priority, initiative, and labels.

**AC-1.3.8**: When I run `linctl project list`, then the table output includes State and Priority columns for each project.

**AC-1.3.9**: All commands support --json and --plaintext output formats with complete field data.

## Traceability Mapping

| AC ID | Spec Section(s) | Component(s)/API(s) | Test Idea |
|-------|----------------|---------------------|-----------|
| **AC-E1** | APIs & Interfaces #1, Workflows #1 | `cmd/issue.go`, `issueCreate/issueUpdate` mutations | Create/update issues with valid project UUIDs; verify project assignment appears in output |
| **AC-E2** | APIs & Interfaces #2, Workflows #2 | `cmd/project.go`, `projectCreate` mutation | Create projects with minimal and full field sets; verify all fields stored correctly |
| **AC-E3** | APIs & Interfaces #3, Workflows #3 | `cmd/project.go`, `projectUpdate` mutation | Update single and multiple fields; verify partial updates work correctly |
| **AC-E4** | APIs & Interfaces #4, Workflows #4 | `cmd/project.go`, `projectArchive` mutation | Archive project; verify archivedAt timestamp set |
| **AC-E5** | APIs & Interfaces #5, Data Models | `cmd/project.go`, `project` query | Run project get/list; verify new fields displayed in table |
| **AC-E6** | Output Formatters | `internal/output/` | Test all commands with --json and --plaintext flags; verify consistent formatting |
| **AC-E7** | Error Handling, Validation Logic | `cmd/issue.go`, `cmd/project.go` | Test invalid inputs (bad state, priority, team); verify clear error messages |
| **AC-E8** | Documentation | `README.md` | Manual review of documentation; verify examples present and accurate |
| **AC-E9** | System Architecture | All existing commands | Run existing commands unchanged; verify no breaking changes |
| **AC-E10** | Code Quality | All Go files | Run `gofmt` on all files; verify no formatting issues |
| **AC-1.1.1** | APIs & Interfaces #1 | `issueCreate` mutation with projectId | Create issue with --project flag; verify issue created with project relationship |
| **AC-1.1.2** | APIs & Interfaces #1 | `issueUpdate` mutation with projectId | Update issue with --project flag; verify project changed |
| **AC-1.1.3** | Workflows #1, Error Handling | `issueUpdate` with projectId=nil | Update issue with --project unassigned; verify project removed |
| **AC-1.1.4** | Error Handling | Validation logic | Attempt issue create with invalid project UUID; verify error message |
| **AC-1.1.5** | Output Formatters | JSON formatter | Create/update issue with --json flag; verify project field in JSON |
| **AC-1.2.1** | APIs & Interfaces #2 | `projectCreate` mutation | Create project with minimal fields; verify defaults applied |
| **AC-1.2.2** | APIs & Interfaces #2 | `projectCreate` mutation | Create project with all optional fields; verify all values set |
| **AC-1.2.3** | APIs & Interfaces #4 | `projectArchive` mutation | Archive project; verify success message displayed |
| **AC-1.2.4** | Validation Logic | Flag parsing in `cmd/project.go` | Attempt project create without --name or --team; verify error |
| **AC-1.2.5** | Team Resolution | Team lookup logic | Attempt project create with invalid team key; verify error message |
| **AC-1.2.6** | Validation Logic | State/priority validation | Attempt project create with invalid state/priority; verify validation error |
| **AC-1.2.7** | Output Formatters | All formatters | Test project create/archive with --json/--plaintext; verify formatting |
| **AC-1.3.1** | APIs & Interfaces #3 | `projectUpdate` mutation | Update project name; verify change persisted |
| **AC-1.3.2** | Validation Logic, APIs #3 | State validation, `projectUpdate` | Update project state; verify enum validation works |
| **AC-1.3.3** | Workflows #3, APIs #3 | Multi-field update logic | Update multiple fields at once; verify all changed |
| **AC-1.3.3.1** | APIs & Interfaces #3 | `projectUpdate` mutation | Update project description; verify change persisted |
| **AC-1.3.3.2** | APIs & Interfaces #3 | `projectUpdate` mutation | Update project summary; verify change persisted |
| **AC-1.3.4** | Validation Logic | Flag detection with `cmd.Flags().Changed()` | Run project update with no flags; verify error message |
| **AC-1.3.5** | Validation Logic | State enum validation | Update with invalid state; verify error lists valid states |
| **AC-1.3.6** | Validation Logic | Priority range validation | Update with priority 5; verify error shows valid range |
| **AC-1.3.7** | APIs & Interfaces #5, Data Models | Enhanced `project` query | Get project by UUID; verify all fields displayed |
| **AC-1.3.8** | APIs & Interfaces #5, Output Formatters | Enhanced `projects` query, table formatter | List projects; verify State and Priority columns present |
| **AC-1.3.9** | Output Formatters | All formatters | Test project get/list/update with --json/--plaintext; verify complete data |

## Risks, Assumptions, Open Questions

### Risks

**RISK-1**: Linear GraphQL API Schema Changes
- **Impact**: High - Breaking changes to GraphQL schema could break all project operations
- **Probability**: Low - Linear maintains backward compatibility
- **Mitigation**: Monitor Linear API changelog; test against Linear API regularly; maintain integration tests
- **Contingency**: Quick patch release if breaking changes detected

**RISK-2**: Label/Initiative Resolution Performance
- **Impact**: Medium - Label name-to-ID resolution requires additional API calls
- **Probability**: Medium - Multi-label updates may be slow
- **Mitigation**: Batch label lookups in single query; consider caching label mappings (future)
- **Contingency**: Document performance characteristics; recommend using label IDs directly for automation

**RISK-3**: Rate Limiting on Heavy Usage
- **Impact**: Medium - Bulk project operations could hit 5,000 requests/hour limit
- **Probability**: Low - Most users perform few operations per hour
- **Mitigation**: Document rate limits; implement exponential backoff on 429 responses
- **Contingency**: Add rate limit awareness and queuing (future enhancement)

**RISK-4**: Incomplete GraphQL Mutation Support
- **Impact**: High - If Linear API doesn't support all fields in mutations, features incomplete
- **Probability**: Low - Linear API documentation shows support for these fields
- **Mitigation**: Validate all GraphQL mutations against Linear API before implementation
- **Contingency**: Document unsupported fields; implement workarounds or defer features

**RISK-5**: Breaking Existing Commands
- **Impact**: Critical - Breaking changes would affect all existing linctl users
- **Probability**: Low - Changes are additive (new flags, new commands)
- **Mitigation**: Comprehensive testing of existing commands; smoke test suite; manual testing
- **Contingency**: Quick rollback; patch release

### Assumptions

**ASSUMPTION-1**: Linear API supports all required GraphQL mutations
- **Validation**: Review Linear API documentation; test mutations in development workspace
- **If False**: Defer unsupported features; document limitations

**ASSUMPTION-2**: Users have appropriate Linear workspace permissions
- **Validation**: Document required permissions; test with restricted user accounts
- **If False**: Clear error messages guide users to request permissions

**ASSUMPTION-3**: Project UUIDs are stable and permanent
- **Validation**: Review Linear API documentation on identifier stability
- **If False**: Handle ID changes gracefully; document ID lifecycle

**ASSUMPTION-4**: Team keys and label names are unique within workspace
- **Validation**: Test with duplicate names in development workspace
- **If False**: Use UUIDs directly; enhance resolution logic to handle ambiguity

**ASSUMPTION-5**: Existing output formatters can handle new fields without changes
- **Validation**: Test table/JSON/plaintext formatters with enhanced data models
- **If False**: Update formatters to handle new field types (arrays, nested objects)

**ASSUMPTION-6**: No new Go dependencies required
- **Validation**: Review existing codebase patterns for validation and GraphQL handling
- **If False**: Evaluate and add minimal dependencies; update go.mod

### Open Questions

**QUESTION-1**: Should label resolution support fuzzy matching or exact match only?
- **Context**: Users may type "back-end" when label is "backend"
- **Decision Needed**: Before Story 1.3 implementation
- **Recommendation**: Start with exact match; add fuzzy matching if user feedback requests it

**QUESTION-2**: Should project archive command require confirmation prompt?
- **Context**: Archive is reversible via Linear UI but not via linctl
- **Decision Needed**: Before Story 1.2 implementation
- **Recommendation**: No confirmation (consistent with existing linctl patterns); document reversibility

**QUESTION-3**: How should multi-word labels be handled in CLI?
- **Context**: Labels with spaces need special handling (quotes or encoding)
- **Decision Needed**: Before Story 1.3 implementation
- **Recommendation**: Comma-separated list, quote labels with spaces: `--label "bug fix,urgent"`

**QUESTION-4**: Should initiative assignment support initiative names or UUIDs only?
- **Context**: Initiative name-to-ID resolution adds complexity
- **Decision Needed**: Before Story 1.3 implementation
- **Recommendation**: UUID only for MVP; add name resolution in future if requested

**QUESTION-5**: Should `project list` include archived projects by default?
- **Context**: Existing linctl patterns don't filter archived items by default
- **Decision Needed**: Before Story 1.3 implementation
- **Recommendation**: Follow existing patterns; add `--include-archived` flag if needed

## Test Strategy Summary

### Test Layers

**1. Unit Tests (Go test)**
- **Scope**: Input validation functions, team resolution logic, flag parsing
- **Coverage Target**: 80%+ for new functions
- **Location**: `*_test.go` files alongside implementation
- **Key Tests**:
  - State enum validation (valid and invalid values)
  - Priority range validation (0-4, reject 5+)
  - UUID format validation
  - "unassigned" special case handling
  - Flag change detection logic

**2. Integration Tests (Smoke Tests)**
- **Scope**: End-to-end command execution against real Linear API
- **Coverage**: All new commands and flags
- **Location**: `tests/smoke_test.sh` (extend existing suite)
- **Key Tests**:
  - `linctl issue create --project UUID` (success case)
  - `linctl project create --name "Test" --team ENG` (minimal fields)
  - `linctl project create` with all optional fields (full fields)
  - `linctl project update UUID --state started` (single field)
  - `linctl project update UUID --state started --priority 1` (multi-field)
  - `linctl project archive UUID` (archive operation)
  - `linctl project get UUID` (enhanced display)
  - `linctl project list` (enhanced display with new columns)
  - All commands with `--json` and `--plaintext` flags

**3. Manual Testing**
- **Scope**: Complex workflows, edge cases, error scenarios
- **Environment**: Development Linear workspace
- **Key Scenarios**:
  - Complete workflow: Create project → Create issue with project → Update project → Archive project
  - Error cases: Invalid team key, invalid project UUID, missing required fields
  - Output validation: Verify table formatting, JSON structure, plaintext readability
  - Backward compatibility: Run existing commands to ensure no regression

**4. Acceptance Testing**
- **Scope**: Validate all acceptance criteria from epics.md
- **Method**: Manual test against each AC
- **Documentation**: Test results recorded in Linear issue comments
- **Sign-off**: All ACs must pass before story marked "done"

### Test Data Requirements

**Test Workspace Setup:**
- Linear development workspace with test data
- Multiple teams (e.g., ENG, PRODUCT, DESIGN)
- Existing projects with various states and priorities
- Sample issues for project assignment testing
- Labels and initiatives configured

**Test Data:**
- Valid project UUIDs (from development workspace)
- Invalid project UUIDs (for error testing)
- Valid and invalid team keys
- Various state values (valid and invalid)
- Various priority values (0-4 and invalid)

### Testing Tools

**Existing Tools:**
- `tests/smoke_test.sh` - Automated smoke testing
- `go test` - Unit testing framework
- `gofmt` - Code formatting validation
- Linear development workspace - Integration testing

**Manual Testing:**
- Terminal for command execution
- Linear web UI for verification
- JSON validation tools (jq, etc.)

### Test Execution Plan

**During Development (Story 1.1, 1.2, 1.3):**
1. Write unit tests for validation logic
2. Test commands manually against Linear API
3. Verify output formats (table, JSON, plaintext)
4. Test error scenarios

**Before Story Completion:**
1. Run all unit tests (`go test ./...`)
2. Run extended smoke test suite
3. Execute acceptance test checklist
4. Manual regression testing of existing commands
5. Code formatting check (`gofmt`)

**Before Epic Completion (Day 7):**
1. Complete integration testing across all stories
2. End-to-end workflow validation
3. Performance testing (command latency)
4. Documentation review and examples testing
5. Backward compatibility verification

### Coverage of Acceptance Criteria

| Test Layer | AC Coverage |
|------------|-------------|
| **Unit Tests** | AC-E7 (validation), AC-E10 (code quality) |
| **Smoke Tests** | AC-E1, AC-E2, AC-E3, AC-E4, AC-E5, AC-E6 (all commands work) |
| **Manual Tests** | All story-level ACs (1.1.1 - 1.3.9) |
| **Documentation Review** | AC-E8 (documentation complete) |
| **Regression Tests** | AC-E9 (backward compatibility) |

### Edge Cases and Special Scenarios

**Edge Case Tests:**
- Empty strings for project name
- Very long project names (>255 characters)
- Special characters in project descriptions
- Multiple label assignment (single, multiple, none)
- Project with no initiative assigned (null handling)
- Archived project operations (should fail gracefully)
- Concurrent updates (last-write-wins via Linear API)
- Network failures and timeouts
- Rate limit scenarios (429 responses)

### Rollback Testing

**Pre-release Verification:**
- Install previous version of linctl
- Verify all existing commands work
- Install new version
- Verify all commands (old and new) work
- Rollback to previous version
- Verify existing commands still work (no data corruption)

### Success Criteria for Testing

**All tests must pass before release:**
- ✅ All unit tests pass (`go test ./...`)
- ✅ All smoke tests pass (extended suite)
- ✅ All acceptance criteria validated manually
- ✅ No regressions in existing commands
- ✅ Code formatting passes (`gofmt`)
- ✅ Documentation accurate and complete
- ✅ README examples tested and working
