# Brownfield Setup Workflow Instructions

<critical>The workflow execution engine is governed by: {project-root}/bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/bmad/bmm/workflows/brownfield-setup/workflow.yaml</critical>
<critical>Communicate in {communication_language} throughout the setup process</critical>

<workflow>

<step n="0" goal="Parse flags and determine execution mode">
  <action>Check for flags in user input: --dry-run, --skip-git, --skip-bmad, --skip-github, --skip-docs, --skip-linear, --skip-wezterm</action>
  <action>Set dry_run_mode = true if --dry-run flag present</action>
  <action>Set skip flags for each component</action>

  <check if="dry_run_mode == true">
    <action>Display: "DRY RUN MODE - Scanning project, no changes will be made"</action>
  </check>

  <template-output>execution_mode</template-output>
</step>

<step n="1" goal="Scan existing project and identify gaps">
  <action>Display: "🔍 Scanning project for missing setup components..."</action>

  <action>Detect what's already present:</action>
  <action>- Git: Check for .git directory</action>
  <action>- GitHub: Check for remote origin (git remote -v)</action>
  <action>- Linear: Check for .envrc/.env with LINEAR_TEAM and LINEAR_PROJECT variables (linctl stores auth separately)</action>
  <action>- WezTerm: Check if project path exists in WezTerm launcher config</action>
  <action>- BMAD: Check for .bmad/ directory and module installation</action>
  <action>- BMAD IDE integrations: Check for Claude Code, Codex, Gemini CLI configurations</action>
  <action>- Logging system: Check for .logs/ directory and .logs/all.ndjson</action>
  <action>- Service launcher: Check for ./service-start script</action>
  <action>- Foundation docs: Check for idea.md, CLAUDE.md, AGENTS.md, README.md</action>

  <action>Build gaps list - components that are missing:</action>
  <action>Store: missing_components = []</action>
  <action>For each component not detected, add to missing_components</action>

  <check if="missing_components is empty">
    <action>Display: "✓ Project is fully configured! No gaps detected."</action>
    <action>Display breakdown of what was found:</action>
    <action>- ✓ Git initialized</action>
    <action>- ✓ GitHub remote configured</action>
    <action>- ✓ Linear project linked</action>
    <action>- ✓ WezTerm launcher updated</action>
    <action>- ✓ BMAD installed</action>
    <action>- ✓ Logging system initialized</action>
    <action>- ✓ Service launcher present</action>
    <action>- ✓ Foundation documents present</action>
    <action>Exit workflow with success message</action>
  </check>

  <action>Display: "Found gaps in project setup:"</action>
  <action>For each missing component, display: "  ⚠️  Missing: {component}"</action>

  <template-output>detected_state, missing_components</template-output>
</step>

<step n="2" goal="Determine project type and configuration">
  <action>Detect project type from existing files:</action>
  <action>- Check for package.json → software (Node.js)</action>
  <action>- Check for requirements.txt → software (Python)</action>
  <action>- Check for go.mod → software (Go)</action>
  <action>- Check for *.ipynb → data-science</action>
  <action>- Check for docs/ with many .md files → writing</action>
  <action>- Otherwise → ask user</action>

  <check if="project_type cannot be auto-detected">
    <ask response="project_type">What type of project is this?
      1. Software Development
      2. Data Science / Analytics
      3. Writing / Documentation
      4. Business / Strategy
      5. Design / Creative
    </ask>
  </check>

  <action>Map project type to BMAD module (currently all → bmm)</action>

  <action>Infer Linear label and WezTerm section from parent directory:</action>
  <action>- Get parent directory name of current project</action>
  <action>- If parent contains "peps" (case-insensitive) → Label: "Peps Ventures", Section: "PEPS VENTURES"</action>
  <action>- If parent contains "raegis" (case-insensitive) → Label: "Raegis Labs", Section: "RAEGIS LABS"</action>
  <action>- If parent contains "xognosis" (case-insensitive) → Label: "XOGNOSIS", Section: "XOGNOSIS"</action>
  <action>- Otherwise → Label: "General", Section: "PROJECTS"</action>

  <template-output>project_type, bmad_module, inferred_label, inferred_section</template-output>
