# linctl - Epic Breakdown

**Date:** 2025-11-06
**Project Level:** 1 (Coherent Feature)

---

## Epic 1: Comprehensive Project Management

**Slug:** project-management

### Goal

Enable complete project lifecycle management from linctl CLI, eliminating the need to context-switch to Linear's web UI for project operations. Provide developers, project managers, and automation tools with comprehensive CLI access to project creation, updates, archival, and issue-project assignment capabilities.

### Scope

**In Scope:**
- Issue-project assignment during issue creation and updates
- Project creation with required and optional fields
- Project updates with multi-field support (name, state, priority, initiative, labels)
- Project archival operations
- Enhanced project display showing all key fields (state, priority, initiative, labels)
- All operations support table/JSON/plaintext output formats

**Out of Scope:**
- Project un-archival (rare use case)
- Advanced project features (milestones, roadmaps, templates)
- Bulk project operations
- Project members management
- Issue filtering by project (future enhancement)

### Success Criteria

1. ✅ **Issue-Project Integration:** Users can assign projects to issues during `issue create` and `issue update` operations using `--project` flag
2. ✅ **Project Creation:** Users can create new projects via `project create` command with team assignment and optional configuration
3. ✅ **Project Updates:** Users can update multiple project fields in a single command using `project update`
4. ✅ **Project Archival:** Users can archive projects via `project archive` command
5. ✅ **Enhanced Display:** `project get` and `project list` commands show state, priority, initiative, and labels
6. ✅ **Output Format Support:** All commands support `--json` and `--plaintext` flags for automation
7. ✅ **Error Handling:** Clear, actionable error messages for validation failures and API errors
8. ✅ **Documentation:** README.md includes examples for all new commands and flags
9. ✅ **Backward Compatibility:** No breaking changes to existing commands
10. ✅ **Code Quality:** All code follows existing linctl conventions and passes `gofmt`

### Dependencies

**External:**
- Linear GraphQL API (mutations: issueCreate, issueUpdate, projectCreate, projectUpdate, projectArchive)
- Linear workspace with appropriate permissions

**Internal:**
- Existing linctl infrastructure (auth, API client, output formatters)
- cmd/issue.go (extend with --project flag)
- cmd/project.go (extend with create/update/archive commands)
- pkg/api/queries.go (add GraphQL mutations)

**Framework:**
- Go 1.23.0+
- Cobra v1.8.0 (CLI framework)
- Viper v1.18.2 (configuration)
- tablewriter v0.0.5 (output formatting)

---

## Story Map - Epic 1

```
Epic: Comprehensive Project Management
├── Story 1.1: Issue-Project Assignment (2 points, 1.5h)
│   Dependencies: None (foundational extension)
│   Deliverable: --project flag on issue create/update
│
├── Story 1.2: Project Creation & Archival (3 points, 3.5h)
│   Dependencies: Story 1.1 (enables issue-project workflow)
│   Deliverable: project create and project archive commands
│
└── Story 1.3: Project Updates & Enhanced Display (5 points, 4h)
    Dependencies: Story 1.2 (requires project CRUD foundation)
    Deliverable: project update command + enhanced display
```

**Dependency Validation:** ✅ Valid sequence
- Story 1.1 is independent (extends existing commands)
- Story 1.2 depends only on Story 1.1
- Story 1.3 depends only on Story 1.2
- No forward dependencies detected

---

## Stories - Epic 1

### Story 1.1: Issue-Project Assignment

As a linctl user (developer, project manager, or automation script),
I want to assign projects to issues during creation and updates,
So that I can manage the complete issue-project workflow from the terminal without switching to Linear's web UI.

**Acceptance Criteria:**

**AC #1:** Given a valid project UUID, when I run `linctl issue create --title "Test" --team ENG --project PROJECT-UUID`, then the issue is created and assigned to the specified project, and the output shows the project assignment.

**AC #2:** Given an existing issue and valid project UUID, when I run `linctl issue update ISS-123 --project PROJECT-UUID`, then the issue's project is updated to the new project.

**AC #3:** Given an issue with an existing project assignment, when I run `linctl issue update ISS-123 --project unassigned`, then the project assignment is removed from the issue.

**AC #4:** Given an invalid project UUID, when I attempt to assign it to an issue, then the command fails with a clear error message: "Project 'INVALID-UUID' not found".

**AC #5:** JSON output includes the project field showing project ID and name when `--json` flag is used.

**Prerequisites:** None (extends existing issue commands)

**Technical Notes:** Add `--project` flag to issueCreateCmd and issueUpdateCmd in cmd/issue.go. Include projectId field in GraphQL mutation input map. Handle "unassigned" special value by setting projectId to nil.

**Estimated Effort:** 2 points (1.5 hours)

---

### Story 1.2: Project Creation & Archival

