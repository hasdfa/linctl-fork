# Health Check Workflow Instructions

<critical>The workflow execution engine is governed by: {project-root}/bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/bmad/bmm/workflows/health-check/workflow.yaml</critical>
<critical>Communicate in {communication_language} throughout the diagnostic process</critical>
<critical>This workflow is CHATTY - provide detailed explanations, context, and recommendations</critical>

<workflow>

<step n="1" goal="Initialize diagnostic and explain process">
  <action>Display: "🏥 Project Health Check"</action>
  <action>Display: "Running comprehensive diagnostics on your project..."</action>
  <action>Display: ""</action>
  <action>Display: "I'll check:"</action>
  <action>Display: "  • Git configuration and status"</action>
  <action>Display: "  • GitHub integration"</action>
  <action>Display: "  • Linear project health"</action>
  <action>Display: "  • WezTerm launcher setup"</action>
  <action>Display: "  • BMAD installation and IDE integrations"</action>
  <action>Display: "  • Dependencies freshness"</action>
  <action>Display: "  • Test coverage"</action>
  <action>Display: "  • Documentation completeness"</action>
  <action>Display: "  • Logging system"</action>
  <action>Display: "  • Service launcher"</action>
  <action>Display: "  • Maintenance schedule"</action>
  <action>Display: ""</action>

  <action>Initialize findings storage:</action>
  <action>findings_critical = []</action>
  <action>findings_warning = []</action>
  <action>findings_info = []</action>
  <action>auto_fixes_available = []</action>

  <template-output>diagnostic_initialized</template-output>
</step>

<step n="2" goal="Check Git status and configuration" if="skip_git == false">
  <action>Display: "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"</action>
  <action>Display: "📁 Git Status"</action>
  <action>Display: ""</action>

  <action>Check if git repository exists:</action>
  <action>Execute: test -d .git && echo "yes" || echo "no"</action>

  <check if="git not initialized">
    <action>Add to findings_critical: {
      "category": "Git",
      "issue": "Git repository not initialized",
      "explanation": "This project is not under version control. Git is essential for tracking changes, collaborating, and deploying.",
      "impact": "Cannot track history, no backup, deployment will fail",
      "fix": "Initialize Git repository",
      "auto_fix_available": true
    }</action>
    <action>Display: "🔴 CRITICAL: Git not initialized"</action>
    <action>Display: "   Without Git, you can't track changes or deploy to most platforms."</action>
    <goto step="3"/>
  </check>

  <action>Check for uncommitted changes:</action>
  <action>Execute: git status --porcelain</action>
  <action>Count modified, staged, untracked files</action>

  <check if="uncommitted changes exist">
    <action>Add to findings_warning: {
      "category": "Git",
      "issue": "{count} uncommitted changes",
      "explanation": "You have uncommitted work. Regular commits help track progress and prevent data loss.",
      "impact": "Risk of losing work, harder to collaborate",
      "fix": "Review and commit changes: git add . && git commit -m 'message'",
      "auto_fix_available": false
    }</action>
    <action>Display: "🟡 WARNING: {count} uncommitted changes"</action>
    <action>Display: "   Consider committing regularly to track your progress."</action>
  </check>

  <action>Check current branch:</action>
  <action>Execute: git branch --show-current</action>

  <check if="on main/master with uncommitted changes">
    <action>Add to findings_warning: {
      "category": "Git",
      "issue": "Working directly on {branch}",
      "explanation": "It's better to use feature branches for development work, keeping main/master stable.",
      "impact": "Makes collaboration harder, main branch becomes unstable",
      "fix": "Create feature branch: git checkout -b feature/your-feature",
      "auto_fix_available": false
    }</action>
    <action>Display: "🟡 WARNING: Working directly on {branch}"</action>
    <action>Display: "   Consider using feature branches for development."</action>
  </check>

  <action>Display: "✓ Git configured and operational"</action>

  <template-output>git_checked</template-output>
</step>

