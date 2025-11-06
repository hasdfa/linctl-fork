# Post-Dev Setup Workflow Instructions

<critical>The workflow execution engine is governed by: {project-root}/bmad/core/tasks/workflow.xml</critical>
<critical>You MUST have already loaded and processed: {project-root}/bmad/bmm/workflows/post-dev-setup/workflow.yaml</critical>
<critical>Communicate in {communication_language} throughout the setup process</critical>
<critical>This workflow is LOW INTERACTIVITY - detect and execute with minimal prompts</critical>

<workflow>

<step n="0" goal="Parse flags and determine execution mode">
  <action>Check for flags in user input: --dry-run, --skip-ports, --skip-logs, --skip-cli, --skip-traefik</action>
  <action>Set dry_run_mode = true if --dry-run flag present</action>
  <action>Set skip flags for each component</action>

  <check if="dry_run_mode == true">
    <action>Display: "DRY RUN MODE - Scanning project, no changes will be made"</action>
  </check>

  <template-output>execution_mode</template-output>
</step>

<step n="1" goal="Detect project structure and services">
  <action>Display: "🔍 Detecting project structure and services..."</action>

  <action>Detect frontend framework:</action>
  <action>- Check for package.json → Look for "next", "react", "vue", "angular" in dependencies</action>
  <action>- Check for src/pages/, src/app/, src/components/ directories</action>
  <action>- Store: frontend_detected = true/false, frontend_type = "Next.js"|"React"|"Vue"|etc</action>

  <action>Detect backend framework:</action>
  <action>- Check for Python files (main.py, app.py) → Look for FastAPI, Django, Flask imports</action>
  <action>- Check for Node.js backend (server.js, index.js in separate directory)</action>
  <action>- Check for Go files (main.go, go.mod)</action>
  <action>- Store: backend_detected = true/false, backend_type = "FastAPI"|"Django"|"Express"|etc</action>

  <action>Detect database:</action>
  <action>- Check for docker-compose.yml with postgres, mysql, mongodb services</action>
  <action>- Check for database connection configs</action>
  <action>- Store: database_detected = true/false, database_type = "postgres"|"mysql"|etc</action>

  <check if="frontend_detected == false and backend_detected == false">
    <action>Display: "⚠️  No frontend or backend detected yet."</action>
    <action>Display: "This command should be run after you have minimally working code."</action>
    <ask>Continue anyway? [y/n]</ask>
    <check if="user answered no">
      <action>Exit workflow with message: "Run this command when you have working frontend/backend code."</action>
    </check>
  </check>

  <action>Display detected services:</action>
  <action if="frontend_detected == true">Display: "  ✓ Frontend: {frontend_type}"</action>
  <action if="backend_detected == true">Display: "  ✓ Backend: {backend_type}"</action>
  <action if="database_detected == true">Display: "  ✓ Database: {database_type}"</action>

  <template-output>detected_services</template-output>
</step>

<step n="2" goal="Configure dynamic port assignment" if="skip_ports == false">
  <action>Display: "📡 Configuring dynamic port assignment..."</action>

  <action>Check if portbase utility is available:</action>
  <action>Execute: which portbase or check ~/bin/portbase exists</action>

  <check if="portbase not found">
    <action>Display: "⚠️  portbase utility not found. Port assignment requires system setup."</action>
    <action>Display: "See: docs/canonical/how-to/port-assignment-setup.md"</action>
    <action>Display: "Or install with: cp scripts/dev-proxy/portbase.py ~/bin/portbase && chmod +x ~/bin/portbase"</action>
    <ask>Skip port assignment for now? [y/n]</ask>
    <check if="user answered yes">
      <goto step="3"/>
    </check>
    <check if="user answered no">
      <action>Exit workflow with message: "Install portbase and re-run this command."</action>
    </check>
  </check>

  <action>Generate deterministic base port:</action>
  <action>Execute: portbase</action>
  <action>Store: BASE_PORT = output</action>

  <action>Calculate service ports:</action>
  <action>WEB_PORT = BASE_PORT + 0 (frontend)</action>
  <action>API_PORT = BASE_PORT + 1 (backend)</action>
  <action>DB_PORT = BASE_PORT + 2 (database)</action>
  <action>DOCS_PORT = BASE_PORT + 3 (documentation)</action>
  <action>WS_PORT = BASE_PORT + 4 (websocket)</action>

  <action>Get project name from directory:</action>
  <action>PROJECT_NAME = basename of current directory (lowercase, hyphens instead of spaces)</action>
  <action>PROJECT_DOMAIN = "{PROJECT_NAME}.localtest.me"</action>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would create .envrc with ports:"</action>
    <action>Display: "   BASE_PORT={BASE_PORT}"</action>
    <action if="frontend_detected">Display: "   WEB_PORT={WEB_PORT} (Frontend: {frontend_type})"</action>
    <action if="backend_detected">Display: "   API_PORT={API_PORT} (Backend: {backend_type})"</action>
    <action if="database_detected">Display: "   DB_PORT={DB_PORT} (Database: {database_type})"</action>
    <action>Display: "   PROJECT_DOMAIN={PROJECT_DOMAIN}"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Create .envrc file with port assignments:</action>
    <action>Write to .envrc:
