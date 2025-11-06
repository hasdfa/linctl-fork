# Sync and Compile Workflow

Synchronize BMAD agents and workflows from the central template repository, then compile all agents to make them immediately available.

## Workflow Steps

### Step 0: Parse Flags and Initialize

**Parse command flags:**
- `--dry-run`: Preview changes without applying
- `--agents-only`: Sync only agents (skip workflows)
- `--workflows-only`: Sync only workflows (skip agents)
- `--skip-compile`: Sync but don't compile agents
- `--verbose`: Show detailed output

**Set verbosity:**
- Default: Silent execution with summary
- `--verbose`: Show all operations
- `--dry-run`: Always verbose

**Load configuration:**
```yaml
agent_customisations_repo: ~/Documents/raegis_labs/agent_customisations
templates_base_path: {agent_customisations_repo}/bmad
sync_modules:
  - bmm          # Custom BMM agents (Atlas, etc.)
  - bmb          # BMad Builder agent and workflows
  - bmm-plus     # Extended agents (PM, SM, TEA, Dev, Analyst, UX) and Linear workflows
  - project-ops  # Custom steward agents (docs, linear, test)
```

**For each module:**
```yaml
templates_agents_path: {templates_base_path}/{module}/agents
templates_workflows_path: {templates_base_path}/{module}/workflows
project_agents_path: {project-root}/bmad/{module}/agents
project_workflows_path: {project-root}/bmad/{module}/workflows
```

---

### Step 1: Verify Prerequisites

**Check template repository exists:**
```bash
if [[ ! -d "{agent_customisations_repo}" ]]; then
  ERROR: Template repository not found
  Expected: {agent_customisations_repo}
  Action: Clone agent_customisations repository first
fi
```

**Check project has BMAD v6 installed:**
```bash
if [[ ! -d "{project-root}/bmad/bmm" ]]; then
  ERROR: BMAD v6 not installed in this project
  Action: Run 'npx bmad-method@alpha install' first
fi
```

**Check for uncommitted changes (optional warning):**
```bash
if git status --porcelain | grep -q "^M.*bmad/"; then
  WARNING: You have uncommitted changes in bmad/ directory
  Consider committing before sync
  Continue? [y/N]
fi
```

---

### Step 2: Preview Changes

**Scan template repository (all modules):**
```bash
# For each module in sync_modules list
for module in bmm bmb bmm-plus project-ops; do
  TEMPLATES_AGENTS_YAML[$module]=$(find {templates_base_path}/$module/agents -name "*.yaml" -type f 2>/dev/null)
  TEMPLATES_AGENTS_MD[$module]=$(find {templates_base_path}/$module/agents -name "*.md" -type f 2>/dev/null)
  TEMPLATES_WORKFLOWS[$module]=$(find {templates_base_path}/$module/workflows -type d -name "workflow.yaml" -exec dirname {} \; 2>/dev/null)
done
```

**Compare with current project (all modules):**
```bash
# For each module
for module in bmm bmb bmm-plus project-ops; do
  echo "Checking module: $module"

  # Check MD files (pre-compiled agents)
  for agent in ${TEMPLATES_AGENTS_MD[$module]}; do
    AGENT_NAME=$(basename "$agent")
    TARGET="{project-root}/bmad/$module/agents/$AGENT_NAME"

    if [[ -f "$TARGET" ]]; then
      if diff -q "$agent" "$TARGET" >/dev/null; then
        STATUS: [$module] Unchanged
      else
        STATUS: [$module] Will UPDATE
      fi
    else
      STATUS: [$module] Will ADD (new)
    fi
  done

  # Check YAML files (source reference)
  for agent in ${TEMPLATES_AGENTS_YAML[$module]}; do
    AGENT_NAME=$(basename "$agent")
    TARGET="{project-root}/bmad/$module/agents/$AGENT_NAME"

    if [[ -f "$TARGET" ]]; then
      if diff -q "$agent" "$TARGET" >/dev/null; then
        STATUS: [$module] Unchanged
      else
        STATUS: [$module] Will UPDATE
      fi
    else
      STATUS: [$module] Will ADD (new)
    fi
  done
done
```

