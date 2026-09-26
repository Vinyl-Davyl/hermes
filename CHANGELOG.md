# Changelog

## 0.3.1

- `--id` picks one session from `hermes list`
- Handoff prints which session was taken
- README explains from → to (Cursor limit → Claude, Antigravity ↔ Cursor)

## 0.3.0

- CLI destinations start seeded (`.hermes/seed.md`), like a Catchup `--into`
- Cursor / editor destinations auto-copy PROMPT.md
- New landing page for demos (`hermes site`)
- Beginner PATH guide and demo script

## 0.2.0

- `hermes handoff` is the main command (`hermes handoff cursor`, `hermes handoff claude cursor`)
- `hermes resume` finds the latest pack — no `./handoff-*` glob (zsh nomatch)
- Shorthand: `hermes to`, `hermes from`
- More readers: Copilot CLI, ZCode, DeepSeek Harness
- Broader discovery paths and env overrides (`CLAUDE_CONFIG_DIR`, `CODEX_HOME`, …)
- `list -q` search and `list --here`
- `--copy` clipboard, `--open` start a destination CLI
- `make install` / `scripts/install.sh` for a local PATH install
- Fake session fixture for testing without other agents
- Landing page + docs updated for the simple command set

## 0.1.0

- `export`, `import`, `list`, `doctor`, `validate`, `landing`
- Git-first handoff packs (`MANIFEST.json` v1)
- Best-effort readers: Claude Code, Cursor, Codex, Cline, Kimi, Antigravity, OpenCode, Pi
- Local landing page (`web/` + `hermes landing`)
