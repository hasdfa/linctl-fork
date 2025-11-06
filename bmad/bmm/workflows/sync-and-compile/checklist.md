# Sync and Compile - Validation Checklist

Use this checklist to verify successful sync and compilation.

## Prerequisites Validation

### Template Repository
- [ ] Template repository exists at `~/Documents/raegis_labs/agent_customisations`
- [ ] Repository has `bmad/bmm/agents/` directory
- [ ] Repository has `bmad/bmm/workflows/` directory
- [ ] Repository is up-to-date (optional: `git pull`)

### Project Setup
- [ ] Current project has `bmad/` directory
- [ ] BMAD v6 structure present (`bmad/bmm/`, `bmad/core/`)
- [ ] Project has `bmad/bmm/agents/` directory
- [ ] Project has `bmad/bmm/workflows/` directory

### Environment
- [ ] Node.js >= 20.0.0 installed (for compilation)
- [ ] Write permissions on project `bmad/` directory
- [ ] No file locks on agent/workflow files

---

## Sync Validation

### Agent Files Synced
- [ ] All `.agent.yaml` files copied from templates
- [ ] File timestamps updated to current time
- [ ] File contents match template source exactly
- [ ] No permission errors during copy
- [ ] Local customizations (`.local.yaml`) preserved (not overwritten)

**Verify:**
```bash
# Check agents were synced
ls -lt bmad/bmm/agents/*.yaml | head -5

# Compare with templates
diff bmad/bmm/agents/atlas.agent.yaml \
     ~/Documents/raegis_labs/agent_customisations/bmad/bmm/agents/atlas.agent.yaml
```

### Workflow Directories Synced
- [ ] All workflow directories copied from templates
- [ ] Each workflow has `workflow.yaml`
- [ ] Each workflow has `instructions.md`
- [ ] Each workflow has `checklist.md` (if applicable)
- [ ] Workflow file contents match template source

**Verify:**
```bash
# Check workflows synced
ls -la bmad/bmm/workflows/

# Verify workflow structure
ls bmad/bmm/workflows/quick-setup/
# Should see: workflow.yaml, instructions.md, checklist.md
```

---

## Compilation Validation

### Compilation Executed
- [ ] `npx bmad-method@alpha install` was run
- [ ] "Compile Agents" option was selected
- [ ] Compilation completed without errors
- [ ] No warnings about missing files or syntax errors

### Compiled Files Created
- [ ] Every `.agent.yaml` has corresponding `.md` file
- [ ] `.md` files have recent timestamps (compilation time)
- [ ] `.md` files contain agent content (not empty)
- [ ] `.md` file size is reasonable (>1KB typically)

**Verify:**
```bash
# Check .md files exist
ls -la bmad/bmm/agents/*.md

# Verify each YAML has MD
for yaml in bmad/bmm/agents/*.agent.yaml; do
  base=$(basename "$yaml" .agent.yaml)
  md="bmad/bmm/agents/${base}.md"
  if [[ ! -f "$md" ]]; then
    echo "Missing: $md"
  else
    echo "OK: $base"
  fi
done

# Check .md file size (should be >1KB)
ls -lh bmad/bmm/agents/*.md | awk '{if ($5 ~ /^[0-9]+B$/) print "WARNING: " $9 " is too small"}'
```

### Compiled Content Quality
- [ ] `.md` files are valid Markdown
- [ ] Agent metadata present (name, title, icon)
- [ ] Menu items rendered correctly
- [ ] Workflow paths preserved
- [ ] No XML/YAML syntax errors in output

**Verify:**
```bash
# Check atlas.md has expected content
grep -q "name: Atlas" bmad/bmm/agents/atlas.md && echo "✓ Atlas metadata OK"
grep -q "*quick-setup" bmad/bmm/agents/atlas.md && echo "✓ Atlas menu OK"

# Check for compilation errors
grep -i "error\|warning" bmad/bmm/agents/*.md
# Should return nothing
```

---

## Functional Validation

### Agents Available in Menu
- [ ] Agents appear in AI tool agent selector
- [ ] Agent names display correctly
- [ ] Agent icons render properly (🌍, 🏗️, etc.)
- [ ] Agent descriptions are readable

**Verify in Claude Code / Codex:**
```bash
# In chat, type: @
# Should see: @atlas, @architect, @developer, etc.
```

### Agents Activate Successfully
- [ ] `@atlas` activates Atlas agent
- [ ] Agent displays menu with all commands
- [ ] Agent loads configuration from `bmad/core/config.yaml`
- [ ] Agent greets with correct user name

