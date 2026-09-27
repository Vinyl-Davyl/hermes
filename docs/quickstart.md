# Quick start

## Install

```bash
curl -fsSL https://tryhermes.pages.dev/install | sh
```

Site: [tryhermes.pages.dev](https://tryhermes.pages.dev)

Then `hermes doctor`.

## The rule

```text
hermes handoff <GOING TO>
hermes handoff <FROM> <TO>
```

Cursor died, continue in Claude:

```bash
hermes handoff claude
hermes handoff cursor claude
```

Antigravity → Cursor:

```bash
hermes handoff antigravity cursor
```

## Pick a chat

```bash
hermes list --here
hermes list --agent cursor --here -q keyword
hermes handoff cursor claude --id e188
```

A prefix of the id is enough. `hermes list` without `--here` is every session on this machine.

## Resume

```bash
hermes resume
hermes resume cursor
hermes resume --print
```

`--print` prints `PROMPT.md`. Never run `./handoff-*`.

## GUI, not the terminal CLI

```bash
hermes handoff claude cursor --no-open --copy -m "what you were doing"
```

Then a new chat in Cursor, the Claude app, or VS Code → paste.

## Build

```bash
go test ./...
go build -o bin/hermes ./cmd/hermes
./bin/hermes version
```

Optional: `make install` then `export PATH="$HOME/.local/bin:$PATH"`.

## First pack

```bash
./bin/hermes doctor
./bin/hermes list
./bin/hermes handoff --from none -m "Describe the work" --no-open
./bin/hermes resume
```

Do not run `./handoff-*`. zsh will fail. Use `resume`.
