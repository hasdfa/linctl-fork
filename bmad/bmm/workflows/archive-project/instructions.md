# Archive Project Workflow Instructions

## Overview
This workflow performs comprehensive project archival - checking for uncommitted work, exporting data, cleaning up integrations, and compressing the project to an archive location. This is designed for retiring projects that are no longer active.

**⚠️ IMPORTANT**: Archival is a significant operation. The project will be removed from active locations and compressed. While restoration is possible, it requires manual steps.

## Execution Flow

<step n="1" goal="Explain archival process and get confirmation to proceed">
  <action>Display archival overview:</action>
  <action>
    "# Project Archival

    This workflow will:

    **1. Pre-Archive Checks**
    - Check for uncommitted changes
    - Check for unpushed commits
    - Check for open Linear issues
    - Check for recent activity (commits in last 30 days)
    - Check for running processes

    **2. Data Export**
    - Export Linear issues to JSON
    - Export git metadata and history
    - Export project statistics
    - Copy documentation files
    - Copy configuration templates (no secrets)

    **3. Cleanup Operations**
    - Remove from WezTerm launcher
    - Close Linear project (set to 'Canceled')
    - Remove from IDE recent projects
    - Clear build artifacts and caches

    **4. Archive Creation**
    - Create archive metadata file
    - Compress project to .tar.gz
    - Move to ~/Documents/archived/{date}-{project-name}

    **⚠️ This is a one-shot operation**: Either archive completely or cancel.
    No partial archival supported.
    "
  </action>

  <ask response="proceed_check">Do you want to proceed with pre-archive checks? [y/n]</ask>

  <check if="user answered n or no">
    <action>Display: "Archival canceled. Project remains unchanged."</action>
    <action>Exit workflow</action>
  </check>
</step>

<step n="2" goal="Pre-Check 1: Git Status" if="skip_checks == false">
  <action>Display: "## Running Pre-Archive Checks (1/5)"</action>

  <action>Check git status:</action>
  <action>- Run: git status --porcelain</action>
  <action>- Identify uncommitted changes (modified, untracked, staged)</action>
  <action>- Count uncommitted files</action>

  <check if="uncommitted changes found">
    <action>Set: has_uncommitted = true</action>
    <action>Display warning:</action>
    <action>
      "⚠️  **Uncommitted Changes Detected** ({count} files)

      {list of uncommitted files}

      Recommendation: Commit or stash these changes before archiving.
      "
    </action>
  </check>

  <check if="no uncommitted changes">
    <action>Set: has_uncommitted = false</action>
    <action>Display: "✅ Git status clean"</action>
  </check>
</step>

<step n="3" goal="Pre-Check 2: Unpushed Commits" if="skip_checks == false">
  <action>Display: "## Running Pre-Archive Checks (2/5)"</action>

  <action>Check for unpushed commits:</action>
  <action>- Run: git log origin/main..HEAD --oneline (or appropriate remote branch)</action>
  <action>- Count unpushed commits</action>

  <check if="unpushed commits found">
    <action>Set: has_unpushed = true</action>
    <action>Display warning:</action>
    <action>
      "⚠️  **Unpushed Commits Detected** ({count} commits)

      {list of commit messages}

      Recommendation: Push to remote before archiving to preserve work.
      "
    </action>
  </check>

  <check if="no unpushed commits">
    <action>Set: has_unpushed = false</action>
    <action>Display: "✅ All commits pushed to remote"</action>
  </check>

  <check if="no remote configured">
    <action>Set: has_no_remote = true</action>
    <action>Display warning:</action>
    <action>
      "⚠️  **No Remote Repository Configured**

      Commits exist only locally. Consider pushing to GitHub before archiving.
      "
    </action>
  </check>
</step>

<step n="4" goal="Pre-Check 3: Open Linear Issues" if="skip_checks == false">
  <action>Display: "## Running Pre-Archive Checks (3/5)"</action>

  <action>Check for Linear project:</action>

  <check if="Linear project exists">
    <action>Get Linear issues:</action>
    <action>- Query issues for this project (all states)</action>
    <action>- Count open issues (Todo, In Progress, In Review)</action>
    <action>- Count done/canceled issues</action>

    <check if="open issues found">
      <action>Set: has_open_issues = true</action>
      <action>Set: open_issue_count = {count}</action>
      <action>Display warning:</action>
      <action>
        "⚠️  **Open Linear Issues** ({count} issues)

        {list of open issues with state and title}

        These issues will be exported to JSON but the Linear project will be set to 'Canceled'.
        Consider completing or closing issues before archiving.
        "
      </action>
    </check>

    <check if="no open issues">
      <action>Set: has_open_issues = false</action>
      <action>Display: "✅ No open Linear issues"</action>
    </check>
  </check>

  <check if="no Linear project">
    <action>Set: has_no_linear = true</action>
    <action>Display: "ℹ️  No Linear project configured"</action>
  </check>