**Show preview:**
```
=== Sync Preview ===

Module: bmm
  Agents:
    ✓ atlas.md (unchanged, pre-compiled)
    ✓ atlas.yaml (unchanged, source)
    → architect.md (will update)
  Workflows:
    ✓ quick-setup/ (unchanged)
    → maintenance/ (will update)

Module: project-ops
  Agents:
    + docs-steward.md (will add, pre-compiled)
    + docs-steward.yaml (will add, source)
    + linear-steward.md (will add, pre-compiled)
    + test-steward.md (will add, pre-compiled)
  Workflows:
    + example-workflow/ (will add)

Note: Pre-compiled .md files copied directly (no compilation needed)

Continue? [y/N]
```

If `--dry-run`: STOP HERE and show preview only

---

### Step 3: Sync Agents (All Modules)

**If `--workflows-only` flag: SKIP this step**

**Copy pre-compiled MD files (ready to use):**
```bash
# For each module
for module in bmm project-ops; do
  for agent in ${TEMPLATES_AGENTS_MD[$module]}; do
    AGENT_NAME=$(basename "$agent")
    TARGET="{project-root}/bmad/$module/agents/$AGENT_NAME"

    # Skip *.local.* files (project-specific customizations)
    if [[ "$AGENT_NAME" =~ \.local\. ]]; then
      SKIP: $AGENT_NAME (local customization)
      continue
    fi

    # Create target directory if needed
    mkdir -p "{project-root}/bmad/$module/agents"

    # Copy pre-compiled agent to bmad/{module}/agents/
    cp "$agent" "$TARGET"

    # Also copy to Claude Code GLOBAL commands directory (available in all projects)
    CLAUDE_GLOBAL="$HOME/.claude/commands/bmad/$module/agents/$AGENT_NAME"
    mkdir -p "$(dirname "$CLAUDE_GLOBAL")"
    cp "$agent" "$CLAUDE_GLOBAL"

    # Also copy to Codex GLOBAL prompts directory (flat structure)
    if [[ -d "$HOME/.codex" ]]; then
      # Codex uses flat file structure: bmad-{module}-agents-{name}.md
      CODEX_FILENAME="bmad-${module}-agents-${AGENT_NAME}"
      CODEX_GLOBAL="$HOME/.codex/prompts/$CODEX_FILENAME"
      cp "$agent" "$CODEX_GLOBAL"
    fi

    if [[ $VERBOSE == true ]]; then
      echo "✓ Synced: [$module] $AGENT_NAME (project + global Claude Code + global Codex)"
    fi
  done
done
```

**Copy YAML source files (optional reference):**
```bash
# For each module
for module in bmm project-ops; do
  for agent in ${TEMPLATES_AGENTS_YAML[$module]}; do
    AGENT_NAME=$(basename "$agent")
    TARGET="{project-root}/bmad/$module/agents/$AGENT_NAME"

    # Skip *.local.* files (project-specific customizations)
    if [[ "$AGENT_NAME" =~ \.local\. ]]; then
      SKIP: $AGENT_NAME (local customization)
      continue
    fi

    # Copy source YAML
    cp "$agent" "$TARGET"

    if [[ $VERBOSE == true ]]; then
      echo "✓ Synced: [$module] $AGENT_NAME (source)"
    fi
  done
done
```

**Preserve local customizations:**
```bash
# *.local.* files are NEVER overwritten
# Example: architect.local.yaml and architect.local.md stay untouched
```

---

### Step 4: Sync Workflows (All Modules)

**If `--agents-only` flag: SKIP this step**

**Copy workflow directories:**
```bash
# For each module
for module in bmm project-ops; do
  for workflow_dir in ${TEMPLATES_WORKFLOWS[$module]}; do
    WORKFLOW_NAME=$(basename "$workflow_dir")
    TARGET="{project-root}/bmad/$module/workflows/$WORKFLOW_NAME"

    # Create target directory if needed
    mkdir -p "$TARGET"

    # Copy all files (workflow.yaml, instructions.md, checklist.md)
    cp -r "$workflow_dir/"* "$TARGET/"

    if [[ $VERBOSE == true ]]; then
      echo "✓ Synced workflow: [$module] $WORKFLOW_NAME"
    fi
  done
done
```

