# Archive Project Workflow Validation Checklist

## Initial Explanation and Confirmation
- [ ] Archival overview displayed clearly
- [ ] Four main phases explained (checks, export, cleanup, archive)
- [ ] One-shot operation warning shown
- [ ] User asked to proceed with checks
- [ ] Workflow exits gracefully if user cancels
- [ ] No changes made if user cancels early

## Pre-Check 1: Git Status
- [ ] Git status command executed successfully
- [ ] Uncommitted changes detected correctly
- [ ] Modified files listed
- [ ] Untracked files listed
- [ ] Staged files listed
- [ ] File count accurate
- [ ] Warning displayed if uncommitted changes found
- [ ] Success message shown if clean
- [ ] has_uncommitted flag set correctly

## Pre-Check 2: Unpushed Commits
- [ ] Remote branch configured correctly
- [ ] git log command executed successfully
- [ ] Unpushed commits counted accurately
- [ ] Commit messages listed
- [ ] Warning displayed if unpushed commits found
- [ ] Success message shown if all pushed
- [ ] No remote scenario detected and warned
- [ ] has_unpushed flag set correctly
- [ ] has_no_remote flag set correctly

## Pre-Check 3: Open Linear Issues
- [ ] Linear project existence checked
- [ ] Linear issues queried successfully
- [ ] Open issues identified (Todo, In Progress, In Review states)
- [ ] Done/Canceled issues counted separately
- [ ] Issue list displayed with state and title
- [ ] Warning shown if open issues exist
- [ ] Success message shown if no open issues
- [ ] No Linear project scenario handled
- [ ] has_open_issues flag set correctly
- [ ] open_issue_count set correctly
- [ ] has_no_linear flag set correctly

## Pre-Check 4: Recent Activity
- [ ] Last commit date retrieved correctly
- [ ] Days since last commit calculated
- [ ] File modification times checked
- [ ] Most recent file identified (.py, .js, .ts, .md)
- [ ] 30-day threshold applied correctly
- [ ] Warning shown if activity in last 30 days
- [ ] Last commit date and days displayed
- [ ] Last file modification displayed
- [ ] Success message shown if no recent activity
- [ ] has_recent_activity flag set correctly

## Pre-Check 5: Running Processes
- [ ] lsof command executed (or equivalent)
- [ ] Processes in project directory detected
- [ ] Common dev servers checked (npm, node, python, etc.)
- [ ] Database processes checked (postgres, mysql, etc.)
- [ ] Process list displayed with PID and command
- [ ] Warning shown if running processes found
- [ ] Success message shown if no processes
- [ ] has_running_processes flag set correctly

## Pre-Check Summary
- [ ] Warning count calculated correctly
- [ ] All passed scenario displays success message
- [ ] Warnings present scenario displays warning summary
- [ ] Each warning category shown with status
- [ ] Options provided (address warnings, proceed, cancel)
- [ ] Summary clear and actionable

## Archive Reason and Confirmation
- [ ] User asked for archive reason
- [ ] Archive reason captured correctly
- [ ] Project name extracted from directory
- [ ] Archive date calculated (YYYY-MM-DD format)
- [ ] Archive name formatted correctly (date-project-name)
- [ ] Archive path constructed correctly
- [ ] Project size calculated (du -sh)
- [ ] File count calculated accurately
- [ ] Preview file generated (.bmad/archive-preview.md)
- [ ] Preview shows all archive details
- [ ] Preview shows pre-check warnings
- [ ] Preview shows all operations (export, cleanup, archive)
- [ ] Preview shows restoration instructions
- [ ] Preview saved and path shown to user

## Final Confirmation Handling
- [ ] Warning confirmation asked if warnings present
- [ ] Workflow exits if user declines to proceed with warnings
- [ ] Final confirmation asks for exact project name
- [ ] User must type exact project name to confirm
- [ ] Cancel option available
- [ ] Workflow exits if user types "cancel"
- [ ] Workflow exits if user types wrong name
- [ ] Archival only proceeds with correct project name
- [ ] Confirmation messages clear and appropriate

## Archive Directory Creation
- [ ] Main archive directory created (~/Documents/archived/{name}/)
- [ ] exported-data subdirectory created
- [ ] exported-data/docs subdirectory created
- [ ] exported-data/config subdirectory created
- [ ] Directory permissions correct
- [ ] Success message shown
- [ ] Archive path verified

## Linear Data Export
- [ ] Skipped correctly if skip_export flag set
- [ ] Skipped correctly if has_no_linear is true
- [ ] All Linear issues queried successfully
- [ ] Issues include all required fields (title, description, state, etc.)
- [ ] Comments included for each issue
- [ ] Created/updated dates included
- [ ] Data formatted as valid JSON
- [ ] Saved to exported-data/linear-issues.json
- [ ] Linear project metadata exported
- [ ] Saved to exported-data/linear-project.json
- [ ] Issue count displayed in success message
- [ ] Export successful confirmation shown