<step n="3" goal="Check GitHub integration" if="skip_git == false">
  <action>Display: ""</action>
  <action>Display: "🐙 GitHub Integration"</action>
  <action>Display: ""</action>

  <action>Check for GitHub remote:</action>
  <action>Execute: git remote get-url origin 2>/dev/null || echo "none"</action>

  <check if="no remote origin">
    <action>Add to findings_info: {
      "category": "GitHub",
      "issue": "No GitHub remote configured",
      "explanation": "Your project isn't connected to GitHub. This means no backup, no collaboration, and deployment challenges.",
      "impact": "No remote backup, can't collaborate easily",
      "fix": "Create GitHub repo: gh repo create {project-name} --private --source=.",
      "auto_fix_available": true
    }</action>
    <action>Display: "🔵 INFO: No GitHub remote configured"</action>
    <action>Display: "   Your code is only local. Consider pushing to GitHub for backup."</action>
    <goto step="4"/>
  </check>

  <action>Check if local is ahead of remote:</action>
  <action>Execute: git rev-list --left-right --count origin/main...HEAD 2>/dev/null || echo "0 0"</action>
  <action>Parse: commits_behind commits_ahead</action>

  <check if="commits_ahead > 0">
    <action>Add to findings_warning: {
      "category": "GitHub",
      "issue": "{commits_ahead} unpushed commits",
      "explanation": "You have local commits that aren't backed up to GitHub. If your machine fails, this work is lost.",
      "impact": "Risk of losing work, team can't see your progress",
      "fix": "Push to GitHub: git push",
      "auto_fix_available": true
    }</action>
    <action>Display: "🟡 WARNING: {commits_ahead} commits not pushed to GitHub"</action>
    <action>Display: "   Push soon to back up your work."</action>
  </check>

  <check if="commits_behind > 0">
    <action>Add to findings_info: {
      "category": "GitHub",
      "issue": "{commits_behind} commits behind remote",
      "explanation": "Your local branch is behind the remote. You might be missing updates from collaborators.",
      "impact": "May have merge conflicts later, missing latest changes",
      "fix": "Pull latest changes: git pull",
      "auto_fix_available": true
    }</action>
    <action>Display: "🔵 INFO: {commits_behind} commits behind remote"</action>
    <action>Display: "   Pull the latest changes when ready."</action>
  </check>

  <action>Display: "✓ GitHub remote configured and synced"</action>

  <template-output>github_checked</template-output>
</step>