**Test:**
```bash
# In chat:
@atlas

# Expected output:
# Greetings, {user_name}! 🌍
# [Atlas menu with 8+ commands]
```

### Workflows Execute
- [ ] `*quick-setup` workflow can be selected
- [ ] Workflow loads configuration from `workflow.yaml`
- [ ] Workflow reads `instructions.md`
- [ ] No "workflow not found" errors

**Test:**
```bash
# In chat:
@atlas
*sync-and-compile --dry-run

# Should execute without errors
```

---

## Edge Case Validation

### Local Customizations Protected
- [ ] Files matching `*.local.*` pattern were NOT overwritten
- [ ] Example: `architect.local.yaml` still has project-specific changes
- [ ] Local symlinks still point to correct files

**Verify:**
```bash
# If you had local customizations, check they're preserved
if [[ -f bmad/bmm/agents/architect.local.yaml ]]; then
  # Check it wasn't overwritten (timestamp should be old)
  ls -l bmad/bmm/agents/architect.local.yaml
fi
```

### No Unexpected Files
- [ ] No `.bak` or backup files created
- [ ] No temporary files left behind
- [ ] No duplicate files (e.g., `atlas.agent.yaml.1`)

**Verify:**
```bash
# Check for unexpected files
find bmad/bmm/ -name "*.bak" -o -name "*~" -o -name "*.tmp"
# Should return nothing
```

### Git Status Clean (Optional)
- [ ] Only expected files changed (agents, workflows)
- [ ] No unintended modifications
- [ ] Changes ready to commit if desired

**Verify:**
```bash
git status bmad/

# Should show:
# Modified: bmad/bmm/agents/*.yaml
# Modified: bmad/bmm/agents/*.md
# Modified: bmad/bmm/workflows/*/
```

---

## Rollback Validation

### Can Revert If Needed
- [ ] Git can restore previous version
- [ ] Backup exists (if created)
- [ ] Know how to re-run sync safely

**Rollback if needed:**
```bash
# Undo changes
git restore bmad/

# Or restore specific file
git restore bmad/bmm/agents/atlas.agent.yaml

# Re-compile if restored
npx bmad-method@alpha install → Compile Agents
```

---

## Success Criteria

All of the following must be true:

1. ✅ **Sync Complete**
   - All template agents copied to project
   - All template workflows copied to project
   - Local customizations preserved

2. ✅ **Compilation Complete**
   - Every `.yaml` has matching `.md` file
   - `.md` files have valid content
   - No compilation errors

3. ✅ **Agents Functional**
   - Agents appear in selector
   - Agents activate correctly
   - Workflows execute without errors

4. ✅ **No Regressions**
   - Existing functionality still works
   - Local customizations intact
   - No unexpected files created

---

## Common Issues and Fixes

### Issue: Agents not appearing after sync
**Cause:** Compilation not run or failed
**Fix:**
```bash
npx bmad-method@alpha install
→ Compile Agents
```

### Issue: "Workflow not found" error
**Cause:** Workflow directory not synced
**Fix:**
```bash
# Re-run sync with verbose
@atlas
*sync-and-compile --verbose
```

### Issue: Local customizations overwritten
**Cause:** Files didn't match `*.local.*` pattern
**Fix:**
```bash
# Restore from git
git restore bmad/bmm/agents/your-file.yaml

# Rename to protect
mv your-file.yaml your-file.local.yaml
```

### Issue: Compilation errors
**Cause:** Invalid YAML syntax in agent file
**Fix:**
```bash
# Check YAML syntax
yamllint bmad/bmm/agents/*.yaml

# Or validate online: yamllint.com
```

---

## Post-Validation Actions

After all checks pass:

1. **Test Key Workflows:**
   ```bash
   @atlas
   *quick-setup --dry-run
   *health-check
   ```

2. **Commit Changes (Optional):**
   ```bash
   git add bmad/
   git commit -m "Sync BMAD agents and workflows from templates"
   ```

3. **Document Custom Changes:**
   - Note any project-specific agents not in templates
   - Document local customizations (`.local.yaml` files)
   - Update project README if needed

4. **Share Success:**
   - Update team if collaborative project
   - Note any new agents/workflows available
   - Document any breaking changes

---

## Validation Complete

If all checks pass:

✅ **Sync and compile successful**
- Templates synced to project
- Agents compiled and functional
- Workflows ready to use
- No regressions or issues

Your project now has the latest BMAD agents and workflows from the central template repository.
