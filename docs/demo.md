# Film the 30-second demo

One project. Two windows. One command.

## Setup

- Use a **real app repo** you already have (not an empty folder).
- Window A: **Cursor**, that repo open, a chat mid-task (“fix the flaky auth test”).
- Window B: a normal terminal, `cd` to the **same** repo.
- Window C (after the command): **Antigravity**, same repo.

Do not open Credo, Hermes, and your portfolio at once. Hermes now only packs this folder.

## What you type (Window B)

```bash
hermes doctor
hermes handoff cursor antigravity -m "fix the flaky auth test"
```

You should see `from cursor` and `to antigravity`, and `copied PROMPT.md`.

## What you do in Antigravity

1. Open the same folder.
2. New chat (do not resume an old one).
3. Paste (clipboard already has the prompt), or attach `@PROMPT.md`.
4. Let it pick up the task. Stop recording.

## If you would rather demo Claude

Same setup, but Window C is a terminal with Claude Code:

```bash
hermes handoff cursor claude -m "fix the flaky auth test"
```

If `claude` is on PATH, it starts with the pack. You do not paste.

## After you have the file

1. Save as `web/handoff.mp4` (keep it under ~15 MB if you can).
2. Uncomment the `<video>` tag in `web/index.html`.
3. Add `handoff.mp4` to the `//go:embed` line in `web/embed.go`.