<step n="4" goal="Check Linear integration and environment variables" if="skip_linear == false">
  <action>Display: ""</action>
  <action>Display: "📋 Linear Project Integration"</action>
  <action>Display: ""</action>

  <action>Check for .envrc or .env:</action>
  <action>Execute: test -f .envrc && echo "envrc" || (test -f .env && echo "env" || echo "none")</action>

  <check if=".envrc/.env does not exist">
    <action>Add to findings_critical: {
      "category": "Linear",
      "issue": ".envrc/.env missing",
      "explanation": "No .envrc or .env found. Linear helpers expect LINEAR_TEAM and LINEAR_PROJECT to be set (linctl uses its own auth).",
      "impact": "Cannot connect to Linear API, no issue tracking",
      "fix": "Run @Atlas *brownfield-setup to write .envrc with LINEAR_TEAM and LINEAR_PROJECT",
      "auto_fix_available": true
    }</action>
    <action>Display: "🔴 CRITICAL: .envrc/.env not found"</action>
    <action>Display: "   Linear helpers expect LINEAR_TEAM and LINEAR_PROJECT."</action>
    <action>Display: "   Fix: Run @Atlas *brownfield-setup to add .envrc"</action>
    <goto step="5"/>
  </check>

  <action>Check linctl authentication:</action>
  <action>Execute: linctl auth status</action>
  <action>Store result as: auth_status</action>

  <check if="auth_status indicates not authenticated">
    <action>Add to findings_critical: {
      "category": "Linear",
      "issue": "linctl not authenticated",
      "explanation": "linctl requires a one-time authentication (stored in ~/.linctl-auth.json)",
      "impact": "Cannot access Linear data",
      "fix": "Run: linctl auth",
      "auto_fix_available": false
    }</action>
    <action>Display: "🔴 CRITICAL: linctl not authenticated (run 'linctl auth')"</action>
  </check>

  <action>Check Linear env in .envrc/.env (for helpers):</action>
  <action>Execute: bash -lc 'source .envrc 2>/dev/null || source .env 2>/dev/null || true; echo ${LINEAR_TEAM:-}'</action>
  <action>Store as: configured_team_key</action>
  <action>Execute: bash -lc 'source .envrc 2>/dev/null || source .env 2>/dev/null || true; echo ${LINEAR_PROJECT:-}'</action>
  <action>Store as: configured_project_id</action>

  <note>linctl stores credentials in ~/.linctl-auth.json; no .env file permissions required for API key</note>

  <check if="configured_team_key == '' or configured_project_id == ''">
    <action>Add to findings_critical: {
      "category": "Linear",
      "issue": "Missing Linear environment variables (helpers)",
      "explanation": "Missing LINEAR_TEAM and/or LINEAR_PROJECT in .envrc/.env. linctl auth is separate and stored in ~/.linctl-auth.json.",
      "impact": "Cannot attach issues to project by default",
      "fix": "Run @Atlas *brownfield-setup to add .envrc with LINEAR_TEAM and LINEAR_PROJECT",
      "auto_fix_available": true
    }</action>
    <action>Display: "🔴 CRITICAL: Missing LINEAR_TEAM and/or LINEAR_PROJECT"</action>
    <action>Display: "   Fix: Run @Atlas *brownfield-setup"</action>
    <goto step="5"/>
  </check>

  <action>Display: "Environment Variables (helpers):"</action>
  <action>Display: "   ✓ LINEAR_TEAM set"</action>
  <action>Display: "   ✓ LINEAR_PROJECT set"</action>
  <action>Display: ""</action>

  <action>Validate Linear API connectivity:</action>
  <action>Display: "Testing Linear API connection..."</action>

  <check if="linctl not installed">
    <action>Add to findings_warning: {
      "category": "Linear",
      "issue": "linctl not installed",
      "explanation": "Environment variables are set but linctl CLI is not installed. Cannot verify API connectivity.",
      "impact": "Cannot validate Linear integration",
      "fix": "Install linctl: brew install linctl",
      "auto_fix_available": false
    }</action>
    <action>Display: "🟡 WARNING: linctl not installed"</action>
    <action>Display: "   Install: brew install linctl"</action>
    <goto step="5"/>
  </check>

  <action>Test Linear connectivity via linctl:</action>
  <action>Execute: linctl team list 2>&1</action>
  <action>Store result as: team_list_result</action>

  <check if="team_list_result contains error or authentication failed">
    <action>Add to findings_critical: {
      "category": "Linear",
      "issue": "Linear API authentication failed",
      "explanation": "linctl authentication failed. You may need to re-run 'linctl auth'.",
      "impact": "Cannot access Linear data",
      "fix": "Run: linctl auth",
      "auto_fix_available": true
    }</action>
    <action>Display: "🔴 CRITICAL: Linear authentication failed (linctl)"</action>
    <action>Display: "   Fix: Run 'linctl auth'"</action>
    <action>Display: "   Tip: Use @Atlas *brownfield-setup to configure .envrc (LINEAR_TEAM/LINEAR_PROJECT)"</action>
    <goto step="5"/>
  </check>

  <action>Parse team list to verify team access:</action>
  <action>Extract team count and names from team_list_result</action>
  <action>Display: "   ✓ API connection successful ({team_count} teams accessible)"</action>

  <action>Verify LINEAR_TEAM matches Raegis Labs (key RAE):</action>
  <action>Execute: bash -lc 'source .envrc 2>/dev/null || source .env 2>/dev/null || true; echo ${LINEAR_TEAM:-}'</action>
  <action>Store as: configured_team_key</action>
  <action>Expected team key: RAE</action>

  <check if="configured_team_key != expected_team_key">
    <action>Add to findings_warning: {
      "category": "Linear",
      "issue": "LINEAR_TEAM mismatch",
      "explanation": "LINEAR_TEAM is set to {configured_team_key} but should be RAE (Raegis Labs team key).",
      "impact": "May be accessing wrong team",
      "fix": "Update LINEAR_TEAM in .envrc/.env",
      "auto_fix_available": true
    }</action>
    <action>Display: "🟡 WARNING: LINEAR_TEAM mismatch"</action>
    <action>Display: "   Current: {configured_team_key}"</action>
    <action>Display: "   Expected: RAE (Raegis Labs)"</action>
  </check>

  <action>Verify LINEAR_PROJECT exists:</action>
  <action>Execute: bash -lc 'source .envrc 2>/dev/null || source .env 2>/dev/null || true; linctl project get ${LINEAR_PROJECT:-} --json 2>&1'</action>
  <action>Store result as: project_result</action>

  <check if="project_result contains error or not found">
    <action>Add to findings_warning: {
      "category": "Linear",
      "issue": "LINEAR_PROJECT not found",
      "explanation": "LINEAR_PROJECT is set but the project doesn't exist in Linear. It may have been deleted or archived.",
      "impact": "Cannot track issues in this project",
      "fix": "Run @Atlas *brownfield-setup to match or create project",
      "auto_fix_available": true
    }</action>
    <action>Display: "🟡 WARNING: LINEAR_PROJECT not found in Linear"</action>
    <action>Display: "   Project may have been deleted or archived"</action>
    <action>Display: "   Fix: Run @Atlas *brownfield-setup"</action>
  </check>

  <check if="project_result successful">
    <action>Parse project name from project_result</action>
    <action>Display: "   ✓ Project exists: {project_name}"</action>
  </check>

  <action>Display: ""</action>
  <action>Display: "✓ Linear integration healthy"</action>

  <template-output>linear_checked</template-output>
</step>

