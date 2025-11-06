# Sync Standards Workflow Instructions

## Overview
This workflow silently syncs templates and standards from the agent_customisations repository to the current project. It operates with minimal prompts, only asking for confirmation when conflicts are detected.

## Execution Flow

<step n="1" goal="Verify sync source repository">
  <action>Check for agent_customisations repository:</action>
  <action>- Expected location: ~/Documents/raegis_labs/agent_customisations</action>
  <action>- Verify templates/ directory structure exists</action>

  <check if="repository not found">
    <action>Display error: "❌ agent_customisations repository not found at ~/Documents/raegis_labs/agent_customisations"</action>
    <action>Display: "   Cannot sync templates. Please clone the repository first."</action>
    <action>Exit workflow with error</action>
  </check>

  <check if="repository found">
    <action>Set: sync_source_root = ~/Documents/raegis_labs/agent_customisations</action>
    <action>Initialize sync tracking variables:</action>
    <action>- synced_count = 0</action>
    <action>- conflict_count = 0</action>
    <action>- skipped_count = 0</action>
    <action>- error_count = 0</action>
  </check>
</step>

<step n="2" goal="Sync agent templates" if="skip_agents == false">
  <action>Compare agent templates:</action>
  <action>- Source: {sync_source_root}/templates/agents/*.yaml</action>
  <action>- Target: {project-root}/.bmad/expansion-packs/*/agents/*.yaml</action>
  <action>- Protected: *.local.* files (never overwrite)</action>

  <action>For each template agent:</action>
  <action>1. Check if target file exists</action>
  <action>2. If exists, compare content (md5 hash or diff)</action>
  <action>3. If local modifications detected AND not *.local.*, mark as conflict</action>
  <action>4. If no conflict or new file, sync automatically</action>

  <check if="conflicts detected AND force == false">
    <action>Display conflict summary:</action>
    <action>
      "⚠️  Agent Template Conflicts Detected ({count} files)

      The following agents have local modifications:
      {list of agents with modification status}

      Options:
      1. Overwrite - Replace local versions with templates
      2. Keep - Preserve local modifications (skip sync)
      3. Review - Show diffs for each file

      Recommendation: Use *.local.yaml files for local customizations.
      "
    </action>

    <ask response="conflict_action">How to handle conflicts? [overwrite/keep/review]</ask>

    <check if="user answered overwrite">
      <action>Overwrite conflicted agents with templates</action>
      <action>Increment synced_count by conflict_count</action>
    </check>

    <check if="user answered keep">
      <action>Skip conflicted agents</action>
      <action>Increment skipped_count by conflict_count</action>
    </check>

    <check if="user answered review">
      <action>For each conflicted agent:</action>
      <action>- Show diff (template vs local)</action>
      <action>- Ask: overwrite/keep for this file</action>
      <action>- Apply user choice</action>
    </check>
  </check>

  <check if="force == true">
    <action>Overwrite all files (including conflicts) without prompting</action>
  </check>

  <action>Sync non-conflicted agents automatically (silent)</action>
  <action>Update symlinks in agents/ directory</action>
  <action>Increment synced_count by number of synced agents</action>
</step>

<step n="3" goal="Sync BMAD framework" if="skip_bmad == false">
  <action>Compare BMAD framework files:</action>
  <action>- Source: {sync_source_root}/templates/bmad/</action>
  <action>- Target: {project-root}/bmad/</action>
  <action>- Protected: bmad/bmm/config.yaml, bmad/*/config.yaml (never overwrite)</action>

  <action>For each BMAD file:</action>
  <action>1. Skip if protected (config.yaml files)</action>
  <action>2. Check if target exists and has local modifications</action>
  <action>3. If conflict, mark for review</action>
  <action>4. If no conflict, sync automatically</action>

  <check if="conflicts detected AND force == false">
    <action>Display conflict summary:</action>
    <action>
      "⚠️  BMAD Framework Conflicts Detected ({count} files)

      The following BMAD files have local modifications:
      {list of files with modification status}

      Options:
      1. Overwrite - Replace with templates
      2. Keep - Preserve local modifications
      "
    </action>

    <ask response="bmad_conflict_action">How to handle BMAD conflicts? [overwrite/keep]</ask>

    <check if="user answered overwrite">
      <action>Overwrite conflicted files</action>
      <action>Increment synced_count</action>
    </check>

    <check if="user answered keep">
      <action>Skip conflicted files</action>
      <action>Increment skipped_count</action>
    </check>
  </check>

  <check if="force == true">
    <action>Overwrite all non-protected files without prompting</action>
  </check>

  <action>Sync non-conflicted BMAD files automatically (silent)</action>
  <action>Increment synced_count</action>
