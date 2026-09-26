# Hermes

**Local-first CLI for coding-agent context handoff.**

Named for the messenger of the gods. Pack the work. Open the next chat. Keep going.

Transfer and resume across **Claude Code, Codex, Cursor, Cline, Kimi, Antigravity, OpenCode, Pi Agent, Copilot CLI, ZCode, and DeepSeek Harness**.

No MCP. No plugin. No cloud. One Go binary.

New to Go? [docs/beginner.md](docs/beginner.md). Film a demo? [docs/demo.md](docs/demo.md).

---

## The only rule (read this twice)

```text
hermes handoff <WHERE YOU ARE GOING>
hermes handoff <WHERE YOU ARE NOW> <WHERE YOU ARE GOING>
```

**One name = destination.**  
**Two names = from, then to.**

You are in a **Cursor** chat. The limit expired. You want **Claude**:

```bash
hermes handoff claude
# same thing, spelled out:
hermes handoff cursor claude
hermes handoff --from cursor --to claude
```

`hermes handoff cursor` is the **wrong** command for that. One name means _go to Cursor_, not _leave Cursor_.

| I am here                     | I want to continue here | Command                                                   |
| ----------------------------- | ----------------------- | --------------------------------------------------------- |
| Cursor (limit hit)            | Claude                  | `hermes handoff claude` or `hermes handoff cursor claude` |
| Claude (limit hit)            | Cursor                  | `hermes handoff cursor` or `hermes handoff claude cursor` |
| Cursor                        | Antigravity             | `hermes handoff cursor antigravity`                       |
| Antigravity                   | Cursor                  | `hermes handoff antigravity cursor`                       |
| Cursor                        | Cline                   | `hermes handoff cursor cline`                             |
| any agent                     | Cursor                  | `hermes handoff cursor`                                   |
| newest session Hermes can see | Claude                  | `hermes handoff claude`                                   |

After the command:

- **Claude / Codex / OpenCode / Pi / Copilot / Cline** (if that CLI is installed): Hermes starts it and points it at `.hermes/seed.md`.
- **Cursor / Antigravity / ZCode / Kimi**: the prompt is on your clipboard. Open a **new** chat and paste, or attach `@PROMPT.md`.

Hermes never writes into Cursor’s database. A new chat is required. That is the product.

---

## What if I have many Cursor chats?

`hermes handoff cursor claude` takes the **newest Cursor session that looks like this folder**. It is not always “the chat I have open right now.”

1. See the list:

```bash
hermes list --agent cursor --here
```

2. Copy the id (the long hash). A prefix is enough.

3. Hand off that one:

```bash
hermes handoff cursor claude --id e1884976
```

Same pattern for Antigravity or Claude:

```bash
hermes list --agent antigravity --here
hermes handoff antigravity cursor --id <id>
```

---

## The flow, step by step

```bash
export PATH="$HOME/.local/bin:$PATH"   # only if this tab is old
hermes doctor                          # what can Hermes see?
hermes list                            # sessions, newest first

hermes handoff cursor                  # go TO Cursor (from the newest session)
hermes handoff claude                  # go TO Claude
hermes handoff cursor claude           # FROM Cursor TO Claude
hermes handoff claude cursor           # FROM Claude TO Cursor
hermes handoff --from none -m "demo" --no-open
```

`--from none` skips every agent chat and packs **git only**. Use it for a dry run.  
`--no-open` writes the pack and does not start Claude/Codex.

Run these from the **project folder** you were editing, not from a random directory.

---

## Install (not npm)

Go builds one file. You copy that file. You do not publish to npm.

```bash
go test ./...
go build -o bin/hermes ./cmd/hermes
./bin/hermes doctor          # always works from this folder

make install                 # ~/.local/bin/hermes
export PATH="$HOME/.local/bin:$PATH"
hermes doctor
```

Add the `export` line to `~/.zshrc` so new terminals find `hermes`.

---

## Why `./handoff-*` failed

zsh expands `*` before Hermes starts. If no pack exists, zsh says `no matches found`. Create a pack with `hermes handoff`, reuse it with `hermes resume`. Never type the star.

---

## Commands

| Command                             | What it does                                 |
| ----------------------------------- | -------------------------------------------- |
| `hermes handoff [from] [to]`        | Pack + send to the next agent                |
| `hermes resume [to]`                | Reuse the latest pack                        |
| `hermes to claude`                  | Same as `handoff claude`                     |
| `hermes from cursor`                | Pack from Cursor, default destination Cursor |
| `hermes list`                       | Sessions Hermes can see                      |
| `hermes list --agent cursor --here` | This project’s Cursor chats                  |
| `hermes list -q auth`               | Search                                       |
| `hermes doctor`                     | Git + which agents exist                     |
| `hermes site`                       | Landing page at http://127.0.0.1:8787        |

### handoff flags

```text
-m, --message         what you were doing
    --from cursor     source agent (or auto|none|…)
    --to claude       destination agent
    --id e1884976     pick one session from hermes list
    --session PATH    explicit transcript file
    --no-open         do not start the next CLI
    --copy            copy PROMPT.md (Cursor already does this)
```

---

## Pack layout

```text
handoff-<timestamp>/
  MANIFEST.json
  PROMPT.md
  SUMMARY.md
  git/diff.patch
  sessions/<agent>.md
```

You can delete an old `handoff-*` folder. It is a note, not source code.

---

## Docs

- [Quick start](docs/quickstart.md)
- [Supported tools](docs/supported-tools.md)
- [FAQ](docs/faq.md)
- [Go / PATH for beginners](docs/beginner.md)
- [Demo script](docs/demo.md)
- Landing: `hermes site`

---

## What Hermes will not do

- Write into Cursor `state.vscdb` or Claude project files
- Merge two agents into one history
- Upload transcripts
- Require npm or an MCP server

---

MIT