</step>

<step n="5" goal="Pre-Check 4: Recent Activity" if="skip_checks == false">
  <action>Display: "## Running Pre-Archive Checks (4/5)"</action>

  <action>Check for recent activity:</action>
  <action>- Get last commit date: git log -1 --format=%cd --date=iso</action>
  <action>- Calculate days since last commit</action>
  <action>- Check file modification times (most recent .py, .js, .ts, .md file)</action>

  <check if="activity in last 30 days">
    <action>Set: has_recent_activity = true</action>
    <action>Display warning:</action>
    <action>
      "⚠️  **Recent Activity Detected**

      Last commit: {days} days ago ({date})
      Last file modification: {days} days ago ({file})

      This project appears to be active. Are you sure you want to archive it?
      "
    </action>
  </check>

  <check if="no activity in last 30 days">
    <action>Set: has_recent_activity = false</action>
    <action>Display: "✅ No recent activity (inactive for {days} days)"</action>
  </check>
</step>

<step n="6" goal="Pre-Check 5: Running Processes" if="skip_checks == false">
  <action>Display: "## Running Pre-Archive Checks (5/5)"</action>

  <action>Check for running processes related to this project:</action>
  <action>- Check for processes in current directory: lsof +D {project-root}</action>
  <action>- Check for common dev servers: npm, node, python, uvicorn, next, vite</action>
  <action>- Check for database processes: postgres, mysql, mongo, redis</action>

  <check if="running processes found">
    <action>Set: has_running_processes = true</action>
    <action>Display warning:</action>
    <action>
      "⚠️  **Running Processes Detected** ({count} processes)

      {list of processes with PID and command}

      Recommendation: Stop these processes before archiving.
      "
    </action>
  </check>

  <check if="no running processes">
    <action>Set: has_running_processes = false</action>
    <action>Display: "✅ No running processes"</action>
  </check>
</step>

<step n="7" goal="Compile pre-check results and show summary">
  <action>Display: "## Pre-Check Summary"</action>

  <action>Count warnings:</action>
  <action>- warning_count = has_uncommitted + has_unpushed + has_no_remote + has_open_issues + has_recent_activity + has_running_processes</action>

  <check if="warning_count == 0">
    <action>Display:</action>
    <action>
      "✅ **All Pre-Checks Passed**

      This project is ready for archival:
      - Git status clean
      - All commits pushed
      - No open Linear issues
      - No recent activity
      - No running processes

      Safe to proceed with archival.
      "
    </action>
  </check>

  <check if="warning_count > 0">
    <action>Display:</action>
    <action>
      "⚠️  **{warning_count} Warning(s) Detected**

      {list each warning category with status}

      You can proceed with archival, but these warnings should be addressed:
      - Uncommitted changes will be included in archive (not lost)
      - Unpushed commits will only exist in archive (not on remote)
      - Open Linear issues will be exported but project will be canceled
      - Recent activity suggests project may still be in use
      - Running processes should be stopped before archiving

      Options:
      1. Address warnings and run archival again
      2. Proceed anyway (with --force if needed)
      3. Cancel archival
      "
    </action>
  </check>
</step>

