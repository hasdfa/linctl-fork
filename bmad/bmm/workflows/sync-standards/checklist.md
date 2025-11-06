# Sync Standards Workflow Validation Checklist

## Repository Verification
- [ ] agent_customisations repository location checked
- [ ] Repository exists at ~/Documents/raegis_labs/agent_customisations
- [ ] templates/ directory structure verified
- [ ] Error displayed and workflow exits if repository not found
- [ ] Sync source root path set correctly
- [ ] Tracking variables initialized (synced, conflict, skipped, error counts)

## Agent Templates Sync
- [ ] Source path correct: templates/agents/*.yaml
- [ ] Target path correct: .bmad/expansion-packs/*/agents/*.yaml
- [ ] All template agents enumerated
- [ ] Target agents enumerated correctly
- [ ] *.local.* files identified and protected
- [ ] Content comparison performed (md5 or diff)
- [ ] New template agents detected
- [ ] Modified local agents detected as conflicts
- [ ] Non-conflicted agents synced automatically (silent)
- [ ] No prompts when no conflicts
- [ ] Conflict summary shown when conflicts detected
- [ ] User given overwrite/keep/review options
- [ ] Overwrite option replaces local with template
- [ ] Keep option preserves local versions
- [ ] Review option shows diffs and asks per file
- [ ] --force flag overwrites without prompting
- [ ] *.local.* files never overwritten (even with --force)
- [ ] Symlinks in agents/ directory updated after sync
- [ ] Synced count incremented correctly
- [ ] Conflict count tracked accurately
- [ ] Skipped count tracked accurately

## BMAD Framework Sync
- [ ] Source path correct: templates/bmad/
- [ ] Target path correct: bmad/
- [ ] config.yaml files identified and protected
- [ ] All config.yaml files skipped automatically
- [ ] BMAD core files compared with templates
- [ ] Modified BMAD files detected as conflicts
- [ ] Non-conflicted files synced automatically (silent)
- [ ] Conflict summary shown when conflicts detected
- [ ] User given overwrite/keep options
- [ ] Overwrite option replaces local files
- [ ] Keep option preserves local files
- [ ] --force flag overwrites non-protected files
- [ ] Protected config files never touched
- [ ] Synced count incremented correctly
- [ ] Conflict count tracked accurately

## Workflow Templates Sync
- [ ] Source path correct: templates/workflows/
- [ ] Target path correct: .bmad/expansion-packs/*/workflows/
- [ ] *.local.* workflow files identified and protected
- [ ] Workflow templates compared with local workflows
- [ ] Modified workflows detected as conflicts
- [ ] New workflow templates detected
- [ ] Non-conflicted workflows synced automatically (silent)
- [ ] Conflict summary shown when conflicts detected
- [ ] User given overwrite/keep options
- [ ] Recommendation shown for *.local.* naming
- [ ] --force flag overwrites non-protected workflows
- [ ] *.local.* workflows never overwritten
- [ ] Synced count incremented correctly
- [ ] Conflict count tracked accurately

## Code Style Templates Sync
- [ ] Source path correct: templates/code-style/
- [ ] Target path correct: {project-root}/
- [ ] .editorconfig file handled correctly
- [ ] .prettierrc file handled correctly
- [ ] .eslintrc.js file handled correctly
- [ ] pyproject.toml sections identified correctly
- [ ] Only [tool.black] and [tool.ruff] sections synced in pyproject.toml
- [ ] Other pyproject.toml sections preserved
- [ ] New style files created if missing
- [ ] Existing style files compared with templates
- [ ] Modified style files detected as conflicts
- [ ] Non-conflicted files synced automatically (silent)
- [ ] Conflict summary shown when conflicts detected
- [ ] User given overwrite/keep options
- [ ] pyproject.toml merge note shown
- [ ] --force flag applies template styles
- [ ] Synced count incremented correctly
- [ ] Conflict count tracked accurately

## GitHub Workflow Templates Sync
- [ ] Source path correct: templates/github/workflows/
- [ ] Target path correct: .github/workflows/
- [ ] *.local.yml files identified and protected
- [ ] GitHub workflow templates compared with local
- [ ] Modified workflows detected as conflicts
- [ ] New templates detected
- [ ] Non-conflicted workflows synced automatically (silent)
- [ ] Conflict summary shown when conflicts detected
- [ ] User given overwrite/keep options
- [ ] Recommendation shown for *.local.yml naming
- [ ] --force flag overwrites non-protected workflows
- [ ] *.local.yml workflows never overwritten
- [ ] .github/workflows/ directory created if missing
- [ ] Synced count incremented correctly
- [ ] Conflict count tracked accurately

## Conflict Handling
- [ ] Conflicts detected using md5 hash or diff
- [ ] Protected patterns never counted as conflicts
- [ ] Conflict count accurate across all categories
- [ ] Conflict threshold checked (>10 conflicts)
- [ ] "Skip all" offered if too many conflicts
- [ ] Individual conflict resolution works (review mode)
- [ ] Diffs shown correctly in review mode
- [ ] Per-file decisions applied correctly
- [ ] Conflict summary clear and actionable
- [ ] Impact and risk explained in summaries