</step>

<step n="3" goal="Plan gap-filling execution">
  <action>Generate execution plan based on missing_components and skip flags</action>

  <action>Build step list for only missing components:</action>
  <action>If "git" in missing_components and skip_git == false → Add: Git initialization</action>
  <action>If "github" in missing_components and skip_github == false → Add: GitHub repo creation</action>
  <action>If "linear" in missing_components and skip_linear == false → Add: Linear project setup</action>
  <action>If "wezterm" in missing_components and skip_wezterm == false → Add: WezTerm launcher update</action>
  <action>If "bmad" in missing_components and skip_bmad == false → Add: BMAD installation</action>
  <action>If "logging" in missing_components and skip_logs == false → Add: Logging system initialization</action>
  <action>If "service_start" in missing_components and skip_logs == false → Add: Service launcher creation</action>
  <action>If "foundation_docs" in missing_components and skip_docs == false → Add: Foundation docs creation</action>

  <action>Display execution plan with estimated time</action>

  <check if="execution plan is empty">
    <action>Display: "All components either present or skipped. Nothing to do."</action>
    <action>Exit workflow</action>
  </check>

  <check if="dry_run_mode == false">
    <ask>Ready to fill these gaps? [y/n]</ask>
    <check if="user answered no">
      <action>Exit workflow with message: "Setup cancelled. Re-run when ready."</action>
    </check>
  </check>

  <template-output>execution_plan</template-output>
</step>

<step n="4" goal="Initialize Git repository" if="git in missing_components and skip_git == false">
  <action>Get current directory name for repo naming</action>
  <ask response="repo_name">Repository name (default: current directory)?</ask>

  <action>Detect project type for .gitignore template (from step 2)</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would initialize git and create .gitignore for {detected_type}"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Execute: git init</action>
    <action>Create .gitignore based on detected type</action>
    <action>Add common ignores: .bmad/, .claude/, .env, node_modules/, __pycache__/</action>
    <action>Display: "✓ Git initialized with .gitignore"</action>
  </check>

  <template-output>git_initialized, gitignore_type</template-output>
</step>

<step n="5" goal="Create foundation documents" if="foundation_docs in missing_components and skip_docs == false">
  <action>Check which specific docs are missing</action>
  <action>Only create missing docs, preserve existing ones</action>

  <ask response="project_description" if="idea.md missing">Briefly describe this project (1-2 sentences):</ask>

  <action if="idea.md missing">Create idea.md with structure from quick-setup</action>
  <action if="CLAUDE.md missing">Create CLAUDE.md with project context</action>
  <action if="AGENTS.md missing">Create AGENTS.md with agent config</action>
  <action if="README.md missing">Create README.md with setup instructions</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would create missing foundation docs: {list}"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Write missing foundation documents</action>
    <action>Display: "✓ Foundation documents created: {list}"</action>
  </check>

  <template-output>foundation_docs_created</template-output>
</step>

<step n="5.5" goal="Configure BMAD IDE integrations" if="(bmad in missing_components or bmad_ide_integrations in missing_components) and skip_bmad == false">
  <check if="bmad in missing_components">
    <action>Display: "BMAD will be installed. Configure IDE integrations:"</action>
  </check>

  <check if="bmad NOT in missing_components and bmad_ide_integrations in missing_components">
    <action>Display: "BMAD is installed but missing IDE integrations. Add them now:"</action>
  </check>

  <ask response="ide_integrations">Install BMAD IDE integrations (Claude Code, Codex, Gemini CLI)?

  1. Yes - Install all integrations (recommended)
  2. Custom - Choose which integrations to install
  3. No - Skip IDE integrations
  </ask>

  <check if="user answered 1 or yes">
    <action>Set: install_claude_code = true</action>
    <action>Set: install_codex = true</action>
    <action>Set: install_gemini = true</action>
    <action>Display: "Will install: Claude Code ✓ Codex ✓ Gemini CLI ✓"</action>
  </check>

  <check if="user answered 2 or custom">
    <ask response="install_claude_code">Install Claude Code integration? [y/n]</ask>
    <ask response="install_codex">Install OpenAI Codex integration? [y/n]</ask>
    <ask response="install_gemini">Install Gemini CLI integration? [y/n]</ask>

    <action>Display selected integrations:</action>
    <action if="install_claude_code == yes">- Claude Code ✓</action>
    <action if="install_codex == yes">- Codex ✓</action>
    <action if="install_gemini == yes">- Gemini CLI ✓</action>
  </check>

  <check if="user answered 3 or no">
    <action>Set: install_claude_code = false</action>
    <action>Set: install_codex = false</action>
    <action>Set: install_gemini = false</action>
    <action>Display: "Skipping IDE integrations"</action>
  </check>

  <template-output>bmad_ide_preferences</template-output>