<step n="8" goal="Get archive reason and final confirmation">
  <ask response="archive_reason">Why are you archiving this project? (e.g., 'Completed', 'Abandoned', 'Superseded by X', etc.)</ask>

  <action>Set: archive_reason = {user response}</action>

  <action>Calculate archive details:</action>
  <action>- Project name: {extract from current directory}</action>
  <action>- Archive date: {current date YYYY-MM-DD}</action>
  <action>- Archive name: {archive_date}-{project_name}</action>
  <action>- Archive path: ~/Documents/archived/{archive_name}</action>
  <action>- Project size: {calculate du -sh}</action>
  <action>- File count: {calculate find . -type f | wc -l}</action>

  <template-output path="{project-root}/.bmad/archive-preview.md">
    # Archive Preview - {project_name}

    ## Archive Details
    - **Project**: {project_name}
    - **Current Path**: {project-root}
    - **Archive Path**: ~/Documents/archived/{archive_name}
    - **Archive Reason**: {archive_reason}
    - **Archive Date**: {archive_date}
    - **Project Size**: {size}
    - **File Count**: {file_count}

    ## Pre-Check Warnings
    {list of warnings if any, or "None" if all passed}

    ## What Will Be Done

    ### 1. Data Export
    - Linear issues → `exported-data/linear-issues.json`
    - Git metadata → `exported-data/git-metadata.json`
    - Project stats → `exported-data/project-stats.json`
    - Documentation → `exported-data/docs/` (all .md files)
    - Configuration → `exported-data/config/` (templates, no secrets)

    ### 2. Cleanup Operations
    - Remove from WezTerm launcher config
    - Set Linear project state to "Canceled"
    - Clear from IDE recent projects
    - Delete node_modules, __pycache__, dist/, build/

    ### 3. Archive Creation
    - Create `ARCHIVE_METADATA.json` with project details
    - Create `RESTORATION.md` with unarchive instructions
    - Compress to `{archive_name}.tar.gz`
    - Move archive to `~/Documents/archived/{archive_name}/`

    ### 4. Post-Archive
    - Original project directory removed from current location
    - Archive available at `~/Documents/archived/{archive_name}/`

    ## Restoration
    To restore this project later:
    1. Extract archive: `tar -xzf {archive_name}.tar.gz`
    2. Move to desired location
    3. Run `git remote -v` to verify remote URL
    4. Run `npm install` or `pip install -r requirements.txt`
    5. Add back to WezTerm launcher if needed
    6. Reopen Linear project if needed

    ## ⚠️ IMPORTANT
    - This is a **one-shot operation** - full archive or cancel
    - Original project will be removed from current location
    - Archive will be compressed and moved to archive directory
    - Restoration requires manual steps
  </template-output>

  <action>Display: "📋 Archive preview saved to .bmad/archive-preview.md"</action>

  <check if="warning_count > 0 AND force == false">
    <ask response="proceed_with_warnings">Proceed with archival despite {warning_count} warning(s)? [y/n]</ask>

    <check if="user answered n or no">
      <action>Display: "Archival canceled. Project remains unchanged."</action>
      <action>Display: "Address the warnings and run archive-project again when ready."</action>
      <action>Exit workflow</action>
    </check>
  </check>

  <ask response="final_confirmation">
    **FINAL CONFIRMATION**

    Archive "{project_name}" to ~/Documents/archived/{archive_name}?

    This will:
    - Export all data
    - Clean up integrations
    - Compress and move project
    - Remove from current location

    Type the project name "{project_name}" to confirm, or "cancel" to abort:
  </ask>

  <check if="user typed project name correctly">
    <action>Display: "✅ Confirmation received. Starting archival process..."</action>
  </check>

  <check if="user typed cancel">
    <action>Display: "Archival canceled. Project remains unchanged."</action>
    <action>Exit workflow</action>
  </check>

  <check if="user typed anything else">
    <action>Display: "❌ Incorrect project name. Archival canceled for safety."</action>
    <action>Exit workflow</action>
  </check>
</step>

<step n="9" goal="Create archive directory structure">
  <action>Display: "## Creating Archive Structure"</action>

  <action>Create archive directory:</action>
  <action>- mkdir -p ~/Documents/archived/{archive_name}</action>
  <action>- mkdir -p ~/Documents/archived/{archive_name}/exported-data</action>
  <action>- mkdir -p ~/Documents/archived/{archive_name}/exported-data/docs</action>
  <action>- mkdir -p ~/Documents/archived/{archive_name}/exported-data/config</action>

  <action>Display: "✅ Archive directory created at ~/Documents/archived/{archive_name}"</action>
</step>

