# Project Scratchpad: Add Parent Issue Linking to `linctl issue update`

## Overview
Add `--parent` flag to `linctl issue update` command to set, change, or remove parent issue relationships.

## Status: COMPLETE

---

## Implementation Plan (Finalized)

### Files Modified
| File | Changes |
|------|---------|
| `cmd/issue.go` | Added `--parent` flag (line 1160), handling logic (lines 1054-1080), help text (lines 928-929), enhanced output (lines 1099-1109) |
| `pkg/api/queries.go` | Added `parent` field to UpdateIssue mutation response (lines 1122-1126) |
| `smoke_test.sh` | Added test cases for `--parent` flag (lines 151-152) |

---

## Execution Log

### 2024-12-06 - Analysis Phase
- [x] Read implementation plan from `linctl-parent-issue-linking-plan.md`
- [x] Analyzed `cmd/issue.go` issueUpdateCmd patterns (flag handling, validation, output)
- [x] Analyzed `pkg/api/queries.go` UpdateIssue mutation structure
- [x] Analyzed test coverage (0% Go tests, smoke tests at smoke_test.sh)
- [x] Created project scratchpad

### 2024-12-06 - Implementation Phase
- [x] Step 1: Added `--parent` flag registration in init()
- [x] Step 2: Added parent flag handling logic with validation
- [x] Step 3: Updated GraphQL mutation response to include parent
- [x] Step 4: Updated command help text with examples
- [x] Step 5: Enhanced output to display parent info
- [x] Added smoke tests for --parent flag documentation
- [x] Code review passed (engineering-agent, code-review-agent)
- [x] Build and vet passed
- [x] Smoke tests passed (39/41 - 2 pre-existing failures unrelated to this feature)

### Key Findings from Analysis:
1. **Flag Pattern**: Use `cmd.Flags().Changed()` before processing
2. **Null Pattern**: Assign `nil` for "none"/"null"/"" values (see assignee/due-date)
3. **Validation Pattern**: Call `GetIssue` to validate parent exists
4. **Error Pattern**: Use `output.Error()` + `os.Exit(1)`
5. **Output Pattern**: Three modes - JSON (full object), plaintext (simple), styled (with colors)
6. **Test Baseline**: 0% Go test coverage, rely on smoke_test.sh

---

## Task Progress

| Task | Status | Agent | Notes |
|------|--------|-------|-------|
| Analysis | ✅ Complete | code-analysis-agent | Patterns identified |
| Step 1: Flag registration | ✅ Complete | engineering-agent | Line 1160 |
| Step 2: Flag handling | ✅ Complete | engineering-agent | Lines 1054-1080 |
| Step 3: GraphQL mutation | ✅ Complete | engineering-agent | Lines 1122-1126 |
| Step 4: Help text | ✅ Complete | engineering-agent | Lines 928-929 |
| Step 5: Output enhancement | ✅ Complete | engineering-agent | Lines 1099-1109 |
| Code review | ✅ Complete | code-review-agent | Approved |
| Smoke tests | ✅ Complete | - | Lines 151-152 |
| Final validation | ✅ Complete | - | Build, vet, tests pass |

---

## Validation Results

### Build & Static Analysis:
- [x] `go build ./...` passes
- [x] `go vet ./...` passes
- [x] Code follows established patterns (verified by code-review-agent)

### Smoke Tests (39/41 passed):
- [x] `issue update help` - PASS
- [x] `issue update --parent flag` - PASS (verifies flag in help)
- [x] All other issue tests - PASS
- [ ] `project get` - PRE-EXISTING FAILURE (Entity not found)
- [ ] `project get (plaintext)` - PRE-EXISTING FAILURE

### Manual Testing Required:
- [ ] `linctl issue update CHILD-123 --parent PARENT-456` (set parent)
- [ ] `linctl issue update CHILD-123 --parent none` (remove parent)
- [ ] `linctl issue update CHILD-123 --parent CHILD-123` (self-reference error)
- [ ] `linctl issue update CHILD-123 --parent INVALID-999` (non-existent parent error)
- [ ] `linctl issue update CHILD-123 --parent PARENT-456 --json` (JSON output)
- [ ] `linctl issue update CHILD-123 --parent PARENT-456 -p` (plaintext output)

---

## Implementation Details

### Features Implemented:
1. **Set parent**: `linctl issue update ISSUE-1 --parent PARENT-1`
2. **Remove parent**: `linctl issue update ISSUE-1 --parent none` (also accepts "null" or "")
3. **Validation**: Checks parent issue exists before update
4. **Self-reference prevention**: Error if issue is set as its own parent
5. **Enhanced output**: Shows parent info after successful update

### Code Quality:
- Follows existing patterns exactly (assignee, state, due-date handling)
- Proper error handling with context-specific messages
- Consistent output formatting across JSON/plaintext/styled modes
- No security vulnerabilities identified

---

## Notes

### Edge Cases Handled:
1. Self-reference prevention (issue cannot be its own parent)
2. Non-existent parent (validate with GetIssue before update)
3. Cross-team parenting (allowed by Linear API)
4. Circular references (handled by Linear API, returns error)

### Patterns Reference:
- Assignee "unassigned" handling: `cmd/issue.go:969-970`
- State validation lookup: `cmd/issue.go:997-1034`
- Parent display in get command: `cmd/issue.go:634-639`
