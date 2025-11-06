# Quick Setup Workflow Instructions

<critical>The workflow execution engine is governed by: {project-root}/bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/bmad/bmm/workflows/quick-setup/workflow.yaml</critical>
<critical>Communicate in {communication_language} throughout the setup process</critical>

<workflow>

<step n="0" goal="Parse flags and determine execution mode">
  <action>Check for flags in user input: --dry-run, --skip-git, --skip-bmad, --skip-github, --skip-docs</action>
  <action>Set dry_run_mode = true if --dry-run flag present</action>
  <action>Set skip flags for each component</action>

  <check if="dry_run_mode == true">
    <action>Display: "DRY RUN MODE - No changes will be made"</action>
  </check>

  <template-output>execution_mode</template-output>
</step>

<step n="1" goal="Determine project type and BMAD module">
  <action>Ask user to select project type from: software, data-science, writing, business, design</action>

  <ask response="project_type">What type of project is this?
    1. Software Development
    2. Data Science / Analytics
    3. Writing / Documentation
    4. Business / Strategy
    5. Design / Creative
  </ask>

  <action>Map selection to BMAD module:</action>
  <action>- software → bmm (Business Management Module)</action>
  <action>- data-science → bmm</action>
  <action>- writing → bmm</action>
  <action>- business → bmm</action>
  <action>- design → bmm</action>

  <note>Currently all map to bmm. Future: specialized modules per type.</note>

  <template-output>project_type, bmad_module</template-output>
</step>

<step n="2" goal="Plan execution steps" optional="false">
  <action>Generate execution plan based on flags and project type</action>

  <action>Build step list:</action>
  <action>1. Directory scaffolding (unless --skip-dirs)</action>
  <action>2. Git initialization (unless --skip-git or git exists)</action>
  <action>3. Foundation documents (unless --skip-docs)</action>
  <action>4. BMAD v6 installation (unless --skip-bmad)</action>
  <action>5. GitHub repository creation (unless --skip-github)</action>
  <action>6. Linear project setup (unless --skip-linear)</action>
  <action>7. WezTerm launcher integration (unless --skip-wezterm)</action>
  <action>8. Template sync from remote (unless --skip-sync)</action>

  <action>Display execution plan with estimated time</action>

  <check if="dry_run_mode == false">
    <ask>Ready to proceed with setup? [y/n]</ask>
    <check if="user answered no">
      <action>Exit workflow with message: "Setup cancelled. Re-run when ready."</action>
    </check>
  </check>

  <template-output>execution_plan</template-output>
</step>

<step n="3" goal="Create directory structure" if="skip_dirs == false">
  <action>Create standard directories:</action>
  <action>- src/ (source code)</action>
  <action>- tests/ (test files)</action>
  <action>- docs/ (documentation)</action>
  <action>- .bmad/ (BMAD framework files)</action>
  <action>- .claude/commands/ (Claude slash commands)</action>

  <check if="project_type == 'software'">
    <action>Add: lib/, scripts/, config/</action>
  </check>

  <check if="project_type == 'data-science'">
    <action>Add: data/, notebooks/, models/</action>
  </check>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would create directories: [list]"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Execute: mkdir -p [directories]</action>
    <action>Display: "✓ Directories created"</action>
  </check>

  <template-output>directories_created</template-output>
</step>

<step n="4" goal="Initialize Git repository" if="skip_git == false">
  <action>Check if .git directory already exists</action>

  <check if=".git exists">
    <action>Display: "✓ Git already initialized, skipping"</action>
    <goto step="5"/>
  </check>

  <action>Get current directory name for repo naming</action>
  <ask response="repo_name">Repository name (default: current directory)?</ask>

  <action>Detect project type for .gitignore template:</action>
  <action>- Check for package.json → Node.js</action>
  <action>- Check for requirements.txt → Python</action>
  <action>- Check for go.mod → Go</action>
  <action>- Default → Generic</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would initialize git and create .gitignore for [detected_type]"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Execute: git init</action>
    <action>Create .gitignore based on detected type</action>
    <action>Add common ignores: .bmad/, .claude/, .env, node_modules/, __pycache__/</action>
    <action>Display: "✓ Git initialized with .gitignore"</action>
  </check>

  <template-output>git_initialized, gitignore_type</template-output>
</step>

<step n="5" goal="Create foundation documents" if="skip_docs == false">
  <action>Generate foundation documents with project context</action>

  <ask response="project_description">Briefly describe this project (1-2 sentences):</ask>

  <action>Create idea.md with structure:</action>
  <action>- Project name and description</action>
  <action>- Core purpose and goals</action>
  <action>- Target users</action>
  <action>- Success criteria</action>

  <action>Create CLAUDE.md with structure:</action>
  <action>- Project overview</action>
  <action>- Development guidelines</action>
  <action>- File locations and structure</action>
  <action>- Context for Claude Code</action>

  <action>Create AGENTS.md with structure:</action>
  <action>- Agent configuration for Codex</action>
  <action>- Available agents and workflows</action>
  <action>- Integration points</action>

  <action>Create README.md with structure:</action>
  <action>- Project title and description</action>
  <action>- Setup instructions</action>
  <action>- Usage</action>
  <action>- Development</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would create: idea.md, CLAUDE.md, AGENTS.md, README.md"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Write all foundation documents</action>
    <action>Display: "✓ Foundation documents created"</action>
  </check>

  <template-output>foundation_docs_created</template-output>