```bash
# Deterministic port assignment
# Generated by Atlas post-dev-setup

export BASE_PORT=$(portbase)
export WEB_PORT=$BASE_PORT
export API_PORT=$((BASE_PORT+1))
export DB_PORT=$((BASE_PORT+2))
export DOCS_PORT=$((BASE_PORT+3))
export WS_PORT=$((BASE_PORT+4))

# Custom domain for Traefik proxy
export PROJECT_DOMAIN="{PROJECT_DOMAIN}"
```
    </action>

    <action>Activate direnv:</action>
    <action>Execute: direnv allow</action>

    <action>Verify ports are loaded:</action>
    <action>Execute: bash -c "source .envrc && echo $WEB_PORT"</action>

    <action>Display: "✓ Port assignment configured"</action>
    <action if="frontend_detected">Display: "   Frontend: http://localhost:{WEB_PORT}"</action>
    <action if="backend_detected">Display: "   Backend: http://localhost:{API_PORT}"</action>
    <action if="database_detected">Display: "   Database: localhost:{DB_PORT}"</action>

    <action>Update .gitignore to exclude .envrc:</action>
    <check if=".gitignore exists">
      <action if=".envrc not in .gitignore">Append ".envrc" to .gitignore</action>
    </check>
  </check>

  <template-output>ports_configured</template-output>
</step>

<step n="2.5" goal="Update application code to use ports" if="skip_ports == false and dry_run_mode == false">
  <action>Display: "Updating application code to use assigned ports..."</action>

  <check if="frontend_detected and frontend_type == Next.js">
    <action>Check if package.json exists</action>
    <action>Update package.json dev script to use PORT environment variable:</action>
    <action>Pattern: "dev": "PORT=${PORT:-${WEB_PORT:-3000}} next dev"</action>
    <action>Display: "✓ Updated Next.js to use $WEB_PORT"</action>
  </check>

  <check if="backend_detected and backend_type == FastAPI">
    <action>Find main Python file (main.py, app.py, or similar)</action>
    <action>Display: "💡 Add port configuration to your FastAPI app:"</action>
    <action>Display: '```python
import os
port = int(os.getenv("PORT", os.getenv("API_PORT", "8000")))
uvicorn.run(app, host="0.0.0.0", port=port)
```'</action>
  </check>

  <check if="backend_detected and backend_type == Express">
    <action>Display: "💡 Add port configuration to your Express app:"</action>
    <action>Display: '```javascript
