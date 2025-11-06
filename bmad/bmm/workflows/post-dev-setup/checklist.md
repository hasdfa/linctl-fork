# Post-Dev Setup Workflow Validation Checklist

## Service Detection
- [ ] Frontend framework correctly detected (if present)
- [ ] Backend framework correctly detected (if present)
- [ ] Database service correctly detected (if present)
- [ ] User warned if no services detected

## Port Assignment (if not skipped)
- [ ] portbase utility is available
- [ ] Deterministic base port calculated correctly
- [ ] .envrc file created with correct port assignments
- [ ] WEB_PORT, API_PORT, DB_PORT, etc. properly exported
- [ ] PROJECT_DOMAIN set correctly ({project-name}.localtest.me)
- [ ] direnv activated successfully
- [ ] Ports verified and loaded
- [ ] .envrc added to .gitignore

## Application Code Updates (if not skipped)
- [ ] Frontend code updated to use $WEB_PORT (if applicable)
- [ ] Backend code updated to use $API_PORT (if applicable)
- [ ] User informed of required code changes
- [ ] PORT priority pattern documented (PORT || WEB_PORT || default)

## Traefik Proxy (if not skipped and user opted in)
- [ ] Traefik availability checked
- [ ] Traefik started if not running (user confirmed)
- [ ] Project registered with Traefik using add-proj
- [ ] Traefik configuration reloaded
- [ ] Proxy URLs displayed to user
- [ ] Dashboard URL provided

## Unified Logging (if not skipped)
- [ ] .logs/ directory created
- [ ] .logs/inbox/ subdirectory created
- [ ] .logs/all.ndjson main log file created
- [ ] .logs/redaction.json config created
- [ ] Redaction patterns properly formatted
- [ ] .logs/ added to .gitignore
- [ ] !.logs/redaction.json kept in git
- [ ] User instructed on how to feed logs

## CLI Menu (if not skipped)
- [ ] ./cli-menu script created
- [ ] Script is executable (chmod +x)
- [ ] Logging to .logs/inbox/cli-menu.log configured
- [ ] Menu functions properly templated
- [ ] Project name substituted correctly
- [ ] cli-menu added to .gitignore
- [ ] User informed how to run and customize

## User Experience
- [ ] Dry-run mode previewed all changes correctly
- [ ] Skip flags respected user preferences
- [ ] Minimal prompts (low interactivity)
- [ ] Clear next steps provided
- [ ] All URLs and paths displayed correctly

## Integration
- [ ] All created files are functional
- [ ] No conflicting configurations
- [ ] Documentation complete and accurate
- [ ] Project remains in valid operational state

## Final Validation
- [ ] direnv loads .envrc successfully
- [ ] Environment variables accessible in shell
- [ ] Logging directory structure correct
- [ ] CLI menu runs without errors
- [ ] Traefik proxy accessible (if configured)