---

### Step 5: Sync Workflow Commands to User-Level (Claude Code & Codex)

**Sync slash commands from template repository to user-level directories:**

```bash
echo "📋 Syncing workflow commands to user-level..."

TEMPLATE_BASE="{agent_customisations_repo}/.claude/commands/bmad"
CLAUDE_USER_BASE="$HOME/.claude/commands/bmad"
CODEX_USER_BASE="$HOME/.codex/commands"

# For each module
for module in bmm bmb bmm-plus project-ops; do
  # Check if template has commands for this module
  if [[ -d "$TEMPLATE_BASE/$module" ]]; then

    # Sync to Claude Code user-level
    echo "  Syncing $module commands to Claude Code..."
    mkdir -p "$CLAUDE_USER_BASE/$module"
    cp -r "$TEMPLATE_BASE/$module/"* "$CLAUDE_USER_BASE/$module/"

    # Count commands
    CMD_COUNT=$(find "$CLAUDE_USER_BASE/$module" -name "*.md" -type f | wc -l)

    if [[ $VERBOSE == true ]]; then
      echo "✓ [$module] $CMD_COUNT command(s) synced to ~/.claude/commands/bmad/$module/"
    fi

    # Sync to Codex user-level (if .codex directory exists)
    if [[ -d "$HOME/.codex" ]]; then
      echo "  Syncing $module commands to Codex..."
      mkdir -p "$CODEX_USER_BASE"

      # Codex uses flat file structure: bmad-{module}-{path}-{name}.md
      # Convert .claude hierarchical structure to Codex flat structure
      find "$TEMPLATE_BASE/$module" -name "*.md" -type f | while read cmd_file; do
        # Get relative path from module directory
        REL_PATH="${cmd_file#$TEMPLATE_BASE/$module/}"
        # Convert path separators to dashes: workflows/linear-create-epic.md → workflows-linear-create-epic.md
        FLAT_NAME="${REL_PATH//\//-}"
        # Create Codex filename: bmad-bmm-plus-workflows-linear-create-epic.md
        CODEX_FILENAME="bmad-${module}-${FLAT_NAME}"

        cp "$cmd_file" "$CODEX_USER_BASE/$CODEX_FILENAME"
      done

      if [[ $VERBOSE == true ]]; then
        echo "✓ [$module] Commands synced to ~/.codex/commands/ (flat structure)"
      fi
    fi
  fi
done

if [[ $VERBOSE == false ]]; then
  echo "✓ Workflow commands synced to user-level (Claude Code + Codex)"
fi
```

**What this does:**
- Copies workflow slash commands from template repository
- Claude Code: Maintains directory structure in `~/.claude/commands/bmad/`
- Codex: Flattens to `~/.codex/commands/bmad-{module}-{path}-{name}.md`
- Makes commands available across ALL projects immediately

---

### Step 6: Verify Agent Files (All Modules)

**Check that .md files were copied:**

Since we copy pre-compiled `.md` files directly from templates, no compilation is needed.

```bash
echo "📋 Verifying agents..."

TOTAL_MD_COUNT=0

# For each module
for module in bmm project-ops; do
  MODULE_MD_COUNT=0

  for md_file in {project-root}/bmad/$module/agents/*.md; do
    if [[ -f "$md_file" ]]; then
      MODULE_MD_COUNT=$((MODULE_MD_COUNT + 1))
      TOTAL_MD_COUNT=$((TOTAL_MD_COUNT + 1))

      if [[ $VERBOSE == true ]]; then
        echo "✓ [$module] $(basename "$md_file") ready"
      fi
    fi
  done

  if [[ $MODULE_MD_COUNT -gt 0 && $VERBOSE == false ]]; then
    echo "✓ [$module] $MODULE_MD_COUNT agent(s) ready"
  fi
done

if [[ $TOTAL_MD_COUNT -eq 0 ]]; then
  echo "⚠ Warning: No agent .md files found"
  echo "Check template repository has pre-compiled agents"
else
  echo "✓ Total: $TOTAL_MD_COUNT agent(s) across all modules"
fi
```