</step>

<step n="6" goal="Install BMAD v6" if="skip_bmad == false">
  <action>Check if BMAD already installed (.bmad/ exists and populated)</action>

  <check if="BMAD already installed">
    <action>Display: "✓ BMAD already installed, skipping"</action>
    <goto step="7"/>
  </check>

  <action>Install BMAD {bmad_module} module</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would install BMAD v6:"</action>
    <action>Display: "   Module: {bmad_module}"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Execute BMAD installation for {bmad_module}</action>
    <action>Create .bmad/ structure</action>
    <action>Initialize module configuration</action>
    <action>Set up agent symlinks</action>
    <action>Display: "✓ BMAD v6 installed ({bmad_module} module)"</action>
  </check>

  <template-output>bmad_installed, bmad_module</template-output>
</step>

<step n="7" goal="Create GitHub repository" if="skip_github == false">
  <action>Check if GitHub remote already configured</action>

  <check if="remote origin exists">
    <action>Display: "✓ GitHub remote already configured, skipping"</action>
    <goto step="8"/>
  </check>

  <action>Verify gh CLI is installed</action>

  <check if="gh not installed">
    <action>Display: "⚠️  GitHub CLI (gh) not found. Skipping GitHub repo creation."</action>
    <action>Display: "Install with: brew install gh"</action>
    <goto step="8"/>
  </check>

  <action>Use repo_name from step 4 or ask again if not set</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would create private GitHub repo: {repo_name}"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Execute: gh repo create {repo_name} --private --source=. --remote=origin</action>
    <action>Execute: git add .</action>
    <action>Execute: git commit -m "Initial commit"</action>
    <action>Execute: git push -u origin main</action>
    <action>Get repository URL from gh</action>
    <action>Display: "✓ GitHub repository created and pushed"</action>
    <action>Display: "   URL: {repo_url}"</action>
  </check>

  <template-output>github_created, repo_url</template-output>
</step>

<step n="7.5" goal="Configure Linear integration with environment variables" if="skip_linear == false">
  <action>Verify required tools installed</action>

  <check if="linctl not installed">
    <action>Display: "⚠️  linctl not found. Skipping Linear integration."</action>
    <action>Display: "Install with: brew install linctl (or see https://github.com/dorkitude/linctl)"</action>
    <goto step="7.6"/>
  </check>

  <ask response="configure_linear">Configure Linear integration with environment variables? [y/n]</ask>

  <check if="user answered no">
    <action>Display: "Skipping Linear integration"</action>
    <goto step="7.6"/>
  </check>

  <action>Infer Linear label and WezTerm section from parent directory:</action>
  <action>- Get parent directory name of current project</action>
  <action>- If parent contains "peps" (case-insensitive) → Label: "Peps Ventures", Section: "PEPS VENTURES"</action>
  <action>- If parent contains "raegis" (case-insensitive) → Label: "Raegis Labs", Section: "RAEGIS LABS"</action>
  <action>- If parent contains "xognosis" (case-insensitive) → Label: "XOGNOSIS", Section: "XOGNOSIS"</action>
  <action>- Otherwise → Label: "General", Section: "PROJECTS"</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would configure Linear integration:"</action>
    <action>Display: "   1. Verify linctl authentication (linctl auth status)"</action>
    <action>Display: "   2. Set LINEAR_TEAM=RAE (Raegis Labs)"</action>
    <action>Display: "   3. Match or create project '{repo_name}' via scripts/linear-create-project.sh"</action>
    <action>Display: "   4. Write to .envrc (LINEAR_TEAM, LINEAR_PROJECT) and run direnv allow"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Step 1: Verify linctl authentication</action>
    <action>Execute: linctl auth status</action>

    <check if="not authenticated">
      <action>Display: "⚠️  linctl not authenticated. Run: linctl auth"</action>
      <action>Display: "Skipping Linear integration"</action>
      <goto step="7.6"/>
    </check>

    <action>Step 2: Set Linear team (hardcoded Raegis Labs)</action>
    <action>Set: linear_team_id = "b8ff8916-3e03-435d-809f-9d45ef4199c8"</action>
    <action>Set: linear_team_key = "RAE"</action>

    <action>Step 3: Find or create Linear project</action>
    <action>Execute: linctl project list --newer-than all_time --json | jq -r '.[] | select(.name == "{repo_name}") | .id'</action>

    <check if="project not found">
      <action>Display: "Creating new Linear project: {repo_name}"</action>
      <action>Execute: scripts/linear-create-project.sh "{repo_name}" RAE</action>
      <action>Store result as: linear_project_id</action>
    </check>

    <check if="project found">
      <action>Display: "Using existing Linear project: {repo_name}"</action>
      <action>Store found ID as: linear_project_id</action>
    </check>

    <action>Step 4: Create/update .envrc</action>
    <action>Append or update entries:</action>
    <action>export LINEAR_TEAM=RAE</action>
    <action>export LINEAR_PROJECT={linear_project_id}</action>
    <action>Execute if direnv installed: direnv allow</action>

    <action>Display: "✓ Linear integration configured (linctl)"</action>
    <action>Display: "   Team: Raegis Labs ({linear_team_key})"</action>
    <action>Display: "   Project: {repo_name}"</action>
  </check>

  <template-output>linear_configured, linear_project_id, linear_project_url, env_created, inferred_label, inferred_section</template-output>