## Git Metadata Export
- [ ] Skipped correctly if skip_export flag set
- [ ] Remote URLs captured (git remote -v)
- [ ] Branch list captured (git branch -a)
- [ ] Commit history captured (last 100 commits)
- [ ] Contributors captured (git shortlog -sn)
- [ ] Tags captured (git tag)
- [ ] Data formatted as valid JSON
- [ ] Saved to exported-data/git-metadata.json
- [ ] Success message shown

## Project Statistics Export
- [ ] Skipped correctly if skip_export flag set
- [ ] Total file count calculated (excluding node_modules, .git)
- [ ] Total size calculated (du -sh)
- [ ] Lines of code calculated (cloc or manual)
- [ ] LOC broken down by language
- [ ] Dependencies identified (package.json, requirements.txt)
- [ ] Last modified date captured
- [ ] Data formatted as valid JSON
- [ ] Saved to exported-data/project-stats.json
- [ ] Success message shown

## Documentation Export
- [ ] Skipped correctly if skip_export flag set
- [ ] All .md files found (excluding node_modules)
- [ ] File count accurate
- [ ] Directory structure preserved
- [ ] Files copied to exported-data/docs/
- [ ] Nested directories handled correctly
- [ ] File permissions preserved
- [ ] Doc count displayed in success message
- [ ] Export successful confirmation shown

## Configuration Export
- [ ] Skipped correctly if skip_export flag set
- [ ] .env.example copied (.env excluded)
- [ ] config.example.yaml copied
- [ ] .editorconfig copied
- [ ] .prettierrc copied
- [ ] .eslintrc.js copied
- [ ] package.json copied
- [ ] requirements.txt copied
- [ ] pyproject.toml copied
- [ ] tsconfig.json copied
- [ ] Other config files identified and copied
- [ ] No secrets included (verified)
- [ ] Files copied to exported-data/config/
- [ ] Success message shown with "(secrets excluded)"

## WezTerm Launcher Removal
- [ ] Skipped correctly if skip_cleanup flag set
- [ ] WezTerm config file located
- [ ] Project path searched in config
- [ ] Project entry removed if found
- [ ] Config file remains valid after edit
- [ ] Success message shown if removed
- [ ] Info message shown if not found in config
- [ ] Info message shown if WezTerm config not found
- [ ] No errors if WezTerm not used

## Linear Project Closure
- [ ] Skipped correctly if skip_cleanup flag set
- [ ] Skipped correctly if has_no_linear is true
- [ ] Linear project state updated to "Canceled"
- [ ] Archive note added to project description
- [ ] Note includes archive date, reason, and location
- [ ] Update successful
- [ ] Success message shown
- [ ] Project not deleted (just canceled)

## IDE Recent Projects Removal
- [ ] Skipped correctly if skip_cleanup flag set
- [ ] VSCode storage.json located (if applicable)
- [ ] Claude Code recent projects checked (if applicable)
- [ ] Codex recent projects checked (if applicable)
- [ ] Project removed from recent lists
- [ ] Best-effort approach (no errors if not found)
- [ ] Success message shown

## Cache and Artifact Cleanup
- [ ] Skipped correctly if skip_cleanup flag set
- [ ] node_modules/ deleted (if exists)
- [ ] __pycache__/ and *.pyc deleted (if exist)
- [ ] dist/ deleted (if exists)
- [ ] build/ deleted (if exists)
- [ ] .next/, .nuxt/, .vite/ deleted (if exist)
- [ ] .pytest_cache/, .mypy_cache/, .ruff_cache/ deleted (if exist)
- [ ] .DS_Store, Thumbs.db deleted (if exist)
- [ ] Space freed calculated accurately
- [ ] Success message shows space freed

## Archive Metadata Creation
- [ ] ARCHIVE_METADATA.json created
- [ ] All required fields present (archive_date, project_name, etc.)
- [ ] Timestamps in ISO 8601 format
- [ ] Git remote URL included (or null)
- [ ] Linear project ID included (or null)
- [ ] Last commit date correct
- [ ] Last modified date correct
- [ ] File count accurate
- [ ] Total size accurate
- [ ] Archive reason included
- [ ] Archived by (user_name) included
- [ ] Warnings list accurate (empty if none)
- [ ] Restoration notes reference included
- [ ] Valid JSON format
- [ ] Saved to archive directory root
- [ ] Success message shown