<step n="5" goal="Check WezTerm launcher" if="skip_wezterm == false">
  <action>Display: ""</action>
  <action>Display: "🖥️  WezTerm Launcher"</action>
  <action>Display: ""</action>

  <action>Find WezTerm config and search for current project:</action>
  <action>Execute: bash -c 'WEZTERM_CONFIG=""; if [ -f ~/.wezterm.lua ]; then WEZTERM_CONFIG=~/.wezterm.lua; elif [ -f ~/.config/wezterm/wezterm.lua ]; then WEZTERM_CONFIG=~/.config/wezterm/wezterm.lua; fi; if [ -z "$WEZTERM_CONFIG" ]; then echo "config_not_found"; exit 0; fi; CURRENT_FULL=$(pwd); CURRENT_TILDE=$(pwd | sed "s|^$HOME|~|"); if grep -qE "$CURRENT_FULL|$CURRENT_TILDE" "$WEZTERM_CONFIG"; then echo "found|$WEZTERM_CONFIG"; else echo "not_found|$WEZTERM_CONFIG"; fi'</action>
  <action>Store result as: launcher_check_result</action>

  <check if="launcher_check_result contains config_not_found">
    <action>Display: "🔵 INFO: WezTerm config not found, skipping launcher check"</action>
    <goto step="6"/>
  </check>

  <action>Parse result:</action>
  <action>Extract status (found or not_found) from launcher_check_result before |</action>
  <action>Extract config_path from launcher_check_result after |</action>

  <check if="status == not_found">
    <action>Add to findings_info: {
      "category": "WezTerm",
      "issue": "Project not in WezTerm launcher",
      "explanation": "Adding to launcher makes it faster to navigate to this project.",
      "impact": "Minor convenience issue",
      "fix": "Add via quick-setup or brownfield-setup",
      "auto_fix_available": true
    }</action>
    <action>Display: "Config: {config_path}"</action>
    <action>Display: "Finding: 🔵 INFO — Project not in WezTerm launcher. Add via *brownfield-setup for quick access."</action>
  </check>

  <check if="status == found">
    <action>Display: "Config: {config_path}"</action>
    <action>Display: "Finding: ✓ Project in WezTerm launcher"</action>
  </check>

  <template-output>wezterm_checked</template-output>
</step>

<step n="6" goal="Check BMAD installation and IDE integrations" if="skip_bmad == false">
  <action>Display: ""</action>
  <action>Display: "🎯 BMAD Installation"</action>
  <action>Display: ""</action>

  <action>Detect BMAD installation (v6 or legacy):</action>
  <action>Execute: bash -c 'if [ -d bmad ]; then echo "v6"; elif [ -d .bmad ]; then echo "legacy"; else echo "none"; fi'</action>
  <action>Store result as: bmad_version</action>

  <check if="bmad_version == none">
    <action>Add to findings_info: {
      "category": "BMAD",
      "issue": "BMAD not installed",
      "explanation": "BMAD provides agent workflows for project management. Install it to access powerful automation.",
      "impact": "Missing automation capabilities",
      "fix": "Install BMAD via quick-setup or brownfield-setup",
      "auto_fix_available": true
    }</action>
    <action>Display: "Finding: 🔵 INFO — BMAD not installed. Install via *quick-setup or *brownfield-setup."</action>
    <goto step="7"/>
  </check>

  <action>Set IDE config path based on version:</action>
  <action if="bmad_version == v6">Set: ide_config_path = "bmad/_cfg/ides"</action>
  <action if="bmad_version == legacy">Set: ide_config_path = ".bmad"</action>

  <action>Check for IDE integration configs:</action>
  <action>Execute: bash -c 'IDE_COUNT=0; if [ -f {ide_config_path}/claude-code.yaml ] || [ -f {ide_config_path}/claude-code.yml ]; then echo -n "claude-code,"; IDE_COUNT=$((IDE_COUNT+1)); fi; if [ -f {ide_config_path}/codex.yaml ] || [ -f {ide_config_path}/codex.yml ]; then echo -n "codex,"; IDE_COUNT=$((IDE_COUNT+1)); fi; if [ -f {ide_config_path}/gemini.yaml ] || [ -f {ide_config_path}/gemini.yml ]; then echo -n "gemini,"; IDE_COUNT=$((IDE_COUNT+1)); fi; if [ "$IDE_COUNT" -eq 0 ]; then echo "none"; fi' | sed 's/,$//'</action>
  <action>Store result as: configured_ides</action>

  <check if="configured_ides == none">
    <action>Add to findings_warning: {
      "category": "BMAD",
      "issue": "No IDE integrations configured",
      "explanation": "BMAD is installed but not connected to your IDEs. You're missing out on AI-powered development assistance.",
      "impact": "Can't use AI agents from your IDE",
      "fix": "Configure IDE integrations via brownfield-setup",
      "auto_fix_available": true
    }</action>
    <action>Display: "Version: BMAD {bmad_version}"</action>
    <action>Display: "Finding: 🟡 WARNING — BMAD installed but no IDE integrations configured."</action>
  </check>

  <check if="configured_ides != none">
    <action>Parse configured_ides into list (comma-separated)</action>
    <action>Display: "Version: BMAD {bmad_version}"</action>
    <action>Display: "Finding: ✓ BMAD installed with IDE integrations:"</action>
    <action if="configured_ides contains claude-code">Display: "   • Claude Code ✓"</action>
    <action if="configured_ides contains codex">Display: "   • Codex ✓"</action>
    <action if="configured_ides contains gemini">Display: "   • Gemini CLI ✓"</action>

    <check if="not all three configured">
      <action>Count configured IDEs</action>
      <action>Display: "   ({count}/3 integrations configured)"</action>
    </check>
  </check>

  <template-output>bmad_checked</template-output>
</step>