<step n="10" goal="Export Linear data" if="skip_export == false AND has_no_linear == false">
  <action>Display: "## Exporting Linear Data"</action>

  <action>Export Linear issues (via linctl):</action>
  <action>- Query all issues for this project (all states)</action>
  <action>- Include: title, description, state, assignee, labels, comments, created/updated dates</action>
  <action>- Command example: linctl issue list --team $LINEAR_TEAM --include-completed --newer-than all_time --json > linear-issues.json</action>
  <action>- Format as JSON</action>
  <action>- Save to ~/Documents/archived/{archive_name}/exported-data/linear-issues.json</action>

  <action>Export Linear project metadata (via linctl):</action>
  <action>- Project name, description, state, members, milestones</action>
  <action>- Command example: linctl project get {linear_project_id} --json > linear-project.json</action>
  <action>- Save to ~/Documents/archived/{archive_name}/exported-data/linear-project.json</action>

  <action>Display: "✅ Linear data exported ({issue_count} issues)"</action>
</step>

<step n="11" goal="Export Git metadata" if="skip_export == false">
  <action>Display: "## Exporting Git Metadata"</action>

  <action>Export git information:</action>
  <action>- Remote URLs: git remote -v</action>
  <action>- Branch list: git branch -a</action>
  <action>- Commit history summary: git log --oneline --all --graph (last 100 commits)</action>
  <action>- Contributors: git shortlog -sn</action>
  <action>- Tags: git tag</action>

  <action>Format as JSON and save to ~/Documents/archived/{archive_name}/exported-data/git-metadata.json</action>

  <action>Display: "✅ Git metadata exported"</action>
</step>

<step n="12" goal="Export project statistics" if="skip_export == false">
  <action>Display: "## Exporting Project Statistics"</action>

  <action>Calculate project statistics:</action>
  <action>- Total files: find . -type f | wc -l (excluding node_modules, .git)</action>
  <action>- Total size: du -sh</action>
  <action>- Lines of code by language: cloc . (if available) or manual count</action>
  <action>- Dependencies: package.json, requirements.txt, etc.</action>
  <action>- Last modified: stat (most recent file)</action>

  <action>Format as JSON and save to ~/Documents/archived/{archive_name}/exported-data/project-stats.json</action>

  <action>Display: "✅ Project statistics exported"</action>
</step>

<step n="13" goal="Export documentation" if="skip_export == false">
  <action>Display: "## Exporting Documentation"</action>

  <action>Copy all documentation files:</action>
  <action>- Find all .md files: find . -name "*.md" (excluding node_modules)</action>
  <action>- Copy to ~/Documents/archived/{archive_name}/exported-data/docs/ (preserve directory structure)</action>

  <action>Display: "✅ Documentation exported ({doc_count} files)"</action>
</step>

<step n="14" goal="Export configuration templates" if="skip_export == false">
  <action>Display: "## Exporting Configuration"</action>

  <action>Copy configuration files (no secrets):</action>
  <action>- .env.example (NOT .env)</action>
  <action>- config.example.yaml</action>
  <action>- .editorconfig, .prettierrc, .eslintrc.js</action>
  <action>- package.json, requirements.txt, pyproject.toml</action>
  <action>- tsconfig.json, vite.config.js, etc.</action>

  <action>Copy to ~/Documents/archived/{archive_name}/exported-data/config/</action>

  <action>Display: "✅ Configuration exported (secrets excluded)"</action>
</step>

<step n="15" goal="Remove from WezTerm launcher" if="skip_cleanup == false">
  <action>Display: "## Cleaning Up Integrations (1/4)"</action>

  <action>Find WezTerm launcher config:</action>
  <action>- Expected location: ~/.config/wezterm/wezterm.lua or similar</action>

  <check if="WezTerm config found">
    <action>Search for project path in config file</action>

    <check if="project found in config">
      <action>Remove project entry from launcher</action>
      <action>Display: "✅ Removed from WezTerm launcher"</action>
    </check>

    <check if="project not in config">
      <action>Display: "ℹ️  Project not found in WezTerm launcher"</action>
    </check>
  </check>

  <check if="WezTerm config not found">
    <action>Display: "ℹ️  WezTerm config not found, skipping"</action>
  </check>
</step>

<step n="16" goal="Close Linear project" if="skip_cleanup == false AND has_no_linear == false">
  <action>Display: "## Cleaning Up Integrations (2/4)"</action>

  <action>Update Linear project state (via linctl):</action>
  <action>- Set project state to "Canceled"</action>
  <action>- Add archive note to project description</action>
  <action>- Note: "Archived on {date}. Reason: {archive_reason}. Archive location: ~/Documents/archived/{archive_name}"</action>

  <action>Display: "✅ Linear project set to 'Canceled'"</action>
</step>

