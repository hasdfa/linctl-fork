# Quick Setup Workflow - Validation Checklist

## Pre-Execution Checks

- [ ] User has gh CLI installed (if GitHub creation needed)
- [ ] User has git installed
- [ ] Current directory is intended project location
- [ ] No conflicting .git directory (or user wants to use existing)

## Execution Validation

### Directory Structure
- [ ] Standard directories created (src/, tests/, docs/, .bmad/, .claude/commands/)
- [ ] Project-type specific directories created if needed
- [ ] All directories have proper permissions

### Git Setup
- [ ] Git repository initialized (.git/ exists)
- [ ] .gitignore created with appropriate template
- [ ] .gitignore includes BMAD and Claude directories
- [ ] Common patterns added (.env, node_modules/, etc.)

### Foundation Documents
- [ ] idea.md created with project description
- [ ] CLAUDE.md created with development guidelines
- [ ] AGENTS.md created with agent configuration
- [ ] README.md created with setup instructions
- [ ] All documents contain project-specific content (not just templates)

### BMAD Installation
- [ ] .bmad/ directory structure created
- [ ] Correct module installed (bmm, etc.)
- [ ] Module configuration files present
- [ ] Agent symlinks created
- [ ] Manifest files updated

### GitHub Repository
- [ ] Private repository created (if not skipped)
- [ ] Remote origin configured correctly
- [ ] Initial commit made
- [ ] Code pushed to main branch
- [ ] Repository URL captured and displayed

### Template Sync
- [ ] Templates fetched from agent_customisations repo
- [ ] Agents copied to correct location
- [ ] Workflows copied to correct location
- [ ] Preferences copied to correct location
- [ ] Local customizations preserved (*.local.* files not overwritten)

## Post-Execution Validation

- [ ] All skipped steps properly bypassed
- [ ] Dry-run mode didn't create any files
- [ ] Summary statistics accurate
- [ ] Next steps displayed clearly
- [ ] No errors in execution log
- [ ] User can immediately start development

## Error Recovery

### Common Issues

**Git already initialized:**
- Skip git initialization step
- Verify .gitignore still created/updated

**gh CLI not found:**
- Skip GitHub creation
- Inform user to install gh CLI
- Provide manual instructions

**BMAD already installed:**
- Skip BMAD installation
- Verify version compatibility
- Offer to update/sync instead

**Network issues (template sync):**
- Retry with exponential backoff
- Fall back to local templates if available
- Continue without templates as last resort

**Permission errors:**
- Display clear error message
- Suggest sudo if appropriate
- Provide manual recovery steps

## Success Criteria

- [ ] Project is ready for immediate development
- [ ] All foundation documents in place
- [ ] Version control configured
- [ ] BMAD framework operational
- [ ] Templates up-to-date
- [ ] User knows next steps