<step n="7" goal="Check dependencies freshness" if="skip_deps == false">
  <action>Display: ""</action>
  <action>Display: "📦 Dependencies"</action>
  <action>Display: ""</action>

  <action>Detect package manager:</action>
  <action>Check for package.json → npm/yarn/pnpm</action>
  <action>Check for requirements.txt → pip</action>
  <action>Check for go.mod → go modules</action>
  <action>Check for Gemfile → bundler</action>

  <check if="no package manager detected">
    <action>Display: "🔵 INFO: No package manager detected, skipping dependency check"</action>
    <goto step="8"/>
  </check>

  <check if="package.json exists">
    <action>Check for outdated npm packages:</action>
    <action>Execute: npm outdated --json 2>/dev/null || echo "{}"</action>
    <action>Count outdated packages, categorize by severity</action>

    <check if="outdated packages > 10">
      <action>Add to findings_warning: {
        "category": "Dependencies",
        "issue": "{count} outdated npm packages",
        "explanation": "Many packages are outdated. This can cause security vulnerabilities and compatibility issues.",
        "impact": "Security risks, potential bugs, harder to upgrade later",
        "fix": "Update packages: npm update (safe) or npm update --latest (major versions)",
        "auto_fix_available": true
      }</action>
      <action>Display: "🟡 WARNING: {count} outdated npm packages"</action>
      <action>Display: "   Run 'npm outdated' to see details, 'npm update' to update safely."</action>
    </check>

    <check if="outdated packages <= 10 and > 0">
      <action>Display: "🔵 INFO: {count} packages have updates available"</action>
      <action>Display: "   Run 'npm outdated' to review."</action>
    </check>

    <check if="outdated packages == 0">
      <action>Display: "✓ All npm packages up to date"</action>
    </check>
  </check>

  <check if="requirements.txt exists">
    <action>Check for outdated Python packages:</action>
    <action>Execute: pip list --outdated --format=json 2>/dev/null || echo "[]"</action>
    <action>Count outdated packages</action>

    <check if="outdated packages > 5">
      <action>Add to findings_warning: {
        "category": "Dependencies",
        "issue": "{count} outdated Python packages",
        "explanation": "Your Python dependencies are outdated. This can cause security issues.",
        "impact": "Security vulnerabilities, compatibility problems",
        "fix": "Update packages: pip install --upgrade -r requirements.txt",
        "auto_fix_available": true
      }</action>
      <action>Display: "🟡 WARNING: {count} outdated Python packages"</action>
      <action>Display: "   Run 'pip list --outdated' to see which ones."</action>
    </check>

    <check if="outdated packages <= 5 and > 0">
      <action>Display: "🔵 INFO: {count} Python packages have updates"</action>
    </check>

    <check if="outdated packages == 0">
      <action>Display: "✓ All Python packages up to date"</action>
    </check>
  </check>

  <template-output>dependencies_checked</template-output>
</step>

<step n="8" goal="Check test coverage" if="skip_tests == false">
  <action>Display: ""</action>
  <action>Display: "🧪 Tests"</action>
  <action>Display: ""</action>

  <action>Search for test files and directories:</action>
  <action>Execute: bash -c 'TEST_COUNT=0; if [ -d tests ] || [ -d test ] || [ -d __tests__ ]; then TEST_COUNT=$((TEST_COUNT + 1)); fi; FILE_COUNT=$(find . -type f \( -name "*test*.js" -o -name "*test*.ts" -o -name "*test*.py" -o -name "*_test.go" -o -name "*spec*.js" -o -name "*spec*.ts" \) -not -path "*/node_modules/*" -not -path "*/.git/*" -not -path "*/venv/*" -not -path "*/.venv/*" 2>/dev/null | wc -l | xargs); TOTAL=$((TEST_COUNT + FILE_COUNT)); if [ "$TOTAL" -eq 0 ]; then echo "none"; else echo "$FILE_COUNT"; fi'</action>
  <action>Store result as: test_file_count</action>

  <check if="test_file_count == none">
    <action>Add to findings_critical: {
      "category": "Tests",
      "issue": "No tests found",
      "explanation": "Your project has no automated tests. Tests catch bugs early and make refactoring safer.",
      "impact": "Higher bug risk, harder to maintain, scary deployments",
      "fix": "Add test framework and write initial tests",
      "auto_fix_available": false
    }</action>
    <action>Display: "Finding: 🔴 CRITICAL — No test files found. Add a test framework and initial tests."</action>
    <goto step="9"/>
  </check>

  <action>Display: "Finding: ✓ Found {test_file_count} test files"</action>

  <action>Check for test coverage configuration:</action>
  <check if="package.json exists">
    <action>Execute: grep -q '"test:coverage"' package.json 2>/dev/null && echo "yes" || echo "no"</action>
    <action>Store as: has_npm_coverage</action>
    <check if="has_npm_coverage == yes">
      <action>Display: "   Run 'npm run test:coverage' to see detailed coverage."</action>
    </check>
  </check>

  <check if="pytest available">
    <action>Execute: command -v pytest >/dev/null 2>&1 && ([ -f pytest.ini ] || [ -f setup.cfg ] || [ -f pyproject.toml ]) && echo "yes" || echo "no"</action>
    <action>Store as: has_pytest_coverage</action>
    <check if="has_pytest_coverage == yes">
      <action>Display: "   Run 'pytest --cov' to see detailed coverage."</action>
    </check>
  </check>

  <template-output>tests_checked</template-output>