<step n="17" goal="Remove from IDE recent projects" if="skip_cleanup == false">
  <action>Display: "## Cleaning Up Integrations (3/4)"</action>

  <action>Check for IDE recent project lists:</action>
  <action>- VSCode: ~/.config/Code/storage.json</action>
  <action>- Claude Code: (if recent projects tracked)</action>
  <action>- Codex: (if recent projects tracked)</action>

  <action>Remove project from recent lists (best effort)</action>

  <action>Display: "✅ Removed from IDE recent projects (where possible)"</action>
</step>

<step n="18" goal="Clear build artifacts and caches" if="skip_cleanup == false">
  <action>Display: "## Cleaning Up Integrations (4/4)"</action>

  <action>Delete build artifacts and caches:</action>
  <action>- node_modules/ (Node.js)</action>
  <action>- __pycache__/, *.pyc (Python)</action>
  <action>- dist/, build/ (build outputs)</action>
  <action>- .next/, .nuxt/, .vite/ (framework caches)</action>
  <action>- .pytest_cache/, .mypy_cache/, .ruff_cache/</action>
  <action>- .DS_Store, Thumbs.db</action>

  <action>Calculate space freed</action>

  <action>Display: "✅ Cleared caches and build artifacts ({size} freed)"</action>
</step>

<step n="19" goal="Create archive metadata">
  <action>Display: "## Creating Archive Metadata"</action>

  <action>Create ARCHIVE_METADATA.json:</action>
  <action>
    {
      "archive_date": "{ISO 8601 timestamp}",
      "project_name": "{project_name}",
      "project_path_original": "{project-root}",
      "git_remote_url": "{remote URL or null}",
      "linear_project_id": "{Linear project ID or null}",
      "last_commit_date": "{last commit date}",
      "last_modified_date": "{last file modification}",
      "file_count": {file_count},
      "total_size": "{size}",
      "archive_reason": "{archive_reason}",
      "archived_by": "{user_name}",
      "warnings": [{list of warning categories that were present}],
      "restoration_notes": "See RESTORATION.md for unarchive instructions"
    }
  </action>

  <action>Save to ~/Documents/archived/{archive_name}/ARCHIVE_METADATA.json</action>

  <action>Display: "✅ Archive metadata created"</action>
</step>