As a project manager or team lead using linctl,
I want to create new projects and archive completed ones from the CLI,
So that I can manage project lifecycle without leaving my terminal workflow.

**Acceptance Criteria:**

**AC #1:** Given a project name and team key, when I run `linctl project create --name "Q1 Backend" --team ENG`, then a new project is created in Linear with default values (state: planned, priority: 0).

**AC #2:** Given project creation with optional fields, when I run `linctl project create --name "Test" --team ENG --state started --priority 1 --description "Test project"`, then the project is created with all specified field values.

**AC #3:** Given a valid project UUID, when I run `linctl project archive PROJECT-UUID`, then the project is archived in Linear and a success message is displayed.

**AC #4:** Given missing required fields, when I attempt project creation without --name or --team, then the command fails with error: "Both --name and --team are required".

**AC #5:** Given an invalid team key, when I attempt project creation, then the command fails with error: "Team 'INVALID' not found. Use 'linctl team list' to see available teams."

**AC #6:** Given invalid state or priority values, when I attempt project creation, then validation fails before API call with clear error messages.

**AC #7:** All commands support --json, --plaintext output formats and follow existing linctl formatting conventions.

**Prerequisites:** Story 1.1 complete (enables full issue-project workflow)

**Technical Notes:** Create projectCreateCmd and projectArchiveCmd in cmd/project.go. Implement CreateProject() and ArchiveProject() methods in pkg/api/queries.go with proper GraphQL mutations. Validate required fields and resolve team key to UUID before API call.

**Estimated Effort:** 3 points (3.5 hours)

---

### Story 1.3: Project Updates & Enhanced Display

As a project manager using linctl for project tracking,
I want to update project fields and see complete project information in list/detail views,
So that I can manage project state, priority, initiatives, labels, descriptions, and summary entirely from the CLI.

**Acceptance Criteria:**

**AC #1:** Given a project UUID and updated field, when I run `linctl project update PROJECT-UUID --name "New Name"`, then the project name is updated in Linear.

**AC #2:** Given a project UUID, when I run `linctl project update PROJECT-UUID --state started`, then the project state is updated and validated against allowed values (planned, started, paused, completed, canceled).

**AC #3:** Given multiple field updates, when I run `linctl project update PROJECT-UUID --state started --priority 1 --label "urgent,backend"`, then all specified fields are updated in a single API call.

**AC #3.1:** Given a project UUID, when I run `linctl project update PROJECT-UUID --description "Full description"`, then the project description is updated.

**AC #3.2:** Given a project UUID, when I run `linctl project update PROJECT-UUID --summary "Short summary"`, then the project shortSummary field is updated.

**AC #4:** Given no field flags provided, when I run `linctl project update PROJECT-UUID`, then the command fails with error: "At least one field to update is required".

**AC #5:** Given invalid state value, when I attempt update, then validation fails with error: "Invalid state. Must be one of: planned, started, paused, completed, canceled".

**AC #6:** Given invalid priority value, when I attempt update, then validation fails with error: "Priority must be between 0 and 4".

**AC #7:** Given a project UUID, when I run `linctl project get PROJECT-UUID`, then the output displays all fields including description, shortSummary, state, priority, initiative, and labels.

**AC #8:** When I run `linctl project list`, then the table output includes State and Priority columns for each project.

**AC #9:** All commands support --json and --plaintext output formats with complete field data.

**Prerequisites:** Story 1.2 complete (requires project CRUD foundation)

**Technical Notes:** Create projectUpdateCmd in cmd/project.go with multi-field flag support (including --description and --summary). Implement UpdateProject() method in pkg/api/queries.go. Enhance GraphQL queries in projectGetCmd and projectListCmd to include new fields. Update output formatters to display description, shortSummary, state, priority, initiative, labels. Use cmd.Flags().Changed() to detect which fields were explicitly provided.

**Estimated Effort:** 5 points (4.5 hours)

---

## Implementation Timeline - Epic 1

**Total Story Points:** 10 points

**Estimated Timeline:** 1 sprint (1-2 weeks at 1-2 points per day)

**Implementation Sequence:**
1. **Story 1.1** (Days 1-2) → Issue-project assignment foundation
2. **Story 1.2** (Days 3-4) → Project CRUD operations
3. **Story 1.3** (Days 5-6) → Project updates and enhanced display
4. **Integration Testing** (Day 7) → Complete workflow validation
5. **Documentation** (Day 7) → README updates and manual test procedures

---

## Tech-Spec Reference

See [tech-spec.md](./tech-spec.md) for complete technical implementation details including:
- Brownfield codebase analysis and existing patterns
- GraphQL mutation specifications
- Complete file paths and code locations
- Validation rules and error handling
- Testing strategy and procedures
- Deployment process and rollback plan
