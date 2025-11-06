# Maintenance Workflow Validation Checklist

## Initialization
- [ ] Last maintenance timestamp file read correctly
- [ ] Days since last maintenance calculated accurately
- [ ] 30-day threshold warning displayed when overdue
- [ ] First-time run detected when no timestamp exists
- [ ] Maintenance scope explanation shown clearly
- [ ] All 9 categories listed in overview

## Dead Code Cleanup (Category 1/9)
- [ ] Unused imports detected (Python and JS/TS)
- [ ] Commented code blocks identified
- [ ] Unreachable code found (after return/break, if(false))
- [ ] Unused functions/variables detected
- [ ] Preview file generated with file:line references
- [ ] User confirmation requested
- [ ] Dead code removed when confirmed
- [ ] Git history preserved during removal
- [ ] Summary updated with removal count
- [ ] Skipped gracefully when user declines
- [ ] --skip-dead-code flag respected

## Old Branches Cleanup (Category 2/9)
- [ ] Local branches enumerated correctly
- [ ] Merged branches identified (git branch --merged)
- [ ] Stale branches detected (>90 days no commits)
- [ ] Current branch excluded from cleanup
- [ ] main/master branches excluded
- [ ] Preview shows last commit dates
- [ ] User confirmation requested
- [ ] Branches deleted with git branch -d when confirmed
- [ ] Remote branches left untouched (local only)
- [ ] Summary updated with deletion count
- [ ] Skipped gracefully when user declines
- [ ] --skip-branches flag respected

## Stale Dependencies Cleanup (Category 3/9)
- [ ] Node.js project detected (package.json)
- [ ] Python project detected (requirements.txt or pyproject.toml)
- [ ] npm outdated run successfully (if Node.js)
- [ ] pip list --outdated run successfully (if Python)
- [ ] Unused dependencies detected (npx depcheck for Node.js)
- [ ] Preview shows current vs latest versions
- [ ] Major vs minor updates distinguished
- [ ] User confirmation requested
- [ ] Unused packages removed when confirmed
- [ ] Patch/minor updates applied when confirmed
- [ ] Major updates flagged for manual review
- [ ] Summary updated with update/removal count
- [ ] Skipped gracefully when user declines
- [ ] --skip-deps flag respected

## Old Logs Cleanup (Category 4/9)
- [ ] .logs/ directory scanned correctly
- [ ] Log files older than 30 days identified
- [ ] Rotated logs detected (*.log.1, *.log.2)
- [ ] Total size calculated accurately
- [ ] Preview shows age and size for each file
- [ ] User given choice: archive/delete/keep
- [ ] .logs/archive/ directory created when archiving
- [ ] Logs compressed by month (tar.gz) when archiving
- [ ] Compressed archives moved to correct location
- [ ] Logs deleted when user chooses delete
- [ ] Summary updated with action and space saved
- [ ] Skipped gracefully when user declines
- [ ] --skip-logs flag respected

## Temp Files Cleanup (Category 5/9)
- [ ] .DS_Store files detected (macOS)
- [ ] Thumbs.db detected (Windows)
- [ ] *.tmp and *.temp files found
- [ ] __pycache__/ directories identified
- [ ] node_modules/.cache/ detected
- [ ] .pytest_cache/, .mypy_cache/ found
- [ ] Build artifacts in dist/, build/ identified (if not gitignored)
- [ ] Total size calculated
- [ ] Preview shows categorized items
- [ ] User confirmation requested
- [ ] All temp files/directories deleted when confirmed
- [ ] Summary updated with deletion count and space saved
- [ ] Skipped gracefully when user declines
- [ ] --skip-temp flag respected