</step>

<step n="9" goal="Check documentation completeness" if="skip_docs == false">
  <action>Display: ""</action>
  <action>Display: "📚 Documentation"</action>
  <action>Display: ""</action>

  <action>Check for essential documentation files:</action>
  <action>Check for README.md</action>
  <action>Check for idea.md</action>
  <action>Check for CLAUDE.md</action>
  <action>Check for AGENTS.md</action>

  <action>missing_docs = []</action>

  <check if="README.md missing">
    <action>Add to missing_docs: "README.md"</action>
    <action>Add to findings_warning: {
      "category": "Documentation",
      "issue": "README.md missing",
      "explanation": "README is the entry point for anyone looking at your project. It should explain what it does and how to use it.",
      "impact": "Confusing for collaborators and future you",
      "fix": "Create README.md with project overview, setup, and usage",
      "auto_fix_available": true
    }</action>
  </check>

  <check if="idea.md missing">
    <action>Add to missing_docs: "idea.md"</action>
  </check>

  <check if="CLAUDE.md missing">
    <action>Add to missing_docs: "CLAUDE.md"</action>
  </check>

  <check if="missing_docs not empty">
    <action>Display: "🟡 WARNING: Missing documentation files: {missing_docs}"</action>
    <action>Display: "   Add these via quick-setup or brownfield-setup to improve project clarity."</action>
  </check>

  <check if="missing_docs is empty">
    <action>Display: "✓ Essential documentation files present"</action>
  </check>

  <action>Check README.md freshness (if exists):</action>
  <check if="README.md exists">
    <action>Get last modified date of README.md</action>
    <action>Calculate days since last update</action>

    <check if="days_since_update > 90">
      <action>Add to findings_info: {
        "category": "Documentation",
        "issue": "README.md not updated in {days_since_update} days",
        "explanation": "Documentation might be outdated. Keep README current with major project changes.",
        "impact": "Misleading information, harder onboarding",
        "fix": "Review and update README.md",
        "auto_fix_available": false
      }</action>
      <action>Display: "🔵 INFO: README.md last updated {days_since_update} days ago"</action>
      <action>Display: "   Consider reviewing if it's still accurate."</action>
    </check>
  </check>

  <template-output>documentation_checked</template-output>
</step>

<step n="10" goal="Check logging system" if="skip_logs == false">
  <action>Display: ""</action>
  <action>Display: "📝 Logging System"</action>
  <action>Display: ""</action>

  <action>Check if unified logging is set up:</action>
  <action>Look for .logs/ directory</action>
  <action>Check for .logs/all.ndjson</action>
  <action>Check for .logs/inbox/</action>

  <check if="logging not set up">
    <action>Add to findings_info: {
      "category": "Logging",
      "issue": "Unified logging not configured",
      "explanation": "Centralized logging makes debugging much easier. All service logs in one place helps AI agents help you.",
      "impact": "Harder debugging, scattered logs",
      "fix": "Run /logs.init or post-dev-setup",
      "auto_fix_available": true
    }</action>
    <action>Display: "🔵 INFO: Unified logging not set up"</action>
    <action>Display: "   Quick fix: /logs.init"</action>
    <action>Display: "   Or run: post-dev-setup (full setup with ports, etc.)"</action>
    <goto step="10.5"/>
  </check>

  <action>Check if .logs/all.ndjson has recent entries:</action>
  <action>Get last modified time of .logs/all.ndjson</action>
  <action>Calculate minutes since last log entry</action>

  <check if="minutes_since_last_log > 1440">
    <action>Display: "🔵 INFO: No recent log entries (last entry {hours} hours ago)"</action>
    <action>Display: "   Make sure your services are feeding logs to .logs/inbox/"</action>
  </check>

  <check if="minutes_since_last_log <= 1440">
    <action>Display: "✓ Logging system active (last entry {minutes} minutes ago)"</action>
  </check>

  <template-output>logging_checked</template-output>
</step>