<step n="20" goal="Create restoration instructions">
  <action>Display: "## Creating Restoration Instructions"</action>

  <action>Create RESTORATION.md:</action>
  <action>
    # Restoration Instructions - {project_name}

    **Archived**: {archive_date}
    **Reason**: {archive_reason}

    ## Quick Restore

    \`\`\`bash
    # 1. Extract archive
    cd ~/Documents/archived/{archive_name}
    tar -xzf {project_name}.tar.gz

    # 2. Move to desired location
    mv {project_name} ~/path/to/restore/location/

    # 3. Install dependencies
    cd ~/path/to/restore/location/{project_name}
    {dependency install commands based on project type}

    # 4. Verify git remote
    git remote -v

    # 5. (Optional) Add back to WezTerm launcher
    # (Manual step, see exported-data/wezterm-entry.txt if available)

    # 6. (Optional) Reopen Linear project
    # (Manual step, project was set to 'Canceled')
    \`\`\`

    ## Exported Data

    All project data has been exported to \`exported-data/\`:

    - \`linear-issues.json\` - All Linear issues and comments
    - \`git-metadata.json\` - Git history, branches, tags, contributors
    - \`project-stats.json\` - File counts, LOC, dependencies
    - \`docs/\` - All documentation files (.md)
    - \`config/\` - Configuration templates (no secrets)

    ## Notes

    {any project-specific restoration notes}

    ## Archive Metadata

    See \`ARCHIVE_METADATA.json\` for complete archive details.
  </action>

  <action>Save to ~/Documents/archived/{archive_name}/RESTORATION.md</action>

  <action>Display: "✅ Restoration instructions created"</action>
</step>

<step n="21" goal="Compress project">
  <action>Display: "## Compressing Project"</action>

  <action>Create compressed archive:</action>
  <action>- cd to parent directory of project</action>
  <action>- Run: tar -czf ~/Documents/archived/{archive_name}/{project_name}.tar.gz {project_name}/</action>
  <action>- Exclude: .git/ (optional, can include for full history), node_modules/, __pycache__/, dist/, build/</action>

  <action>Calculate compressed size</action>

  <action>Display: "✅ Project compressed to {compressed_size} (from {original_size})"</action>
</step>

<step n="22" goal="Move archive and remove original">
  <action>Display: "## Finalizing Archive"</action>

  <action>Verify archive integrity:</action>
  <action>- Test extraction: tar -tzf ~/Documents/archived/{archive_name}/{project_name}.tar.gz > /dev/null</action>

  <check if="archive valid">
    <action>Display: "✅ Archive integrity verified"</action>

    <action>Remove original project directory:</action>
    <action>- Display: "🗑️  Removing original project from {project-root}"</action>
    <action>- Run: rm -rf {project-root}</action>

    <action>Display: "✅ Original project removed"</action>
  </check>

  <check if="archive invalid">
    <action>Display error: "❌ Archive integrity check failed! Original project NOT removed."</action>
    <action>Display: "   Please verify archive manually before removing original."</action>
    <action>Exit workflow with error</action>
  </check>
</step>

<step n="23" goal="Generate final summary and provide next steps">
  <action>Display final summary:</action>
  <action>
    "✅ **Archive Complete**

    **Project**: {project_name}
    **Archive Location**: ~/Documents/archived/{archive_name}/
    **Archive Size**: {compressed_size}
    **Archive Reason**: {archive_reason}

    ## Archive Contents

    - ✅ Project compressed to {project_name}.tar.gz
    - ✅ Linear data exported ({issue_count} issues)
    - ✅ Git metadata exported
    - ✅ Project statistics exported
    - ✅ Documentation exported ({doc_count} files)
    - ✅ Configuration exported (no secrets)
    - ✅ Archive metadata created (ARCHIVE_METADATA.json)
    - ✅ Restoration instructions created (RESTORATION.md)

    ## Cleanup Performed

    - ✅ Removed from WezTerm launcher
    - ✅ Linear project set to 'Canceled'
    - ✅ Removed from IDE recent projects
    - ✅ Caches cleared ({cache_size} freed)

    ## Original Project

    - 🗑️  Removed from {project-root}

    ## Restoration

    To restore this project, see:
    ~/Documents/archived/{archive_name}/RESTORATION.md

    Or quick restore:
    \`\`\`bash
    cd ~/Documents/archived/{archive_name}
    tar -xzf {project_name}.tar.gz
    mv {project_name} ~/desired/location/
    cd ~/desired/location/{project_name}
    # Follow RESTORATION.md for next steps
    \`\`\`

    ## Archive Management

    Archives are organized by date in ~/Documents/archived/
    Consider periodic review to remove very old archives if needed.
    "
  </action>

  <action>Log archival to system log (if available):</action>
  <action>- Record: Project "{project_name}" archived on {date} to ~/Documents/archived/{archive_name}/</action>
</step>

## Notes for Implementation

### Safety Features
- **Type-to-confirm**: User must type exact project name to confirm archival
- **Archive integrity check**: tar test extraction before removing original
- **Full data export**: All data exported before cleanup
- **Restoration instructions**: Clear instructions included in archive
- **Metadata preservation**: Complete project metadata saved
- **No partial operations**: All-or-nothing approach (no partial archival)

### Pre-Check Philosophy
- Run all checks but allow proceeding with warnings
- User can address warnings and re-run if desired
- Force flag available for automated scenarios
- Running processes check prevents file locks

### Data Export Coverage
- Linear: All issues, comments, project metadata
- Git: History, branches, tags, contributors, remotes
- Statistics: File counts, LOC, size, dependencies
- Docs: All .md files with directory structure preserved
- Config: Templates only, no secrets (.env excluded)

### Cleanup Strategy
- Remove from active integrations (WezTerm, IDE, Linear)
- Clear caches and artifacts to reduce archive size
- Set Linear project to "Canceled" (not deleted, recoverable)
- Local cleanup only (no remote operations)

### Archive Organization
- Date-prefixed naming: YYYY-MM-DD-project-name
- Chronological sorting in archive directory
- Self-contained archives (metadata + data + project)
- Compression: tar.gz for universal compatibility

### Restoration Approach
- Clear, step-by-step instructions
- Dependency installation commands included
- Manual steps for integrations (WezTerm, Linear)
- Exported data available for reference

### Error Handling
- Archive integrity check required before removal
- Failed archival leaves original untouched
- Partial exports still useful even if some fail
- Clear error messages with recovery steps

### Dry Run Support
- Preview all operations without executing
- Show warnings and data export details
- Generate preview report
- Safe for testing archival process