</step>

<step n="6" goal="Install BMAD v6" if="bmad in missing_components and skip_bmad == false">
  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would install BMAD v6:"</action>
    <action>Display: "   Module: {bmad_module}"</action>
    <action if="install_claude_code == true">Display: "   + Claude Code integration"</action>
    <action if="install_codex == true">Display: "   + Codex integration"</action>
    <action if="install_gemini == true">Display: "   + Gemini CLI integration"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Build BMAD install command with IDE preferences from step 5.5</action>
    <action>Execute BMAD installation for {bmad_module}</action>
    <action>Create .bmad/ structure</action>
    <action>Initialize module configuration</action>
    <action>Set up agent symlinks</action>
    <action>Configure selected IDE integrations</action>
    <action>Sync templates from agent_customisations repo</action>
    <action>Display: "✓ BMAD v6 installed ({bmad_module} module)"</action>
    <action if="install_claude_code == true">Display: "   ✓ Claude Code integration configured"</action>
    <action if="install_codex == true">Display: "   ✓ Codex integration configured"</action>
    <action if="install_gemini == true">Display: "   ✓ Gemini CLI integration configured"</action>
  </check>

  <template-output>bmad_installed, bmad_module, ide_integrations_configured</template-output>
</step>

<step n="6.5" goal="Add missing IDE integrations to existing BMAD" if="bmad NOT in missing_components and bmad_ide_integrations in missing_components and skip_bmad == false">
  <action>Display: "Adding IDE integrations to existing BMAD installation..."</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would add IDE integrations:"</action>
    <action if="install_claude_code == true">Display: "   + Claude Code integration"</action>
    <action if="install_codex == true">Display: "   + Codex integration"</action>
    <action if="install_gemini == true">Display: "   + Gemini CLI integration"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Configure IDE integrations in existing BMAD installation</action>
    <action if="install_claude_code == true">Add Claude Code configuration to .bmad/</action>
    <action if="install_codex == true">Add Codex configuration to .bmad/</action>
    <action if="install_gemini == true">Add Gemini CLI configuration to .bmad/</action>
    <action>Display: "✓ IDE integrations added to BMAD"</action>
    <action if="install_claude_code == true">Display: "   ✓ Claude Code integration configured"</action>
    <action if="install_codex == true">Display: "   ✓ Codex integration configured"</action>
    <action if="install_gemini == true">Display: "   ✓ Gemini CLI integration configured"</action>
  </check>

  <template-output>ide_integrations_added</template-output>
</step>

<step n="7" goal="Create GitHub repository" if="github in missing_components and skip_github == false">
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
    <action>If git has commits → git push -u origin main</action>
    <action>Else → Stage and commit existing files first, then push</action>
    <action>Get repository URL from gh</action>
    <action>Display: "✓ GitHub repository created and pushed"</action>
    <action>Display: "   URL: {repo_url}"</action>
  </check>

  <template-output>github_created, repo_url</template-output>
</step>