<step n="10.5" goal="Check service-start launcher" if="skip_logs == false">
  <action>Display: ""</action>
  <action>Display: "🚀 Service Launcher"</action>
  <action>Display: ""</action>

  <action>Check if ./service-start exists:</action>

  <check if="./service-start missing">
    <action>Add to findings_info: {
      "category": "Service Launcher",
      "issue": "./service-start not found",
      "explanation": "service-start is a unified launcher that manages frontend/backend/docker services with automatic log aggregation. Makes development much easier.",
      "impact": "Manual service management, harder to coordinate logging",
      "fix": "Run post-dev-setup to create service-start",
      "auto_fix_available": true
    }</action>
    <action>Display: "🔵 INFO: service-start launcher not found"</action>
    <action>Display: "   Run post-dev-setup to create unified service launcher."</action>
    <goto step="11"/>
  </check>

  <action>Check if ./service-start is executable:</action>
  <check if="not executable">
    <action>Add to findings_warning: {
      "category": "Service Launcher",
      "issue": "./service-start not executable",
      "explanation": "The service-start file exists but isn't executable. It needs +x permissions.",
      "impact": "Can't run the launcher",
      "fix": "chmod +x ./service-start",
      "auto_fix_available": true
    }</action>
    <action>Display: "🟡 WARNING: service-start not executable"</action>
    <action>Display: "   Fix: chmod +x ./service-start"</action>
  </check>

  <check if="executable">
    <action>Display: "✓ service-start launcher present and executable"</action>
    <action>Display: "   Commands: ./service-start start|stop|restart|status"</action>
  </check>

  <template-output>service_start_checked</template-output>
</step>

<step n="11" goal="Check maintenance schedule" if="skip_maintenance == false">
  <action>Display: ""</action>
  <action>Display: "🔧 Maintenance Schedule"</action>
  <action>Display: ""</action>

  <action>Check for .bmad/last-maintenance timestamp file:</action>
  <action>If exists, read last maintenance date</action>
  <action>Calculate days since last maintenance</action>

  <check if="no maintenance timestamp found">
    <action>Display: "🔵 INFO: No maintenance record found"</action>
    <action>Display: "   Run 'maintenance' workflow to clean up and sync."</action>
    <action>Set days_since_maintenance = "never"</action>
  </check>

  <check if="days_since_maintenance > 30">
    <action>Add to findings_warning: {
      "category": "Maintenance",
      "issue": "Maintenance overdue by {days_since_maintenance - 30} days",
      "explanation": "Regular maintenance (every 30 days) keeps your project healthy: updates templates, cleans up dead code, syncs standards.",
      "impact": "Accumulating technical debt, outdated templates",
      "fix": "Run maintenance workflow",
      "auto_fix_available": true
    }</action>
    <action>Display: "🟡 WARNING: Maintenance overdue ({days_since_maintenance} days since last run)"</action>
    <action>Display: "   Run 'maintenance' to clean up and sync templates."</action>
  </check>

  <check if="days_since_maintenance <= 30 and days_since_maintenance != never">
    <action>Display: "✓ Maintenance up to date (last run {days_since_maintenance} days ago)"</action>
  </check>

  <template-output>maintenance_checked</template-output>
</step>

<step n="12" goal="Compile and present findings summary">
  <action>Display: ""</action>
  <action>Display: "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"</action>
  <action>Display: "📊 Health Check Summary"</action>
  <action>Display: ""</action>

  <action>Count findings by priority:</action>
  <action>critical_count = count(findings_critical)</action>
  <action>warning_count = count(findings_warning)</action>
  <action>info_count = count(findings_info)</action>
  <action>total_count = critical_count + warning_count + info_count</action>

  <check if="total_count == 0">
    <action>Display: "🎉 Excellent! Your project is in great health!"</action>
    <action>Display: "   No issues found. Keep up the good work!"</action>
    <goto step="13"/>
  </check>

  <action>Display: "Found {total_count} items needing attention:"</action>
  <action if="critical_count > 0">Display: "  🔴 {critical_count} CRITICAL issues"</action>
  <action if="warning_count > 0">Display: "  🟡 {warning_count} warnings"</action>
  <action if="info_count > 0">Display: "  🔵 {info_count} suggestions"</action>
  <action>Display: ""</action>

  <check if="critical_count > 0">
    <action>Display: "🔴 CRITICAL ISSUES (fix these first):"</action>
    <action>Display: ""</action>
    <action>For each finding in findings_critical:</action>
    <action>  Display: "• {finding.category}: {finding.issue}"</action>
    <action>  Display: "  Why it matters: {finding.explanation}"</action>
    <action>  Display: "  Impact: {finding.impact}"</action>
    <action>  Display: "  Fix: {finding.fix}"</action>
    <action if="finding.auto_fix_available">  Display: "  ⚡ Auto-fix available"</action>
    <action>  Display: ""</action>
  </check>

  <check if="warning_count > 0">
    <action>Display: "🟡 WARNINGS (address soon):"</action>
    <action>Display: ""</action>
    <action>For each finding in findings_warning:</action>
    <action>  Display: "• {finding.category}: {finding.issue}"</action>
    <action>  Display: "  Why it matters: {finding.explanation}"</action>
    <action>  Display: "  Impact: {finding.impact}"</action>
    <action>  Display: "  Fix: {finding.fix}"</action>
    <action if="finding.auto_fix_available">  Display: "  ⚡ Auto-fix available"</action>
    <action>  Display: ""</action>
  </check>

  <check if="info_count > 0">
    <action>Display: "🔵 SUGGESTIONS (nice to have):"</action>
    <action>Display: ""</action>
    <action>For each finding in findings_info:</action>
    <action>  Display: "• {finding.category}: {finding.issue}"</action>
    <action>  Display: "  {finding.explanation}"</action>
    <action>  Display: "  Fix: {finding.fix}"</action>
    <action if="finding.auto_fix_available">  Display: "  ⚡ Auto-fix available"</action>
    <action>  Display: ""</action>
  </check>

  <template-output>findings_summary</template-output>