## Silent Execution
- [ ] No prompts when zero conflicts
- [ ] No category-by-category confirmations (unlike maintenance)
- [ ] Only errors and conflicts shown
- [ ] Progress indicators minimal or absent
- [ ] Final summary always shown
- [ ] Verbosity level "silent" respected
- [ ] Only non-zero counts displayed

## Skip Flags
- [ ] --skip-agents skips agent template sync
- [ ] --skip-bmad skips BMAD framework sync
- [ ] --skip-workflows skips workflow template sync
- [ ] --skip-code-style skips code style sync
- [ ] --skip-github skips GitHub workflow sync
- [ ] --force overwrites all conflicts without prompting
- [ ] --dry-run previews changes without applying
- [ ] Multiple skip flags can be combined

## Dry Run Mode
- [ ] --dry-run flag detected correctly
- [ ] All comparisons performed
- [ ] All conflicts detected
- [ ] Preview report generated
- [ ] No actual file changes made
- [ ] Synced count shows would-be syncs
- [ ] Conflict count shows would-be conflicts
- [ ] Preview file created with detailed report
- [ ] User informed this is dry run only

## Force Mode
- [ ] --force flag detected correctly
- [ ] All conflicts overwritten automatically
- [ ] No prompts shown for conflicts
- [ ] Protected patterns still respected (*.local.*, config.yaml)
- [ ] Warning shown that force mode is active
- [ ] Synced count includes force-overwritten files
- [ ] Conflict count zero (all auto-resolved)

## Logging and Summary
- [ ] Sync log written to .bmad/last-sync-standards.log
- [ ] Log includes timestamp (ISO 8601)
- [ ] Log includes all counts (synced, conflict, skipped, error)
- [ ] Log includes category breakdown
- [ ] Summary displayed to user
- [ ] Non-zero counts highlighted
- [ ] Zero counts omitted from summary
- [ ] Next steps shown when files synced
- [ ] "Already in sync" message when nothing to sync
- [ ] Log file path shown in summary
- [ ] Log file appended (not overwritten) for history

## Protected Patterns
- [ ] *.local.* files never overwritten
- [ ] *.local.yaml agent files protected
- [ ] *.local.yml GitHub workflows protected
- [ ] *.local.md workflow files protected (if any)
- [ ] config.yaml files protected
- [ ] bmm/config.yaml protected
- [ ] All module config.yaml files protected
- [ ] Protected files automatically skipped (silent)
- [ ] Protected files not counted in conflicts
- [ ] User not prompted about protected files

## File Operations
- [ ] File comparisons use efficient method (md5/checksum)
- [ ] File copies preserve permissions
- [ ] File copies preserve timestamps (if appropriate)
- [ ] Directories created if missing
- [ ] Symlinks handled correctly
- [ ] No broken symlinks created
- [ ] File operations atomic (no partial writes)
- [ ] Failed operations don't corrupt files

## Error Handling
- [ ] Missing source files handled gracefully
- [ ] Missing target directories created
- [ ] Permission errors reported clearly
- [ ] I/O errors don't crash workflow
- [ ] Errors logged to sync log file
- [ ] Error count incremented on failures
- [ ] Remaining files synced even after errors
- [ ] Partial sync still recorded in log
- [ ] Error messages actionable and specific

## Integration Points
- [ ] Uses rsync or similar efficient sync
- [ ] Respects .gitignore patterns
- [ ] Git repo status checked (if git project)
- [ ] File permissions preserved
- [ ] No conflicts with running processes
- [ ] Safe to run while project is open

## Safety Features
- [ ] No data loss from sync operations
- [ ] Protected files absolutely never touched
- [ ] All changes reversible via git
- [ ] Dry run available before real sync
- [ ] Audit trail in log file
- [ ] Force mode warns user
- [ ] Conflict detection prevents accidental overwrites

## Performance
- [ ] Batch operations for efficiency
- [ ] Checksums avoid unnecessary copies
- [ ] Parallel operations where safe
- [ ] Large file sync optimized
- [ ] No unnecessary file reads
- [ ] Minimal disk I/O
- [ ] Completes in reasonable time (<30s for typical project)

## User Experience
- [ ] Very low interactivity (as designed)
- [ ] Silent when everything syncs cleanly
- [ ] Clear messages on conflicts only
- [ ] Final summary always shown
- [ ] Summary concise and informative
- [ ] Next steps actionable
- [ ] Log file available for details
- [ ] Can use --dry-run for preview
- [ ] Can use skip flags for selective sync

## Final Validation
- [ ] Workflow completes successfully
- [ ] All enabled categories processed
- [ ] Counts accurate (synced + conflict + skipped + error = total)
- [ ] Log file written correctly
- [ ] Summary accurate and complete
- [ ] Project remains in valid state
- [ ] No corrupted files
- [ ] Symlinks valid
- [ ] Tests still pass after sync
- [ ] Project builds successfully
- [ ] Git status clean or shows expected changes
