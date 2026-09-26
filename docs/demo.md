# Film a demo (Catchup-style, Hermes commands)

Record from a **real git repo** (not this Hermes folder if it has no `.git`). A dirty tree looks better on camera.

## 30-second A-roll

1. Terminal left: pretend Claude just died (or use a real `claude` session).
2. Terminal right:

```bash
hermes doctor
hermes list
hermes handoff claude -m "fix the flaky auth test"
```

If `claude` is installed, it opens already pointed at `.hermes/seed.md`.  
If not, say on camera: “Claude is not on this machine, so I open Cursor and paste.”

## Cursor on camera

```bash
hermes handoff cursor -m "fix the flaky auth test"
```

PROMPT.md is on the clipboard. New Cursor chat → paste, or type `@PROMPT.md`.

## Other agents on camera

```bash
hermes handoff cursor claude    # Cursor → Claude
hermes handoff claude cursor    # Claude → Cursor
hermes handoff --from none -m "git only demo"
hermes list --agent cursor
hermes list -q auth
```

`--no-open` if you only want the pack, not a launched CLI.

## Do not type this

```bash
hermes import ./handoff-*
```

zsh will eat the star. Use `hermes resume`.