**Note on compilation:**
- Pre-compiled `.md` files copied from templates
- No BMAD compiler needed for custom agents (like Atlas)
- Agents ready to use immediately after sync

---

### Step 7: Final Summary

**Silent mode (default):**
```
✓ Sync complete
  • {N} agents synced (pre-compiled .md files)
  • {M} workflows synced
  • Workflow commands synced to user-level
  • Global commands updated (Claude Code + Codex)
  • Available in ALL projects immediately
```

**Verbose mode:**
```
=== Sync Summary ===

Module: bmm
  Agents synced:
    ✓ atlas.md (pre-compiled, ready to use)
    ✓ atlas.yaml (source reference)
    ✓ architect.md (from templates)
  Workflows synced:
    ✓ quick-setup/
    ✓ maintenance/
    ✓ sync-and-compile/

Module: project-ops
  Agents synced:
    ✓ docs-steward.md (pre-compiled, ready to use)
    ✓ linear-steward.md (pre-compiled, ready to use)
    ✓ test-steward.md (pre-compiled, ready to use)
  Workflows synced:
    ✓ example-workflow/

Status:
  ✓ All agents ready (pre-compiled .md files)
  ✓ No compilation needed
  ✓ {K} agents available across {M} modules
  ✓ Workflow commands synced to user-level
  ✓ Slash commands available in Claude Code and Codex

Next steps:
  • Claude Code agents: @atlas or /bmad:bmm:agents:atlas
  • Claude Code workflows: /bmad:bmm-plus:workflows:linear-create-epic
  • Codex agents: /prompts:bmad-bmm-agents-atlas
  • Codex workflows: /bmad-bmm-plus-workflows-linear-create-epic
  • Test workflows: @atlas → *quick-setup or /bmad:bmm:workflows:quick-setup
```

---

## Error Handling

**Template repository not found:**
```
ERROR: Template repository not found
Expected: ~/Documents/raegis_labs/agent_customisations
Solution: Clone the repository or update path in workflow.yaml
```

**BMAD not installed:**
```
ERROR: BMAD v6 not detected
Expected: bmad/bmm/ directory
Solution: Run 'npx bmad-method@alpha install' first
```

**Missing agent files:**
```
ERROR: No .md files found in template repository
Expected: ~/Documents/raegis_labs/agent_customisations/bmad/bmm/agents/*.md
Solution: Ensure template repository has pre-compiled agents
```

**Permission denied:**
```
ERROR: Cannot write to bmad/bmm/agents/
Solution: Check directory permissions
```

---

## Flags Reference

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview changes without applying |
| `--agents-only` | Sync only agents (skip workflows) |
| `--workflows-only` | Sync only workflows (skip agents) |
| `--verbose` | Show detailed operation output |

---

## Examples

**Full sync (agents + workflows):**
```bash
@atlas
*sync-and-compile
# Pre-compiled agents ready immediately
```

**Preview changes first:**
```bash
@atlas
*sync-and-compile --dry-run
```

**Sync only agents:**
```bash
@atlas
*sync-and-compile --agents-only
# Only updates agent .md and .yaml files
```

**Sync only workflows:**
```bash
@atlas
*sync-and-compile --workflows-only
# Only updates workflow directories
```

---

## Notes

**Local Customizations:**
- `*.local.*` files are NEVER overwritten
- Example: `architect.local.md` and `architect.local.yaml` preserved
- Use local versions for project-specific tweaks

**Pre-Compiled Agents:**
- Template repository has pre-compiled `.md` files
- Sync copies `.md` files directly (ready to use)
- Also copies `.yaml` files (optional source reference)
- No BMAD compiler needed for custom agents like Atlas

**Template Repository:**
- Single source of truth: `agent_customisations` repo
- Updates flow: Templates → Project
- Never reverse: Project → Templates (use manual promotion)
- Maintains both `.yaml` (source) and `.md` (compiled) versions