<step n="8" goal="Fill missing Linear environment variables" if="linear in missing_components and skip_linear == false">
  <action>Verify required tools installed</action>

  <check if="linctl not installed">
    <action>Display: "⚠️  linctl not found. Skipping Linear integration."</action>
    <action>Display: "Install with: brew install linctl (or see https://github.com/dorkitude/linctl)"</action>
    <goto step="9"/>
  </check>

  <ask response="configure_linear">Add missing Linear environment variables? [y/n]</ask>

  <check if="user answered no">
    <action>Display: "Skipping Linear integration"</action>
    <goto step="9"/>
  </check>

  <action>Check which Linear env vars are already set in .envrc/.env (if file exists)</action>
  <action>Store: missing_linear_vars = []</action>
  <action>If LINEAR_TEAM not present → Add to missing_linear_vars</action>
  <action>If LINEAR_PROJECT not present → Add to missing_linear_vars</action>

  <check if="missing_linear_vars is empty">
    <action>Display: "✓ Linear environment variables already configured"</action>
    <action>Display: "   LINEAR_TEAM: Set"</action>
    <action>Display: "   LINEAR_PROJECT: Set"</action>
    <goto step="9"/>
  </check>

  <action>Display: "Missing Linear variables: {missing_linear_vars}"</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would add missing Linear env vars:"</action>
    <action if="LINEAR_TEAM in missing_linear_vars">Display: "   + LINEAR_TEAM (Raegis Labs)"</action>
    <action if="LINEAR_PROJECT in missing_linear_vars">Display: "   + LINEAR_PROJECT (match or create)"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Step 1: Verify linctl authentication</action>
    <action>Execute: linctl auth status</action>
    <check if="not authenticated">
      <action>Display: "⚠️  linctl not authenticated. Run: linctl auth"</action>
      <action>Display: "Skipping Linear configuration steps"</action>
      <goto step="9"/>
    </check>

    <action>Step 2: Set LINEAR_TEAM if missing</action>
    <check if="LINEAR_TEAM in missing_linear_vars">
      <action>Set: linear_team_id = "b8ff8916-3e03-435d-809f-9d45ef4199c8"</action>
      <action>Set: linear_team_key = "RAE"</action>
    </check>

    <action>Step 3: Find or create LINEAR_PROJECT if missing</action>
    <check if="LINEAR_PROJECT in missing_linear_vars">
      <action>Get repo name from git remote or current directory</action>
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
    </check>

    <action>Step 4: Write/update .envrc with missing variables only</action>
    <action>Append or update entries:</action>
    <action if="LINEAR_TEAM in missing_linear_vars">export LINEAR_TEAM=RAE</action>
    <action if="LINEAR_PROJECT in missing_linear_vars">export LINEAR_PROJECT={linear_project_id}</action>
    <action>Execute if direnv installed: direnv allow</action>

    <action>Display: "✓ Linear environment configured (linctl)"</action>
    <action>Display: "   Updated variables:"</action>
    <action if="LINEAR_TEAM in missing_linear_vars">Display: "     + LINEAR_TEAM: RAE (Raegis Labs)"</action>
    <action if="LINEAR_PROJECT in missing_linear_vars">Display: "     + LINEAR_PROJECT: {repo_name}"</action>
    <action>Display: "   Preserved existing variables:"</action>
    <action if="LINEAR_TEAM not in missing_linear_vars">Display: "     ✓ LINEAR_TEAM (already set)"</action>
    <action if="LINEAR_PROJECT not in missing_linear_vars">Display: "     ✓ LINEAR_PROJECT (already set)"</action>
  </check>

  <template-output>linear_env_configured, missing_linear_vars_added</template-output>
</step>

<step n="9" goal="Add to WezTerm launcher" if="wezterm in missing_components and skip_wezterm == false">
  <action>Use inferred_section from step 2</action>

  <action>Choose icon based on project type:</action>
  <action>- software → 💻</action>
  <action>- data-science → 📊</action>
  <action>- writing → 📝</action>
  <action>- business → 💼</action>
  <action>- design → 🎨</action>
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
    <action>Locate WezTerm config file</action>

    <check if="WezTerm config not found">
      <action>Display: "⚠️  WezTerm config not found. Skipping launcher update."</action>
      <action>Display: "Manual: Add project to WezTerm launcher in section '{inferred_section}'"</action>
      <goto step="10"/>
    </check>

    <action>Add project entry to WezTerm launcher config in {inferred_section} section</action>
    <action>Display: "✓ Added to WezTerm launcher"</action>
    <action>Display: "   Section: {inferred_section}"</action>
    <action>Display: "   Reload WezTerm to see changes"</action>
  </check>

  <template-output>wezterm_updated, launcher_section</template-output>
