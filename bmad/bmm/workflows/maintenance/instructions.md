# Maintenance Workflow Instructions

## Overview
This workflow performs project maintenance through cleanup and sync operations. It shows previews for each category and asks for confirmation before making changes, ensuring you maintain control while keeping your project healthy.

## Execution Flow

<step n="1" goal="Check maintenance schedule and explain scope">
  <action>Read last maintenance timestamp from {project-root}/.bmad/last-maintenance</action>

  <check if="timestamp file exists">
    <action>Calculate days since last maintenance</action>

    <check if="days > 30">
      <action>Display: "⚠️  Maintenance is {days} days overdue (threshold: 30 days)"</action>
    </check>

    <check if="days <= 30">
      <action>Display: "✅ Last maintenance was {days} days ago (within 30-day threshold)"</action>
    </check>
  </check>

  <check if="timestamp file does not exist">
    <action>Display: "ℹ️  No previous maintenance recorded. This is your first maintenance run."</action>
  </check>

  <action>Display maintenance scope:</action>
  <action>
    "This workflow will:

    **Cleanup Operations:**
    1. Dead Code - Remove unused imports, commented code, unreachable code
    2. Old Branches - Clean up merged or stale git branches
    3. Stale Dependencies - Remove outdated or unused packages
    4. Old Logs - Archive or delete logs older than 30 days
    5. Temp Files - Remove .DS_Store, *.tmp, __pycache__, caches

    **Sync Operations:**
    6. Agent Templates - Update from agent_customisations repository
    7. BMAD Framework - Check for and apply BMAD updates
    8. Linear Settings - Sync labels, states, and templates
    9. Foundation Docs - Update idea.md, CLAUDE.md templates

    You'll be asked to confirm each category before changes are made."
  </action>
</step>

<step n="2" goal="Initialize tracking variables">
  <action>Set: cleanup_summary = []</action>
  <action>Set: sync_summary = []</action>
  <action>Set: total_changes = 0</action>
  <action>Set: categories_processed = 0</action>
</step>

<!-- CLEANUP OPERATIONS -->