</step>

<step n="7.6" goal="Add to WezTerm launcher" if="skip_wezterm == false">
  <action>Determine WezTerm section from parent directory (use from step 7.5 if available, else infer now):</action>
  <action>- If parent contains "peps" (case-insensitive) → Section: "PEPS VENTURES"</action>
  <action>- If parent contains "raegis" (case-insensitive) → Section: "RAEGIS LABS"</action>
  <action>- If parent contains "xognosis" (case-insensitive) → Section: "XOGNOSIS"</action>
  <action>- Otherwise → Section: "PROJECTS"</action>

  <action>Choose icon based on project type:</action>
  <action>- software → 💻 or 🖥️</action>
  <action>- data-science → 📊 or 🔬</action>
  <action>- writing → 📝 or 📚</action>
  <action>- business → 💼 or 📋</action>
  <action>- design → 🎨 or ✨</action>
  <action>- default → 📁</action>

  <action>Get absolute path to current project directory</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would add to WezTerm launcher:"</action>
    <action>Display: "   Section: {inferred_section}"</action>
    <action>Display: "   Name: {repo_name}"</action>
    <action>Display: "   Icon: {chosen_icon}"</action>
    <action>Display: "   Path: {project_path}"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Locate WezTerm config file (typically ~/.wezterm.lua or ~/.config/wezterm/wezterm.lua)</action>

    <check if="WezTerm config not found">
      <action>Display: "⚠️  WezTerm config not found. Skipping launcher update."</action>
      <action>Display: "Manual: Add project to WezTerm launcher in section '{inferred_section}'"</action>
      <goto step="8"/>
    </check>

    <action>Add project entry to WezTerm launcher config in {inferred_section} section</action>
    <action>Format: { label = "{repo_name}", icon = "{chosen_icon}", cwd = "{project_path}" }</action>
    <action>Display: "✓ Added to WezTerm launcher"</action>
    <action>Display: "   Section: {inferred_section}"</action>
    <action>Display: "   Reload WezTerm to see changes"</action>
  </check>

  <template-output>wezterm_updated, launcher_section</template-output>
</step>

<step n="8" goal="Sync templates from agent_customisations" if="skip_sync == false">
  <action>Fetch latest templates from {agent_customisations_repo}</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would sync templates from agent_customisations repo"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Clone or fetch agent_customisations repo to temp location</action>
    <action>Copy templates/agents/ → .bmad/expansion-packs/{bmad_module}/agents/</action>
    <action>Copy templates/workflows/ → .bmad/expansion-packs/{bmad_module}/workflows/</action>
    <action>Copy templates/preferences/ → .bmad/preferences/</action>
    <action>Preserve any *.local.* files (don't overwrite)</action>
    <action>Display: "✓ Templates synced from remote"</action>
  </check>

  <template-output>templates_synced</template-output>
</step>

<step n="9" goal="Final summary and next steps">
  <action>Generate setup completion summary</action>

  <action>Show statistics:</action>
  <action>- Directories created: [count]</action>
  <action>- Documents generated: [count]</action>
  <action>- Git initialized: [yes/no/skipped]</action>
  <action>- GitHub repo: [created/skipped]</action>
  <action>- Linear project: [created/skipped]</action>
  <action>- WezTerm launcher: [updated/skipped]</action>
  <action>- BMAD installed: [yes/no/skipped]</action>
  <action>- Templates synced: [yes/no/skipped]</action>

  <check if="dry_run_mode == true">
    <action>Display: "DRY RUN COMPLETE - No changes were made"</action>
    <action>Display: "Run without --dry-run to execute setup"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Display: "✓ Setup complete! Your project is fully wired and ready."</action>

    <action>Display configuration summary:</action>
    <action>- Git: [initialized/existing]</action>
    <action>- GitHub: [repo_url or "not created"]</action>
    <action>- Linear: [linear_project_url or "not created"]</action>
    <action>- WezTerm: [Added to {inferred_section} or "not updated"]</action>
    <action>- BMAD: [installed with {bmad_module} module]</action>

    <action>Display next steps:</action>
    <action>1. Review foundation documents (idea.md, CLAUDE.md, AGENTS.md)</action>
    <action>2. Start coding in src/</action>
    <action>3. Create first Linear issue to track initial development</action>
    <action>4. When frontend/backend ready, run: @Atlas *post-dev-setup</action>
    <action>5. For maintenance: @Atlas *maintenance</action>
  </check>

  <template-output>setup_complete, next_steps</template-output>
</step>

</workflow>