</step>

<step n="4" goal="Sync workflow templates" if="skip_workflows == false">
  <action>Compare workflow templates:</action>
  <action>- Source: {sync_source_root}/templates/workflows/</action>
  <action>- Target: {project-root}/.bmad/expansion-packs/*/workflows/</action>
  <action>- Protected: *.local.* files</action>

  <action>For each workflow template:</action>
  <action>1. Check if target exists</action>
  <action>2. If exists and modified locally, mark as conflict</action>
  <action>3. If no conflict, sync automatically</action>

  <check if="conflicts detected AND force == false">
    <action>Display conflict summary:</action>
    <action>
      "⚠️  Workflow Template Conflicts Detected ({count} files)

      {list of workflows with local modifications}

      Recommendation: If you customized workflows, rename to *.local.* pattern.
      "
    </action>

    <ask response="workflow_conflict_action">Overwrite or keep local workflows? [overwrite/keep]</ask>

    <check if="user answered overwrite">
      <action>Overwrite conflicted workflows</action>
      <action>Increment synced_count</action>
    </check>

    <check if="user answered keep">
      <action>Skip conflicted workflows</action>
      <action>Increment skipped_count</action>
    </check>
  </check>

  <check if="force == true">
    <action>Overwrite all non-protected workflows</action>
  </check>

  <action>Sync non-conflicted workflows automatically (silent)</action>
  <action>Increment synced_count</action>
</step>

<step n="5" goal="Sync code style templates" if="skip_code_style == false">
  <action>Compare code style files:</action>
  <action>- Source: {sync_source_root}/templates/code-style/</action>
  <action>- Target: {project-root}/</action>
  <action>- Files: .editorconfig, .prettierrc, .eslintrc.js, pyproject.toml (style sections)</action>

  <action>For each code style file:</action>
  <action>1. Check if file exists in project</action>
  <action>2. If exists, compare with template</action>
  <action>3. For pyproject.toml: only sync [tool.black], [tool.ruff] sections</action>
  <action>4. If conflicts, mark for review</action>
  <action>5. If no conflict or new file, sync automatically</action>

  <check if="conflicts detected AND force == false">
    <action>Display conflict summary:</action>
    <action>
      "⚠️  Code Style Conflicts Detected ({count} files)

      {list of style files with differences}

      Note: For pyproject.toml, only [tool.black] and [tool.ruff] sections will be synced.
      "
    </action>

    <ask response="style_conflict_action">Apply template code style? [overwrite/keep]</ask>

    <check if="user answered overwrite">
      <action>Apply template style configurations</action>
      <action>For pyproject.toml: merge sections, don't replace entire file</action>
      <action>Increment synced_count</action>
    </check>

    <check if="user answered keep">
      <action>Keep local style configurations</action>
      <action>Increment skipped_count</action>
    </check>
  </check>

  <check if="force == true">
    <action>Apply all template style configurations</action>
  </check>

  <action>Sync non-conflicted style files automatically (silent)</action>
  <action>Increment synced_count</action>
</step>

<step n="6" goal="Sync GitHub workflow templates" if="skip_github == false">
  <action>Compare GitHub workflow files:</action>
  <action>- Source: {sync_source_root}/templates/github/workflows/</action>
  <action>- Target: {project-root}/.github/workflows/</action>
  <action>- Protected: *.local.yml files</action>

  <action>For each GitHub workflow template:</action>
  <action>1. Check if target exists</action>
  <action>2. If exists and modified, mark as conflict</action>
  <action>3. If no conflict, sync automatically</action>

  <check if="conflicts detected AND force == false">
    <action>Display conflict summary:</action>
    <action>
      "⚠️  GitHub Workflow Conflicts Detected ({count} files)

      {list of workflow files with modifications}

      Recommendation: For project-specific workflows, use *.local.yml naming.
      "
    </action>

    <ask response="github_conflict_action">Apply template GitHub workflows? [overwrite/keep]</ask>

    <check if="user answered overwrite">
      <action>Overwrite conflicted workflows</action>
      <action>Increment synced_count</action>
    </check>

    <check if="user answered keep">
      <action>Keep local workflows</action>
      <action>Increment skipped_count</action>
    </check>
  </check>

  <check if="force == true">
    <action>Overwrite all non-protected workflows</action>
  </check>

  <action>Sync non-conflicted GitHub workflows automatically (silent)</action>
  <action>Increment synced_count</action>
</step>

<step n="7" goal="Update sync log and generate summary">
  <action>Write sync log to {project-root}/.bmad/last-sync-standards.log:</action>
  <action>
    Timestamp: {ISO 8601 timestamp}
    Synced: {synced_count} files
    Conflicts: {conflict_count} files
    Skipped: {skipped_count} files
    Errors: {error_count} files

    Categories synced:
    - Agent templates: {agent_count} files
    - BMAD framework: {bmad_count} files
    - Workflows: {workflow_count} files
    - Code style: {style_count} files
    - GitHub workflows: {github_count} files
  </action>

  <action>Display summary (only non-zero counts shown):</action>
  <action>
    "✅ **Sync Standards Complete**

    {synced_count} files synced
    {conflict_count} conflicts handled
    {skipped_count} files skipped
    {error_count} errors (if any)

    Log: .bmad/last-sync-standards.log
    "
  </action>

  <check if="synced_count > 0">
    <action>Display next steps:</action>
    <action>
      "**Next Steps:**
      1. Review changes: git diff
      2. Run tests to verify compatibility
      3. Commit synced changes
      "
    </action>
  </check>

  <check if="synced_count == 0">
    <action>Display: "ℹ️  All templates are already in sync. No changes made."</action>
  </check>
</step>

## Notes for Implementation

### Silent Execution
- No category-by-category prompts (unlike maintenance workflow)
- Automatic sync when no conflicts
- Only prompt when conflicts detected
- Very low verbosity (errors and conflicts only)

### Conflict Detection
- Compare files using md5 hash or diff
- Conflict = file exists + local modifications + not protected pattern
- Protected patterns never trigger conflicts (automatically skipped)
- Threshold: if >10 conflicts, offer "skip all" to avoid prompt fatigue

### Protected Patterns
- `*.local.*` files - User's local customizations
- `config.yaml` files - Project-specific configuration
- `*.local.yml` - Local GitHub workflows

### Force Mode (`--force`)
- Overwrites all conflicts without prompting
- Still respects protected patterns
- Use with caution (recommend --dry-run first)

### Dry Run Mode (`--dry-run`)
- Shows what would be synced
- Shows conflicts that would be prompted
- Makes no actual changes
- Generates preview report

### Integration Points
- Uses same rsync logic as /init-sync-bmad command
- Respects .gitignore
- Preserves file permissions
- Uses git for conflict detection (if in git repo)

### Error Handling
- Missing source repository: error and exit
- Missing source files: warning and skip
- Permission errors: error and skip file
- Continue syncing remaining files even if some fail
- Record all errors in log file

### Performance
- Batch file comparisons for speed
- Use checksums to avoid unnecessary syncs
- Parallel sync where safe (independent files)

### Safety Features
- Protected files never overwritten
- Conflicts always detected
- Can preview with --dry-run
- All changes reversible via git
- Log file for audit trail