## Restoration Instructions Creation
- [ ] RESTORATION.md created
- [ ] Archive date and reason shown
- [ ] Quick restore steps clear and accurate
- [ ] Extract command correct
- [ ] Move command appropriate
- [ ] Dependency install commands project-appropriate
- [ ] Git remote verification step included
- [ ] WezTerm re-add step mentioned (manual)
- [ ] Linear reopen step mentioned (manual)
- [ ] Exported data section lists all export files
- [ ] Project-specific notes included (if any)
- [ ] Metadata reference included
- [ ] File saved to archive directory root
- [ ] Success message shown

## Project Compression
- [ ] Compression command executed from correct directory
- [ ] tar.gz format used
- [ ] Project compressed successfully
- [ ] Exclusions applied correctly (.git optional, node_modules, etc.)
- [ ] Compressed size calculated
- [ ] Original size shown for comparison
- [ ] Compression ratio reasonable
- [ ] Archive created at correct location
- [ ] Success message shows sizes

## Archive Finalization
- [ ] Archive integrity check performed (tar -tzf)
- [ ] Test extraction successful (no corruption)
- [ ] Success message on valid archive
- [ ] Original project removed only after integrity check
- [ ] rm -rf executed on correct path
- [ ] Original project directory deleted
- [ ] Success message on removal
- [ ] Error and exit if archive invalid
- [ ] Original NOT removed if archive invalid
- [ ] Clear error message on integrity failure

## Final Summary
- [ ] Summary displayed with all key information
- [ ] Project name shown
- [ ] Archive location shown (full path)
- [ ] Archive size shown (compressed)
- [ ] Archive reason shown
- [ ] Archive contents listed (all exports and files)
- [ ] Cleanup operations listed (all performed actions)
- [ ] Original project removal noted
- [ ] Restoration section clear with quick commands
- [ ] Archive management note included
- [ ] System log entry created (if applicable)
- [ ] All statistics accurate (issue count, doc count, etc.)

## Skip Flags
- [ ] --skip-checks skips all pre-checks (if used)
- [ ] --skip-export skips all data exports (if used)
- [ ] --skip-cleanup skips all cleanup operations (if used)
- [ ] --dry-run shows preview without executing (if used)
- [ ] --force allows proceeding despite warnings (if used)
- [ ] Multiple flags can be combined
- [ ] Flags respected throughout workflow

## Dry Run Mode
- [ ] --dry-run flag detected correctly
- [ ] All checks performed (not skipped)
- [ ] Preview generated
- [ ] No actual archival performed
- [ ] No data exported
- [ ] No cleanup performed
- [ ] No project removed
- [ ] Preview report created
- [ ] User informed this is dry run
- [ ] Safe to run multiple times

## Error Handling
- [ ] Missing archive directory creation handled
- [ ] Permission errors reported clearly
- [ ] Linear API failures handled gracefully
- [ ] Git command failures don't crash workflow
- [ ] Missing tools (lsof, tar, du) handled
- [ ] Partial exports still saved
- [ ] Failed cleanup doesn't stop archival
- [ ] Archive integrity failure prevents removal
- [ ] Error messages actionable and specific
- [ ] Errors logged appropriately

## Safety Features
- [ ] Type-to-confirm prevents accidental archival
- [ ] Archive integrity check mandatory
- [ ] All data exported before cleanup
- [ ] Restoration instructions always created
- [ ] Metadata always preserved
- [ ] No partial archival (all-or-nothing)
- [ ] Original only removed after verified archive
- [ ] Can cancel at multiple points
- [ ] Clear warnings for all risks

## User Experience
- [ ] Moderate verbosity (explains what's happening)
- [ ] Clear phase indicators
- [ ] Warnings always shown
- [ ] Progress messages informative
- [ ] Success/failure messages clear
- [ ] Next steps always provided
- [ ] Restoration path easy to find
- [ ] Final summary comprehensive
- [ ] No overwhelming output
- [ ] Can preview with --dry-run before real run

## Integration Points
- [ ] Linear API/MCP tools used correctly
- [ ] Git commands executed properly
- [ ] WezTerm config parsed correctly
- [ ] File system operations safe
- [ ] Tar compression compatible
- [ ] JSON exports valid
- [ ] Directory operations atomic where possible

## Final Validation
- [ ] Workflow completes successfully
- [ ] Archive created and valid
- [ ] Original project removed (if not dry run)
- [ ] All data exported
- [ ] All cleanup performed
- [ ] Metadata complete
- [ ] Restoration instructions accurate
- [ ] Archive can be extracted successfully
- [ ] Extracted project matches original
- [ ] Archive organized correctly in archive directory
- [ ] No data loss
- [ ] No corrupted files
