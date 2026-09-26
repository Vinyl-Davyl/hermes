<p align="center">
  <img src="web/logo.png" width="72" height="72" alt="Hermes">
</p>

# Hermes

**Local-first CLI for coding-agent context handoff.**

Named for the messenger of the gods. Pack the work. Open the next chat. Keep going.

Transfer and resume across **Claude Code, Codex, Cursor, Cline, Kimi, Antigravity, OpenCode, Pi Agent, Copilot CLI, ZCode, and DeepSeek Harness**.

No MCP. No plugin. No cloud. One Go binary.

## Install

```bash
curl -fsSL https://tryhermes.pages.dev/install | sh
```

Site: [tryhermes.pages.dev](https://tryhermes.pages.dev)

Then open a new terminal and run `hermes doctor`.

If the shell says `command not found`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

## The only rule

```text
hermes handoff <WHERE YOU ARE GOING>
hermes handoff <WHERE YOU ARE NOW> <WHERE YOU ARE GOING>
```

**One name = destination.**  
**Two names = from, then to.**

You are in **Cursor**. The limit expired. You want **Claude**:

```bash
hermes handoff claude
hermes handoff cursor claude
```

`hermes handoff cursor` means _go to Cursor_, not _leave Cursor_.

| I am here          | I want to continue here | Command                             |
| ------------------ | ----------------------- | ----------------------------------- |
| Cursor (limit hit) | Claude                  | `hermes handoff claude`             |
| Claude (limit hit) | Cursor                  | `hermes handoff cursor`             |
| Cursor             | Antigravity             | `hermes handoff cursor antigravity` |
| Antigravity        | Cursor                  | `hermes handoff antigravity cursor` |

Hermes only packs a chat that belongs to **this folder**. Another project (another window, another repo) is ignored. If nothing here matches, the pack is git-only. Pass `--id` to pick a specific chat.

After the command:

- **Claude / Codex / OpenCode** (if that CLI is installed): Hermes starts it with the pack.
- **Cursor / Antigravity**: prompt is on the clipboard. New chat → paste.

## Many chats?

```bash
hermes list --agent cursor --here
hermes handoff cursor antigravity --id e1884976
```

## Commands

| Command                      | What it does                  |
| ---------------------------- | ----------------------------- |
| `hermes handoff [from] [to]` | Pack + send to the next agent |
| `hermes resume [to]`         | Reuse the latest pack         |
| `hermes list`                | Sessions for this machine     |
| `hermes list --here`         | Sessions for this folder      |
| `hermes doctor`              | What Hermes can see           |
| `hermes site`                | Local landing page            |

```text
-m, --message     what you were doing
    --from        source agent
    --to          destination
    --id          session id from hermes list
    --no-open     do not start the next CLI
```

## Docs

- [Quick start](docs/quickstart.md)
- [Supported tools](docs/supported-tools.md)
- [FAQ](docs/faq.md)

MIT
