# Engineering Backlog

This backlog collects cross-cutting or future action items that emerge from reviews and planning.

Routing guidance:

- Use this file for non-urgent optimizations, refactors, or follow-ups that span multiple stories/epics.
- Must-fix items to ship a story belong in that story’s `Tasks / Subtasks`.
- Same-epic improvements may also be captured under the epic Tech Spec `Post-Review Follow-ups` section.

| Date | Story | Epic | Type | Severity | Owner | Status | Notes |
| ---- | ----- | ---- | ---- | -------- | ----- | ------ | ----- |
| 2025-11-06 | 1.1 | 1 | Bug | High | TBD | Open | Standardize invalid project error to "Project '<value>' not found"; see pkg/api/client.go:102 and cmd/issue.go:913/1094 |
| 2025-11-06 | 1.1 | 1 | Enhancement | Med | TBD | Open | Add UUID format validation for `--project` (allow `unassigned`) |
| 2025-11-06 | 1.2 | 1 | Enhancement | High | TBD | Done | Include project name in archive success output; implemented at cmd/project.go:791-805 |
| 2025-11-06 | 1.2 | 1 | TechDebt | Med | TBD | Done | Add unit/integration tests for project create/archive and validators (JSON/plaintext) — see pkg/api/queries_test.go, cmd/project_cmd_test.go, cmd/project_test.go |
| 2025-11-06 | 1.2 | 1 | Bug | Med | TBD | Done | GetTeam now resolves by key via teams filter with id fallback; see pkg/api/queries.go:1224-1279 |
| 2025-11-06 | 1.2 | 1 | Enhancement | Low | TBD | Done | Validate --target-date format (YYYY-MM-DD) before API call (cmd/project.go:664-670) |