const PORT = process.env.PORT || process.env.API_PORT || 3000;
app.listen(PORT, "0.0.0.0", () => {
  console.log(`Server running on port ${PORT}`);
});
```'</action>
  </check>

  <action>Display: "💡 Note: Applications should check PORT first (production), then WEB_PORT/API_PORT (local)"</action>

  <template-output>app_code_updated</template-output>
</step>

<step n="2.6" goal="Configure Traefik proxy (optional)" if="skip_ports == false and skip_traefik == false">
  <ask>Set up Traefik proxy for friendly URLs ({PROJECT_DOMAIN})? [y/n]</ask>

  <check if="user answered no">
    <action>Display: "Skipping Traefik setup. Use localhost:{WEB_PORT} URLs."</action>
    <goto step="3"/>
  </check>

  <action>Check if Traefik is running:</action>
  <action>Execute: docker compose -f ~/dev-proxy/docker-compose.yml ps | grep traefik</action>

  <check if="Traefik not running">
    <action>Display: "Traefik proxy not running."</action>
    <ask>Start Traefik now? [y/n]</ask>
    <check if="user answered yes">
      <action>Execute: cd ~/dev-proxy && docker compose up -d</action>
      <action>Wait 2 seconds for Traefik to start</action>
      <action>Display: "✓ Traefik proxy started"</action>
    </check>
    <check if="user answered no">
      <action>Display: "Skipping Traefik. Start manually with: cd ~/dev-proxy && docker compose up -d"</action>
      <goto step="3"/>
    </check>
  </check>

  <action>Register project with Traefik:</action>
  <action if="frontend_detected">Execute: ~/bin/add-proj {PROJECT_NAME} http://127.0.0.1:{WEB_PORT}</action>
  <action if="backend_detected and frontend_detected == false">Execute: ~/bin/add-proj {PROJECT_NAME} http://127.0.0.1:{API_PORT}</action>

  <action>Restart Traefik to load new configuration:</action>
  <action>Execute: cd ~/dev-proxy && docker compose restart traefik</action>

  <action>Display: "✓ Traefik proxy configured"</action>
  <action>Display: "   Access via: http://{PROJECT_DOMAIN}"</action>
  <action>Display: "   Dashboard: http://localhost:8080"</action>

  <template-output>traefik_configured</template-output>
</step>

<step n="3" goal="Initialize unified logging system" if="skip_logs == false">
  <action>Display: "📝 Initializing unified logging system..."</action>

  <action>Check if .logs/ directory already exists:</action>

  <check if=".logs/ exists and .logs/all.ndjson exists">
    <action>Display: "✓ Logging already initialized, skipping"</action>
    <goto step="3.5"/>
  </check>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would run /logs.init to create:"</action>
    <action>Display: "   .logs/all.ndjson - Main log file"</action>
    <action>Display: "   .logs/inbox/ - Log sources directory"</action>
    <action>Display: "   .logs/redaction.json - Optional redaction config"</action>
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
      <action>Display: "⚠️  /logs.init failed, creating manually..."</action>
      <action>Create logging directory structure:</action>
      <action>Execute: mkdir -p .logs/inbox</action>

      <action>Create empty main log file:</action>
      <action>Execute: touch .logs/all.ndjson</action>

      <action>Create optional redaction config:</action>
      <action>Write to .logs/redaction.json:
```json
{
  "enabled": false,
  "patterns": [
    {
      "pattern": "AKIA[0-9A-Z]{16}",
      "replacement": "[REDACTED_AWS_KEY]"
    },
    {
      "pattern": "(?i)secret\\s*[:=]\\s*[^\\s,;]+",
      "replacement": "secret:[REDACTED]"
    },
    {
      "pattern": "(?i)(api[_-]?key|token|password|passwd|pwd)\\s*[:=]\\s*[^\\s,;]+",
      "replacement": "\\1:[REDACTED]"
    },
    {
      "pattern": "Bearer\\s+[A-Za-z0-9\\-._~+/]+=*",
      "replacement": "Bearer [REDACTED]"
    }
  ]
}
```
    </action>

      <action>Update .gitignore to include log directory:</action>
      <check if=".gitignore exists">
        <action if=".logs/ not in .gitignore">Append ".logs/" to .gitignore</action>
        <action if=".logs/redaction.json not in .gitignore">Append "!.logs/redaction.json" to .gitignore (keep config in git)</action>
      </check>

      <action>Display: "✓ Logging system initialized (manual fallback)"</action>
      <action>Display: "   Main log: .logs/all.ndjson"</action>
      <action>Display: "   Inbox: .logs/inbox/"</action>
    </check>

    <action>Display: "💡 To feed logs from your services:"</action>
    <action if="frontend_detected">Display: '   Frontend: {frontend command} 2>&1 | tee -a .logs/inbox/frontend.log'</action>
    <action if="backend_detected">Display: '   Backend: {backend command} 2>&1 | tee -a .logs/inbox/backend.log'</action>
    <action>Display: '   Or use: ./service-start start (recommended)'</action>
    <action>Display: '   View logs: tail -f .logs/all.ndjson'</action>
  </check>

  <template-output>logging_initialized</template-output>
</step>

<step n="3.5" goal="Create service-start script" if="skip_logs == false">
  <action>Display: "🚀 Creating service-start launcher..."</action>

  <action>Check if ./service-start already exists:</action>

  <check if="./service-start exists">
    <action>Display: "✓ service-start already exists, skipping"</action>
    <goto step="4"/>
  </check>

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

<step n="4" goal="Create CLI menu script" if="skip_cli == false">
  <action>Display: "🎯 Creating CLI menu script..."</action>

  <action>Check if ./cli-menu already exists:</action>

  <check if="./cli-menu exists">
    <action>Display: "✓ CLI menu already exists, skipping"</action>
    <goto step="5"/>
  </check>

  <check if="dry_run_mode == true">
    <action>Display: "[DRY RUN] Would create ./cli-menu script"</action>
    <action>Display: "   - Interactive menu for backend admin tasks"</action>
    <action>Display: "   - Logs all operations to .logs/inbox/cli-menu.log"</action>
    <action>Display: "   - Customizable per project"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Create CLI menu script template:</action>
    <action>Write to ./cli-menu:
```bash
#!/usr/bin/env bash
# CLI Menu for {PROJECT_NAME}
# Generated by Atlas post-dev-setup
# Customize this script with your backend admin tasks