</step>

<step n="13" goal="Offer auto-fixes and next steps">
  <action>Count auto-fixable issues:</action>
  <action>auto_fix_count = count of findings with auto_fix_available = true</action>

  <check if="auto_fix_count > 0">
    <action>Display: "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"</action>
    <action>Display: "⚡ Quick Fixes Available"</action>
    <action>Display: ""</action>
    <action>Display: "{auto_fix_count} issues can be fixed automatically:"</action>
    <action>Display: ""</action>

    <action>Build list of auto-fixable items:</action>
    <action>For each finding with auto_fix_available:</action>
    <action>  Display: "{index}. {finding.category}: {finding.issue}"</action>
    <action>  Add to auto_fixes_available list</action>

    <action>Display: ""</action>
    <ask>Would you like to apply these fixes? [y/n]</ask>

    <check if="user answered yes">
      <action>Display: "Applying fixes..."</action>
      <action>Display: ""</action>

      <action>For each auto-fix:</action>
      <check if="fix == Initialize Git repository">
        <action>Execute: git init</action>
        <action>Display: "✓ Git initialized"</action>
      </check>

      <check if="fix == Create GitHub repo">
        <action>Execute: gh repo create {project-name} --private --source=.</action>
        <action>Display: "✓ GitHub repository created"</action>
      </check>

      <check if="fix == Push to GitHub">
        <action>Execute: git push</action>
        <action>Display: "✓ Pushed to GitHub"</action>
      </check>

      <check if="fix == Pull latest changes">
        <action>Execute: git pull</action>
        <action>Display: "✓ Pulled latest changes"</action>
      </check>

      <check if="fix contains Set up Linear">
        <action>Display: "✓ Run quick-setup or brownfield-setup to configure Linear"</action>
      </check>

      <check if="fix contains Add via quick-setup or brownfield-setup">
        <action>Display: "✓ Run brownfield-setup to add missing components"</action>
      </check>

      <check if="fix contains Update packages">
        <action>Ask: "Run 'npm update' now? (this is safe for minor versions) [y/n]"</action>
        <check if="user answered yes">
          <action>Execute: npm update</action>
          <action>Display: "✓ Packages updated"</action>
        </check>
      </check>

      <check if="fix contains Create README.md">
        <action>Display: "✓ Run quick-setup or brownfield-setup to create README.md"</action>
      </check>

      <check if="fix == Run maintenance workflow">
        <action>Display: "✓ Run Atlas maintenance command to perform cleanup"</action>
      </check>

      <check if="fix contains chmod +x ./service-start">
        <action>Execute: chmod +x ./service-start</action>
        <action>Display: "✓ Made service-start executable"</action>
      </check>

      <check if="fix contains Run post-dev-setup to create service-start">
        <action>Display: "✓ Run post-dev-setup to create service-start launcher"</action>
      </check>

      <check if="fix contains Run /logs.init">
        <action>Display: "Initializing logging system..."</action>
        <action>Execute: /logs.init</action>
        <action>Display: "✓ Logging system initialized via /logs.init"</action>
      </check>

      <action>Display: ""</action>
      <action>Display: "✓ Auto-fixes applied! Re-run health-check to verify."</action>
    </check>

    <check if="user answered no">
      <action>Display: "Auto-fixes skipped. You can apply them manually using the commands above."</action>
    </check>
  </check>

  <action>Display: ""</action>
  <action>Display: "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"</action>
  <action>Display: "📋 Recommended Next Steps"</action>
  <action>Display: ""</action>

  <check if="critical_count > 0">
    <action>Display: "1. Address CRITICAL issues immediately - they block development"</action>
  </check>

  <check if="warning_count > 0">
    <action>Display: "2. Plan to fix WARNINGS this week - they accumulate technical debt"</action>
  </check>

  <check if="info_count > 0">
    <action>Display: "3. Consider INFO suggestions when you have time - they improve quality of life"</action>
  </check>

  <check if="total_count == 0">
    <action>Display: "Keep maintaining this health with regular checks!"</action>
    <action>Display: "Run health-check monthly or after major changes."</action>
  </check>

  <action>Display: ""</action>
  <action>Display: "💡 Tip: Run health-check regularly (monthly) to catch issues early."</action>

  <template-output>health_check_complete</template-output>
</step>

</workflow>
