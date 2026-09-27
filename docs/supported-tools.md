# Supported tools

Hermes **reads** sessions when it can. It **never writes** into an editor database.

| ID | Product | Typical location | Notes |
|----|---------|------------------|--------|
| `claude` | Claude Code | `$CLAUDE_CONFIG_DIR/projects` or `~/.claude/projects` | Best supported transcript |
| `cursor` | Cursor | `state.vscdb` under app support / `$CURSOR_CONFIG_DIR` | Read-only. Resume = new chat + `@PROMPT.md` |
| `codex` | Codex CLI | `$CODEX_HOME/sessions` or `~/.codex/sessions` | |
| `opencode` | OpenCode | `$XDG_DATA_HOME/opencode`, `~/.local/share/opencode` | JSONL or `opencode.db` |
| `cline` | Cline / Roo | `$CLINE_DIR`, `~/.cline`, VS Code / Cursor `globalStorage` | |
| `kimi` | Kimi | `$KIMI_CODE_HOME`, `~/.kimi-code`, `~/.kimi` | |
| `antigravity` | Antigravity | `~/.gemini/antigravity-cli`, app support | alias: `agy` |
| `pi` | Pi Agent | `$PI_CODING_AGENT_DIR`, `~/.pi/agent` | |
| `copilot` | Copilot CLI | `$COPILOT_HOME/session-state` or `~/.copilot` | |
| `zcode` | ZCode | `$ZCODE_HOME` or `~/.zcode/cli/db` | SQLite, read-only |
| `deepseek` | DeepSeek Harness | `$DSH_HOME/sessions` or `~/.dsh` | alias: `dsh` |

Aliases: `claude-code` → `claude`, `codex-cli` → `codex`, `open-code` → `opencode`, `agy` → `antigravity`, `dsh` → `deepseek`.

There is no `vscode` agent. VS Code chat and the Claude desktop/web app use the same paste path as Cursor: `--no-open --copy`.

`--from none` skips all session readers (git-only pack).

`--from auto` prefers a Claude session for this repo, then the newest session of any kind.

`hermes doctor` shows which of these Hermes can see on this machine.

If a path is wrong on your OS, handoff still succeeds with git. Open an issue with the real directory.