set -euo pipefail

# Logging to unified log system
LOG_FILE=".logs/inbox/cli-menu.log"
mkdir -p "$(dirname "$LOG_FILE")"

log() {
  echo "[$(date -u +"%Y-%m-%dT%H:%M:%SZ")] $*" | tee -a "$LOG_FILE"
}

# Load environment variables
if [ -f .envrc ]; then
  # shellcheck disable=SC1091
  source .envrc
fi

# Menu function
show_menu() {
  echo "╔══════════════════════════════════════╗"
  echo "║  {PROJECT_NAME} - CLI Menu           ║"
  echo "╠══════════════════════════════════════╣"
  echo "║                                      ║"
  echo "║  1) Start development servers        ║"
  echo "║  2) Run database migrations          ║"
  echo "║  3) Seed database with test data     ║"
  echo "║  4) Run tests                        ║"
  echo "║  5) Build production                 ║"
  echo "║  6) View logs                        ║"
  echo "║  7) Database backup                  ║"
  echo "║  8) Custom task (placeholder)        ║"
  echo "║  0) Exit                             ║"
  echo "║                                      ║"
  echo "╚══════════════════════════════════════╝"
  echo
  read -r -p "Select option [0-8]: " choice
  echo
  return 0
}

# Task implementations (customize these)
task_start_dev() {
  log "Starting development servers..."
  # Add your start command here
  echo "TODO: Add your development server start command"
  # Example: npm run dev & uvicorn main:app --reload
}

task_migrate_db() {
  log "Running database migrations..."
  # Add your migration command here
  echo "TODO: Add your migration command"
  # Example: alembic upgrade head
}

task_seed_db() {
  log "Seeding database with test data..."
  # Add your seed command here
  echo "TODO: Add your seed command"
}

task_run_tests() {
  log "Running tests..."
  # Add your test command here
  echo "TODO: Add your test command"
  # Example: pytest
}

task_build_prod() {
  log "Building for production..."
  # Add your build command here
  echo "TODO: Add your build command"
  # Example: npm run build
}

task_view_logs() {
  log "Viewing unified logs..."
  tail -f .logs/all.ndjson
}

task_backup_db() {
  log "Creating database backup..."
  # Add your backup command here
  echo "TODO: Add your backup command"
}

