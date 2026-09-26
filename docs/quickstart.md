# Quick start

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/Vinyl-Davyl/hermes/main/scripts/install.sh | sh
```

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

Many chats? `hermes list --agent cursor --here` then `--id` on the one you want.

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