<step n="3" goal="Cleanup: Dead Code" if="skip_dead_code == false AND skip_cleanup == false">
  <action>Display: "### 1/9: Dead Code Cleanup"</action>

  <action>Scan for dead code:</action>
  <action>- Unused imports (Python: unused names in import statements, JS/TS: unused imports)</action>
  <action>- Commented code blocks (// TODO, /* ... */, # commented functions)</action>
  <action>- Unreachable code (code after return/break, if(false) blocks)</action>
  <action>- Unused functions/variables (based on usage analysis)</action>

  <action>Build list of dead code findings with file:line references</action>

  <check if="findings found">
    <template-output path="{project-root}/.bmad/maintenance-preview-dead-code.md">
      # Dead Code Cleanup Preview

      Found {count} instances of dead code:

      ## Unused Imports ({count})
      {list of files with unused imports}

      ## Commented Code ({count})
      {list of commented blocks with context}

      ## Unreachable Code ({count})
      {list of unreachable code locations}

      ## Unused Functions/Variables ({count})
      {list of unused definitions}

      **Impact**: Removing this code will improve readability and reduce maintenance burden.
      **Risk**: Low - commented and unused code doesn't affect runtime behavior.
    </template-output>

    <ask response="confirm_dead_code">Apply dead code cleanup? [y/n/skip]</ask>

    <check if="user answered y or yes">
      <action>Execute dead code removal (preserve git history)</action>
      <action>Add to cleanup_summary: "Dead code: Removed {count} instances"</action>
      <action>Increment total_changes by {count}</action>
    </check>

    <check if="user answered skip">
      <action>Add to cleanup_summary: "Dead code: Skipped by user"</action>
    </check>
  </check>

  <check if="no findings">
    <action>Display: "✅ No dead code found"</action>
    <action>Add to cleanup_summary: "Dead code: None found"</action>
  </check>

  <action>Increment categories_processed</action>
</step>

<step n="4" goal="Cleanup: Old Branches" if="skip_branches == false AND skip_cleanup == false">
  <action>Display: "### 2/9: Old Branch Cleanup"</action>

  <action>Scan for branches to clean:</action>
  <action>- List all local branches (git branch)</action>
  <action>- Identify merged branches (git branch --merged)</action>
  <action>- Identify stale branches (no commits in >90 days)</action>
  <action>- Exclude current branch and main/master</action>

  <action>Build list of branches to remove</action>

  <check if="branches found">
    <template-output path="{project-root}/.bmad/maintenance-preview-branches.md">
      # Branch Cleanup Preview

      Found {count} branches to clean:

      ## Merged Branches ({count})
      {list of merged branches with last commit date}
      - These branches have been merged to main/master

      ## Stale Branches ({count})
      {list of stale branches with last commit date}
      - No commits in the last 90 days

      **Impact**: Cleaner branch list, easier navigation.
      **Risk**: Low - merged work is in main, can be recovered from reflog.

      **Note**: Remote branches will not be deleted (use git push origin --delete manually if needed).
    </template-output>

    <ask response="confirm_branches">Delete these local branches? [y/n/skip]</ask>

    <check if="user answered y or yes">
      <action>Delete branches: git branch -d {branch_name} for each</action>
      <action>Add to cleanup_summary: "Branches: Deleted {count} branches"</action>
      <action>Increment total_changes by {count}</action>
    </check>

    <check if="user answered skip">
      <action>Add to cleanup_summary: "Branches: Skipped by user"</action>
    </check>
  </check>

  <check if="no branches found">
    <action>Display: "✅ No old branches to clean"</action>
    <action>Add to cleanup_summary: "Branches: None found"</action>
  </check>

  <action>Increment categories_processed</action>
</step>

<step n="5" goal="Cleanup: Stale Dependencies" if="skip_deps == false AND skip_cleanup == false">
  <action>Display: "### 3/9: Stale Dependencies Cleanup"</action>

  <action>Detect package managers:</action>
  <action>- Node.js: Check for package.json</action>
  <action>- Python: Check for requirements.txt or pyproject.toml</action>

  <check if="Node.js project">
    <action>Run: npm outdated --json (or yarn outdated)</action>
    <action>Identify packages with major version updates available</action>
    <action>Check for unused dependencies: npx depcheck</action>
  </check>

  <check if="Python project">
    <action>Run: pip list --outdated --format=json</action>
    <action>Identify packages with updates available</action>
  </check>

  <action>Build list of stale/unused dependencies</action>

  <check if="dependencies found">
    <template-output path="{project-root}/.bmad/maintenance-preview-deps.md">
      # Dependencies Cleanup Preview

      ## Outdated Packages ({count})
      {list of packages with current vs latest version}

      ## Unused Packages ({count})
      {list of packages not imported anywhere}

      **Impact**: Better security, performance, smaller bundle size.
      **Risk**: Medium - updates may have breaking changes. Review changelogs before updating.

      **Recommended Actions:**
      - Update patch/minor versions (low risk)
      - Review major version updates carefully
      - Remove unused packages (low risk)
    </template-output>

    <ask response="confirm_deps">Update outdated and remove unused dependencies? [y/n/skip]</ask>

    <check if="user answered y or yes">
      <action>Remove unused dependencies first</action>
      <action>Update outdated dependencies (patch and minor versions)</action>
      <action>Display: "⚠️  Major version updates require manual review. See preview for details."</action>
      <action>Add to cleanup_summary: "Dependencies: Updated {count} packages, removed {unused_count}"</action>
      <action>Increment total_changes by {count}</action>
    </check>

    <check if="user answered skip">
      <action>Add to cleanup_summary: "Dependencies: Skipped by user"</action>
    </check>
  </check>

  <check if="no dependencies found">
    <action>Display: "✅ All dependencies are up to date"</action>
    <action>Add to cleanup_summary: "Dependencies: All current"</action>
  </check>

  <action>Increment categories_processed</action>
</step>

<step n="6" goal="Cleanup: Old Logs" if="skip_logs == false AND skip_cleanup == false">
  <action>Display: "### 4/9: Old Logs Cleanup"</action>

  <action>Scan for log files:</action>
  <action>- Check .logs/ directory for files older than 30 days</action>
  <action>- Check for rotated logs (*.log.1, *.log.2, etc.)</action>
  <action>- Calculate total size of old logs</action>

  <action>Build list of log files to archive/delete</action>

  <check if="old logs found">
    <template-output path="{project-root}/.bmad/maintenance-preview-logs.md">
      # Old Logs Cleanup Preview

      Found {count} old log files ({total_size}):

      {list of log files with age and size}

      **Options:**
      1. Delete - Permanently remove old logs
      2. Archive - Compress and move to .logs/archive/

      **Impact**: Free disk space ({total_size}).
      **Risk**: Low - logs older than 30 days rarely needed.
    </template-output>

    <ask response="confirm_logs">Archive old logs? [y=archive/d=delete/n=keep/skip]</ask>

    <check if="user answered y or yes">
      <action>Create .logs/archive/ directory</action>
      <action>Compress logs by month: tar -czf logs-YYYY-MM.tar.gz</action>
      <action>Move compressed archives to .logs/archive/</action>
      <action>Add to cleanup_summary: "Logs: Archived {count} files ({total_size})"</action>
      <action>Increment total_changes by {count}</action>
    </check>

    <check if="user answered d or delete">
      <action>Delete old log files</action>
      <action>Add to cleanup_summary: "Logs: Deleted {count} files ({total_size})"</action>
      <action>Increment total_changes by {count}</action>
    </check>

    <check if="user answered skip or n">
      <action>Add to cleanup_summary: "Logs: Kept by user"</action>
    </check>
  </check>

  <check if="no old logs">
    <action>Display: "✅ No old logs to clean (all within 30 days)"</action>
    <action>Add to cleanup_summary: "Logs: None older than 30 days"</action>
  </check>

  <action>Increment categories_processed</action>
</step>

<step n="7" goal="Cleanup: Temp Files" if="skip_temp == false AND skip_cleanup == false">
  <action>Display: "### 5/9: Temp Files Cleanup"</action>

  <action>Scan for temporary files:</action>
  <action>- .DS_Store (macOS)</action>
  <action>- Thumbs.db (Windows)</action>
  <action>- *.tmp, *.temp, *.cache</action>
  <action>- __pycache__/ directories</action>
  <action>- node_modules/.cache/</action>
  <action>- .pytest_cache/, .mypy_cache/</action>
  <action>- Build artifacts in dist/, build/ if not in .gitignore</action>

  <action>Build list of temp files/directories</action>

  <check if="temp files found">
    <template-output path="{project-root}/.bmad/maintenance-preview-temp.md">
      # Temporary Files Cleanup Preview

      Found {count} temporary items ({total_size}):

      ## System Files ({count})
      {.DS_Store, Thumbs.db, etc.}

      ## Cache Directories ({count})
      {__pycache__, .cache, etc.}

      ## Build Artifacts ({count})
      {dist/, build/, etc.}

      **Impact**: Free disk space ({total_size}), cleaner repository.
      **Risk**: Very low - these files are regenerated as needed.
    </template-output>

    <ask response="confirm_temp">Delete temporary files? [y/n/skip]</ask>

    <check if="user answered y or yes">
      <action>Delete all identified temp files and directories</action>
      <action>Add to cleanup_summary: "Temp files: Deleted {count} items ({total_size})"</action>
      <action>Increment total_changes by {count}</action>
    </check>

    <check if="user answered skip">
      <action>Add to cleanup_summary: "Temp files: Skipped by user"</action>
    </check>
  </check>

  <check if="no temp files">
    <action>Display: "✅ No temporary files found"</action>
    <action>Add to cleanup_summary: "Temp files: None found"</action>
  </check>

  <action>Increment categories_processed</action>
</step>

<!-- SYNC OPERATIONS -->

<step n="8" goal="Sync: Agent Templates" if="skip_agent_templates == false AND skip_sync == false">
  <action>Display: "### 6/9: Agent Templates Sync"</action>

  <action>Check for agent_customisations repository:</action>
  <action>- Expected location: ~/Documents/raegis_labs/agent_customisations</action>
  <action>- Verify templates/agents/ directory exists</action>

  <check if="agent_customisations found">
    <action>Compare local agents with templates:</action>
    <action>- Read agents from {project-root}/.bmad/expansion-packs/*/agents/</action>
    <action>- Compare with ~/Documents/raegis_labs/agent_customisations/templates/agents/</action>
    <action>- Identify agents that have updates available</action>
    <action>- Identify new agents in templates</action>
    <action>- Respect *.local.* files (never overwrite)</action>

    <check if="updates available">
      <template-output path="{project-root}/.bmad/maintenance-preview-agents.md">
        # Agent Templates Sync Preview

        ## Agents with Updates ({count})
        {list of agents with diff summary}

        ## New Agents Available ({count})
        {list of new agents not in project}

        ## Protected Local Agents ({count})
        {list of *.local.* agents that will not be touched}

        **Impact**: Get latest agent improvements and bug fixes.
        **Risk**: Low - local customizations in *.local.* files are preserved.

        **Note**: After sync, symlinks will be updated if needed.
      </template-output>

      <ask response="confirm_agents">Sync agent templates? [y/n/skip]</ask>

      <check if="user answered y or yes">
        <action>Run sync from agent_customisations:</action>
        <action>- Copy updated agents to project</action>
        <action>- Add new agents to project</action>
        <action>- Preserve *.local.* files</action>
        <action>- Update symlinks in agents/ directory</action>
        <action>Add to sync_summary: "Agent templates: Synced {count} agents"</action>
        <action>Increment total_changes by {count}</action>
      </check>

      <check if="user answered skip">
        <action>Add to sync_summary: "Agent templates: Skipped by user"</action>
      </check>
    </check>

    <check if="no updates">
      <action>Display: "✅ All agent templates are current"</action>
      <action>Add to sync_summary: "Agent templates: All current"</action>
    </check>
  </check>

  <check if="agent_customisations not found">
    <action>Display: "⚠️  agent_customisations repository not found at ~/Documents/raegis_labs/agent_customisations"</action>
    <action>Display: "   Skipping agent template sync."</action>
    <action>Add to sync_summary: "Agent templates: Repository not found"</action>
  </check>

  <action>Increment categories_processed</action>
</step>

<step n="9" goal="Sync: BMAD Framework" if="skip_bmad == false AND skip_sync == false">
  <action>Display: "### 7/9: BMAD Framework Sync"</action>

  <action>Check BMAD version:</action>
  <action>- Read current version from {project-root}/.bmad/VERSION</action>
  <action>- Check for updates (if update mechanism exists)</action>

  <action>Check BMAD core files:</action>
  <action>- Verify bmad/core/ directory integrity</action>
  <action>- Check for missing or outdated core files</action>

  <check if="updates available or files missing">
    <template-output path="{project-root}/.bmad/maintenance-preview-bmad.md">
      # BMAD Framework Sync Preview

      **Current Version**: {current_version}
      **Latest Version**: {latest_version}

      ## Changes:
      {list of core file updates or additions}

      **Impact**: Get latest BMAD features and bug fixes.
      **Risk**: Low - core updates are designed to be backward compatible.

      **Note**: Your project configuration and custom agents will not be affected.
    </template-output>

    <ask response="confirm_bmad">Update BMAD framework? [y/n/skip]</ask>

    <check if="user answered y or yes">
      <action>Update BMAD core files</action>
      <action>Update VERSION file</action>
      <action>Display: "✅ BMAD framework updated to {latest_version}"</action>
      <action>Add to sync_summary: "BMAD framework: Updated to {latest_version}"</action>
      <action>Increment total_changes</action>
    </check>

    <check if="user answered skip">
      <action>Add to sync_summary: "BMAD framework: Skipped by user"</action>
    </check>
  </check>

  <check if="no updates">
    <action>Display: "✅ BMAD framework is current"</action>
    <action>Add to sync_summary: "BMAD framework: Current version"</action>
  </check>

  <action>Increment categories_processed</action>
</step>

<step n="10" goal="Sync: Linear Settings" if="skip_linear == false AND skip_sync == false">
  <action>Display: "### 8/9: Linear Settings Sync"</action>

  <action>Check for Linear configuration templates:</action>
  <action>- Look for shared Linear settings in agent_customisations or BMAD</action>
  <action>- Check current Linear project settings (labels, states, templates)</action>

  <check if="Linear project exists">
    <action>Compare current settings with templates:</action>
    <action>- Labels: Compare label names and colors</action>
    <action>- States: Verify workflow states match standards</action>
    <action>- Templates: Check issue and project templates</action>

    <check if="differences found">
      <template-output path="{project-root}/.bmad/maintenance-preview-linear.md">
        # Linear Settings Sync Preview

        ## Missing Labels ({count})
        {list of labels to add}

        ## Different State Configuration ({count})
        {comparison of current vs standard states}

        ## Template Updates ({count})
        {list of template changes}

        **Impact**: Consistent Linear configuration across all projects.
        **Risk**: Low - existing issues and projects are not affected.
      </template-output>

      <ask response="confirm_linear">Sync Linear settings? [y/n/skip]</ask>

      <check if="user answered y or yes">
        <action>Add missing labels to Linear project</action>
        <action>Update templates</action>
        <action>Display: "✅ Linear settings synced"</action>
        <action>Add to sync_summary: "Linear settings: Synced {count} changes"</action>
        <action>Increment total_changes by {count}</action>
      </check>

      <check if="user answered skip">
        <action>Add to sync_summary: "Linear settings: Skipped by user"</action>
      </check>
    </check>

    <check if="no differences">
      <action>Display: "✅ Linear settings match standards"</action>
      <action>Add to sync_summary: "Linear settings: Already synced"</action>
    </check>
  </check>

  <check if="no Linear project">
    <action>Display: "ℹ️  No Linear project configured for this project"</action>
    <action>Add to sync_summary: "Linear settings: No project configured"</action>
  </check>

  <action>Increment categories_processed</action>
</step>

<step n="11" goal="Sync: Foundation Docs" if="skip_docs == false AND skip_sync == false">
  <action>Display: "### 9/9: Foundation Docs Sync"</action>

  <action>Check for foundation doc templates:</action>
  <action>- Look for templates in agent_customisations or BMAD</action>
  <action>- Foundation docs: idea.md, CLAUDE.md, AGENTS.md, README.md</action>

  <action>Compare current docs with templates:</action>
  <action>- Check if templates have structural updates</action>
  <action>- Identify new sections in templates</action>
  <action>- Preserve project-specific content</action>

  <check if="updates available">
    <template-output path="{project-root}/.bmad/maintenance-preview-docs.md">
      # Foundation Docs Sync Preview

      ## Structural Updates Available:

      ### idea.md
      {diff showing template structure vs current}

      ### CLAUDE.md
      {diff showing template structure vs current}

      ### AGENTS.md
      {diff showing template structure vs current}

      **Impact**: Consistent documentation structure across projects.
      **Risk**: Medium - manual merge required to preserve project content.

      **Recommendation**: Review diffs carefully. Consider updating structure manually.
      **Note**: This is informational - automatic merge of docs is NOT recommended.
    </template-output>

    <ask response="confirm_docs">Show foundation doc updates for manual review? [y/n/skip]</ask>

    <check if="user answered y or yes">
      <action>Display detailed diffs for each doc</action>
      <action>Display: "ℹ️  Review these changes and update docs manually to preserve your content."</action>
      <action>Add to sync_summary: "Foundation docs: Updates available (manual review needed)"</action>
    </check>

    <check if="user answered skip">
      <action>Add to sync_summary: "Foundation docs: Skipped by user"</action>
    </check>
  </check>

  <check if="no updates">
    <action>Display: "✅ Foundation docs structure is current"</action>
    <action>Add to sync_summary: "Foundation docs: Structure current"</action>
  </check>

  <action>Increment categories_processed</action>
</step>

<!-- FINALIZATION -->

<step n="12" goal="Update maintenance timestamp">
  <action>Write current timestamp to {project-root}/.bmad/last-maintenance</action>
  <action>Format: ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)</action>
  <action>Display: "📅 Maintenance timestamp updated: {timestamp}"</action>
</step>

<step n="13" goal="Generate summary report">
  <action>Display maintenance summary:</action>

  <template-output path="{project-root}/.bmad/last-maintenance-summary.md">
    # Maintenance Summary - {date}

    **Duration**: {start_time} - {end_time}
    **Categories Processed**: {categories_processed}/9
    **Total Changes**: {total_changes}

    ## Cleanup Operations

    {cleanup_summary - each category result}

    ## Sync Operations

    {sync_summary - each category result}

    ## Next Steps

    - Run tests to verify changes: `npm test` or `pytest`
    - Commit maintenance changes: `git add -A && git commit -m "chore: project maintenance"`
    - Next maintenance recommended: {30 days from now}

    ## Preview Files

    Detailed previews saved to:
    {list of .bmad/maintenance-preview-*.md files}

    These files can be deleted after review.
  </template-output>

  <action>Display summary to user:</action>
  <action>
    "✅ **Maintenance Complete**

    - Processed {categories_processed}/9 categories
    - Made {total_changes} changes
    - Summary saved to .bmad/last-maintenance-summary.md

    **Next Steps:**
    1. Review changes: git diff
    2. Run tests to verify
    3. Commit changes
    4. Next maintenance recommended in 30 days
    "
  </action>
</step>

## Notes for Implementation

### Preview Mechanism
- Each category generates a preview file in `.bmad/maintenance-preview-{category}.md`
- Previews show exactly what will change with diffs or summaries
- User confirms category by category
- Preview files can be reviewed before confirming

### Confirmation Options
- `y` or `yes` - Apply changes for this category
- `n` or `no` - Skip this category
- `skip` - Skip and don't ask again (for --yes mode)

### Safety Features
- **Dry Run Mode** (`--dry-run`): Shows all previews without making changes
- **Git Safety**: All cleanup operations preserve git history
- **Protected Files**: *.local.* files are never overwritten
- **Selective Skip**: Can skip individual categories with flags
- **Timestamp Tracking**: Records last maintenance date

### Integration Points
- Uses same logic as `/init-sync-bmad` for agent template syncing
- Respects .gitignore for temp file detection
- Uses git commands for branch cleanup
- Uses package manager commands for dependency updates

### Error Handling
- If a category fails, continue with remaining categories
- Record failures in summary
- Don't update timestamp if critical errors occurred
- Preserve preview files for debugging