task_custom() {
  log "Running custom task..."
  # Add your custom task here
  echo "TODO: Implement custom task"
}

# Main loop
main() {
  log "CLI menu started"

  while true; do
    show_menu

    case $choice in
      1) task_start_dev ;;
      2) task_migrate_db ;;
      3) task_seed_db ;;
      4) task_run_tests ;;
      5) task_build_prod ;;
      6) task_view_logs ;;
      7) task_backup_db ;;
      8) task_custom ;;
      0) log "CLI menu exited"; echo "Goodbye!"; exit 0 ;;
      *) echo "Invalid option. Please try again." ;;
    esac

    echo
    read -r -p "Press Enter to continue..."
    clear
  done
}

# Run main
clear
main
```
    </action>

    <action>Make cli-menu executable:</action>
    <action>Execute: chmod +x ./cli-menu</action>

    <action>Update .gitignore to exclude cli-menu (project-specific):</action>
    <check if=".gitignore exists">
      <action if="cli-menu not in .gitignore">Append "cli-menu" to .gitignore</action>
    </check>

    <action>Display: "✓ CLI menu script created"</action>
    <action>Display: "   Run with: ./cli-menu"</action>
    <action>Display: "   Customize tasks in the script for your project"</action>
  </check>

  <template-output>cli_menu_created</template-output>
</step>

<step n="5" goal="Final summary and next steps">
  <action>Generate post-dev-setup completion summary</action>

  <action>Show what was configured:</action>
  <action if="skip_ports == false">Display: "✓ Port assignment configured"</action>
  <action if="skip_ports == false and skip_traefik == false">Display: "✓ Traefik proxy registered"</action>
  <action if="skip_logs == false">Display: "✓ Unified logging initialized"</action>
  <action if="skip_logs == false">Display: "✓ service-start launcher created"</action>
  <action if="skip_cli == false">Display: "✓ CLI menu script created"</action>

  <check if="dry_run_mode == true">
    <action>Display: "DRY RUN COMPLETE - No changes were made"</action>
    <action>Display: "Run without --dry-run to execute post-dev-setup"</action>
  </check>

  <check if="dry_run_mode == false">
    <action>Display: "✓ Post-dev setup complete! Your project is operationally configured."</action>

    <action>Display configuration summary:</action>
    <action if="skip_ports == false">
      <action>Display: "Ports (via direnv):"</action>
      <action if="frontend_detected">Display: "  - Frontend: http://localhost:{WEB_PORT}"</action>
      <action if="backend_detected">Display: "  - Backend: http://localhost:{API_PORT}"</action>
      <action if="database_detected">Display: "  - Database: localhost:{DB_PORT}"</action>
      <action if="skip_traefik == false">Display: "  - Proxy: http://{PROJECT_DOMAIN}"</action>
    </action>

    <action if="skip_logs == false">
      <action>Display: "Logging:"</action>
      <action>Display: "  - Main log: .logs/all.ndjson"</action>
      <action>Display: "  - Service launcher: ./service-start start"</action>
      <action>Display: "  - Manual feed: command 2>&1 | tee -a .logs/inbox/service.log"</action>
    </action>

    <action if="skip_cli == false">
      <action>Display: "CLI Menu:"</action>
      <action>Display: "  - Run: ./cli-menu"</action>
      <action>Display: "  - Customize tasks in the script"</action>
    </action>

    <action>Display next steps:</action>
    <action if="skip_logs == false">Display: "1. Start services: ./service-start start (auto-manages aggregator + services)"</action>
    <action if="skip_logs == true">Display: "1. Start your development servers (they'll use assigned ports)"</action>
    <action if="skip_ports == false">Display: "2. Verify ports: echo $WEB_PORT $API_PORT"</action>
    <action if="skip_logs == false">Display: "3. Check status: ./service-start status"</action>
    <action if="skip_cli == false">Display: "4. Customize ./cli-menu with your admin tasks"</action>
    <action if="skip_traefik == false">Display: "5. Access app via: http://{PROJECT_DOMAIN}"</action>
  </check>

  <template-output>post_dev_setup_complete</template-output>
</step>

</workflow>
