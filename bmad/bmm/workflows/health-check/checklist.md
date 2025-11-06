# Health Check Workflow Validation Checklist

## Git Status Check
- [ ] Uncommitted changes detected and reported
- [ ] Current branch identified
- [ ] Branch status vs remote shown (ahead/behind)
- [ ] Detached HEAD state detected if applicable
- [ ] Findings added to appropriate priority level

## GitHub Integration Check
- [ ] Remote origin configuration checked
- [ ] Unpushed commits counted and reported
- [ ] Behind remote status detected
- [ ] Missing remote warning issued if not configured
- [ ] GitHub repo link validated if present

## Linear Integration Check
- [ ] Linear project existence verified
- [ ] Team association confirmed
- [ ] Open issues counted by state
- [ ] Overdue tasks identified
- [ ] Label configuration validated
- [ ] GitHub repo linkage verified
- [ ] Missing Linear project flagged if applicable

## WezTerm Launcher Check
- [ ] WezTerm config file located
- [ ] Project path existence in launcher verified
- [ ] Correct section placement validated
- [ ] Missing entry flagged with inferred section
- [ ] Auto-fix option provided for missing entry

## BMAD Installation Check
- [ ] .bmad/ directory existence verified
- [ ] Module installation status checked
- [ ] Expansion packs enumerated
- [ ] Claude Code integration detected
- [ ] Codex integration detected
- [ ] Gemini CLI integration detected
- [ ] Missing IDE integrations flagged
- [ ] Version compatibility noted

## Dependencies Freshness Check
- [ ] package.json detected (Node.js projects)
- [ ] requirements.txt detected (Python projects)
- [ ] Outdated npm packages identified (if applicable)
- [ ] Outdated pip packages identified (if applicable)
- [ ] Security vulnerabilities flagged
- [ ] Update recommendations provided
- [ ] Lock file consistency verified

## Test Coverage Check
- [ ] Test directory existence verified (tests/, __tests__, spec/)
- [ ] Test framework identified (Jest, Vitest, pytest, etc.)
- [ ] Test file count reported
- [ ] Coverage configuration detected
- [ ] Recent test run status checked
- [ ] Missing tests flagged as WARNING
- [ ] Coverage percentage shown (if available)

## Documentation Completeness Check
- [ ] README.md existence verified
- [ ] README content quality assessed (length, sections)
- [ ] idea.md checked (foundation doc)
- [ ] CLAUDE.md checked (AI context)
- [ ] AGENTS.md checked (agent definitions)
- [ ] blueprint.md checked (technical design)
- [ ] Missing docs flagged with explanations
- [ ] Stale docs detected (>90 days without update)
- [ ] Doc freshness relative to code changes assessed

## Logging System Check
- [ ] .logs/ directory existence verified
- [ ] .logs/all.ndjson main log file checked
- [ ] .logs/inbox/ subdirectory verified
- [ ] Redaction configuration checked
- [ ] Log aggregator health verified (if running)
- [ ] Recent log entries validated (staleness check)
- [ ] .gitignore configuration for logs verified
- [ ] Missing logging setup flagged with /logs.init reference

## Maintenance Schedule Check
- [ ] .bmad/last-maintenance timestamp file checked
- [ ] Days since last maintenance calculated
- [ ] 30-day threshold enforcement
- [ ] Overdue maintenance flagged as WARNING
- [ ] First-time project detection (no timestamp) handled
- [ ] Maintenance command suggested if overdue

## Findings Compilation
- [ ] All findings categorized by priority (CRITICAL/WARNING/INFO)
- [ ] Each finding has category, issue, explanation, impact, fix
- [ ] Auto-fix availability accurately marked
- [ ] Findings sorted by priority (CRITICAL first)
- [ ] Duplicate findings eliminated
- [ ] Finding count per category calculated

## Auto-fix Offering
- [ ] Auto-fixable issues identified correctly
- [ ] User asked which auto-fixes to apply
- [ ] Fixes executed safely (with confirmation)
- [ ] Git initialization auto-fix available
- [ ] WezTerm entry addition auto-fix available
- [ ] .gitignore updates auto-fix available
- [ ] Logging initialization auto-fix available
- [ ] Post-fix validation performed

## Summary and Recommendations
- [ ] Executive summary generated (1-3 sentences)
- [ ] Priority breakdown shown (count per level)
- [ ] Critical issues highlighted prominently
- [ ] Next steps provided (prioritized action items)
- [ ] Relevant workflow commands suggested (maintenance, sync-standards, etc.)
- [ ] Workflow completion time noted

## Output Quality (Chatty Mode)
- [ ] Detailed explanations provided for each issue
- [ ] Context given for why issues matter
- [ ] Impact clearly stated
- [ ] Fix instructions actionable and specific
- [ ] Friendly, helpful tone maintained
- [ ] Not overwhelming (issues grouped logically)
- [ ] Visual formatting (emojis, priority indicators) used effectively

## Skip Flags Respected
- [ ] --skip-git honored if provided
- [ ] --skip-linear honored if provided
- [ ] --skip-bmad honored if provided
- [ ] --skip-deps honored if provided
- [ ] --skip-tests honored if provided
- [ ] --skip-docs honored if provided
- [ ] --skip-logs honored if provided
- [ ] --skip-maintenance honored if provided
- [ ] --quick flag limits to critical checks only

## Error Handling
- [ ] Missing tools handled gracefully (npm, pip, git)
- [ ] Missing files don't crash workflow
- [ ] API failures (Linear, GitHub) reported clearly
- [ ] Partial results still useful if some checks fail
- [ ] Error messages informative and actionable

## Final Validation
- [ ] Workflow completes successfully
- [ ] All enabled checks executed
- [ ] Findings accurately reflect project state
- [ ] No false positives in critical category
- [ ] Recommendations relevant to project type
- [ ] User understands next steps clearly
- [ ] Project remains in valid state (no modifications made)