## Agent Templates Sync (Category 6/9)
- [ ] agent_customisations repository location checked
- [ ] Templates directory existence verified
- [ ] Local agents enumerated from .bmad/expansion-packs/*/agents/
- [ ] Template agents compared with local agents
- [ ] Updates identified with diff summary
- [ ] New agents in templates detected
- [ ] *.local.* files identified and protected
- [ ] Preview shows updates, new agents, protected files
- [ ] User confirmation requested
- [ ] Updated agents copied to project when confirmed
- [ ] New agents added to project when confirmed
- [ ] *.local.* files preserved (never overwritten)
- [ ] Symlinks in agents/ directory updated
- [ ] Summary updated with sync count
- [ ] Warning shown if repository not found
- [ ] Skipped gracefully when user declines
- [ ] --skip-agent-templates flag respected

## BMAD Framework Sync (Category 7/9)
- [ ] Current BMAD version read from .bmad/VERSION
- [ ] Latest version checked (if mechanism exists)
- [ ] Core files integrity verified
- [ ] Missing or outdated files identified
- [ ] Preview shows version change and file updates
- [ ] User confirmation requested
- [ ] Core files updated when confirmed
- [ ] VERSION file updated when confirmed
- [ ] Project configuration preserved
- [ ] Custom agents not affected
- [ ] Summary updated with version change
- [ ] Skipped gracefully when user declines
- [ ] --skip-bmad flag respected

## Linear Settings Sync (Category 8/9)
- [ ] Linear project existence checked
- [ ] Current Linear labels retrieved
- [ ] Current Linear states/workflow retrieved
- [ ] Templates or standards identified
- [ ] Missing labels detected
- [ ] State configuration compared
- [ ] Template differences identified
- [ ] Preview shows missing labels, state diffs, template changes
- [ ] User confirmation requested
- [ ] Missing labels added to Linear when confirmed
- [ ] Templates updated when confirmed
- [ ] Existing issues/projects not affected
- [ ] Summary updated with sync count
- [ ] Warning shown if no Linear project
- [ ] Skipped gracefully when user declines
- [ ] --skip-linear flag respected

## Foundation Docs Sync (Category 9/9)
- [ ] Foundation doc templates located
- [ ] Current docs checked: idea.md, CLAUDE.md, AGENTS.md, README.md
- [ ] Structural differences identified
- [ ] New template sections detected
- [ ] Preview shows diffs for each doc
- [ ] User warned about manual merge requirement
- [ ] Detailed diffs displayed when confirmed
- [ ] Manual review instructions provided
- [ ] Project-specific content preservation emphasized
- [ ] Summary notes manual review needed
- [ ] No automatic overwrites performed
- [ ] Skipped gracefully when user declines
- [ ] --skip-docs flag respected

## Timestamp Update
- [ ] Current timestamp written to .bmad/last-maintenance
- [ ] Timestamp in ISO 8601 format
- [ ] File created if doesn't exist
- [ ] Timestamp display shown to user
- [ ] Timestamp only updated on successful run
- [ ] Timestamp not updated if critical errors occurred

## Summary Report Generation
- [ ] Duration calculated (start to end time)
- [ ] Categories processed count accurate
- [ ] Total changes count accurate
- [ ] Cleanup summary includes all categories
- [ ] Sync summary includes all categories
- [ ] Each category result clearly stated
- [ ] Next steps provided
- [ ] Next maintenance date calculated (30 days)
- [ ] Preview files list included
- [ ] Summary saved to .bmad/last-maintenance-summary.md
- [ ] Summary displayed to user
- [ ] Summary includes test/commit instructions

## Preview Files
- [ ] Dead code preview generated if needed
- [ ] Branches preview generated if needed
- [ ] Dependencies preview generated if needed
- [ ] Logs preview generated if needed
- [ ] Temp files preview generated if needed
- [ ] Agent templates preview generated if needed
- [ ] BMAD framework preview generated if needed
- [ ] Linear settings preview generated if needed
- [ ] Foundation docs preview generated if needed
- [ ] All previews saved to .bmad/ directory
- [ ] Preview files include all required information
- [ ] Previews show impact and risk assessment

## Skip Flags Functionality
- [ ] --skip-cleanup skips all cleanup categories
- [ ] --skip-sync skips all sync categories
- [ ] --skip-dead-code skips dead code cleanup
- [ ] --skip-branches skips branch cleanup
- [ ] --skip-deps skips dependency cleanup
- [ ] --skip-logs skips log cleanup
- [ ] --skip-temp skips temp file cleanup
- [ ] --skip-agent-templates skips agent sync
- [ ] --skip-bmad skips BMAD framework sync
- [ ] --skip-linear skips Linear settings sync
- [ ] --skip-docs skips foundation docs sync
- [ ] --dry-run shows previews without making changes
- [ ] --yes auto-approves all categories (skips confirmations)

## User Experience
- [ ] Moderate verbosity (not silent, not chatty)
- [ ] Clear progress indicators (1/9, 2/9, etc.)
- [ ] Category-by-category confirmation flow
- [ ] Explanations provided for each category
- [ ] Impact and risk clearly stated in previews
- [ ] User can skip individual categories
- [ ] Changes are reversible via git
- [ ] No unexpected modifications
- [ ] Next steps are actionable and specific

## Safety Features
- [ ] Git history preserved for all code changes
- [ ] *.local.* files never overwritten
- [ ] Current branch never deleted
- [ ] main/master branches protected
- [ ] Remote branches not affected
- [ ] Project configuration preserved
- [ ] Custom agents not modified
- [ ] Existing Linear issues/projects not affected
- [ ] Foundation docs not auto-merged (manual only)

## Error Handling
- [ ] Missing tools handled gracefully (npm, pip, git)
- [ ] Missing files/directories don't cause crashes
- [ ] agent_customisations not found warning shown
- [ ] Linear API failures handled
- [ ] Package manager failures reported clearly
- [ ] Category failures don't stop entire workflow
- [ ] Errors recorded in summary
- [ ] Timestamp not updated on critical errors
- [ ] Preview files preserved for debugging

## Integration Points
- [ ] Uses /init-sync-bmad logic for agent syncing
- [ ] Respects .gitignore for temp file detection
- [ ] Uses git commands correctly
- [ ] Uses npm/yarn commands correctly (Node.js)
- [ ] Uses pip commands correctly (Python)
- [ ] linctl CLI used correctly for Linear operations
- [ ] File system operations safe and correct

## Final Validation
- [ ] Workflow completes successfully
- [ ] All enabled categories processed
- [ ] Changes made match previews
- [ ] Summary accurate and complete
- [ ] Preview files can be deleted after review
- [ ] Timestamp updated to current time
- [ ] Project remains in valid state
- [ ] Tests still pass after maintenance
- [ ] Commits can be made successfully
- [ ] Next maintenance date clearly communicated