</step>

<step n="9.3" goal="Initialize logging system" if="logging in missing_components and skip_logs == false">
  <action>Display: "Initializing unified logging system..."</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would run /logs.init to create:"</action>
    <action>Display: "   .logs/all.ndjson - Main log file"</action>
    <action>Display: "   .logs/inbox/ - Log sources directory"</action>
    <action>Display: "   .logs/redaction.json - Redaction config"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Run: /logs.init</action>
    <action>Verify .logs/ directory was created</action>

    <check if="logging created successfully">
      <action>Display: "✓ Logging system initialized via /logs.init"</action>
      <action>Display: "   Main log: .logs/all.ndjson"</action>
      <action>Display: "   Inbox: .logs/inbox/"</action>
    </check>

    <check if="logging creation failed">
      <action>Display: "⚠️  /logs.init failed"</action>
      <action>Display: "   You can manually initialize with: mkdir -p .logs/inbox && touch .logs/all.ndjson"</action>
    </check>
  </check>

  <template-output>logging_initialized</template-output>
</step>

<step n="9.5" goal="Create service-start launcher" if="service_start in missing_components and skip_logs == false">
  <action>Display: "Creating service-start launcher..."</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would create ./service-start script"</action>
    <action>Display: "   - Unified launcher for frontend/backend/docker services"</action>
    <action>Display: "   - Integrates with logging aggregator"</action>
    <action>Display: "   - Uses dynamic ports from .envrc"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Run: /init-service-start</action>
    <action>Verify ./service-start was created</action>

    <check if="service-start created successfully">
      <action>Make service-start executable: chmod +x ./service-start</action>
      <action>Display: "✓ service-start launcher created"</action>
      <action>Display: "   Run with: ./service-start start"</action>
      <action>Display: "   Commands: start, stop, restart, status"</action>
    </check>

    <check if="service-start creation failed">
      <action>Display: "⚠️  Failed to create service-start"</action>
      <action>Display: "   You can manually copy it from the agent_customisations repo"</action>
    </check>
  </check>

  <template-output>service_start_created</template-output>
</step>

<step n="10" goal="Final summary and next steps">
  <action>Generate gap-filling completion summary</action>

  <action>Show before/after state:</action>
  <action>Display: "Gap-filling complete! Here's what changed:"</action>
  <action>For each component that was added, show: "  ✓ Added: {component}"</action>
  <action>For each component that was already present, show: "  ✓ Existing: {component}"</action>

  <check if="dry_run_mode == true">
    <action>Display: "DRY RUN COMPLETE - No changes were made"</action>
    <action>Display: "Run without --dry-run to execute gap-filling"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Display: "✓ Project fully configured and ready!"</action>

    <action>Display configuration summary:</action>
    <action>- Git: {initialized/existing}</action>
    <action>- GitHub: {repo_url or "not created"}</action>
    <action>- Linear: {linear_project_url or "not created"}</action>
    <action>- WezTerm: {Added to {inferred_section} or "not updated"}</action>
    <action>- BMAD: {installed with {bmad_module} module}</action>
    <action>- IDE Integrations: {list of configured integrations or "not configured"}</action>
    <action>- Logging system: {initialized/existing}</action>
    <action>- Service launcher: {created/existing}</action>
    <action>- Foundation docs: {created/existing}</action>

    <action>Display next steps:</action>
    <action>1. Review any newly created documentation</action>
    <action if="service_start in components_added">2. Start services with: ./service-start start</action>
    <action if="logging in components_added and service_start not in components_added">2. View logs: tail -f .logs/all.ndjson</action>
    <action>3. Continue development with full tooling support</action>
    <action>4. Use @Atlas *maintenance for ongoing project health</action>
  </check>

  <template-output>gap_filling_complete, next_steps</template-output>
</step>

</workflow>
