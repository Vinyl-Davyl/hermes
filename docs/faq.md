# FAQ

## I am in Cursor and the limit expired. Is the command `handoff cursor` or `handoff claude`?

`hermes handoff claude`. One name is the place you are **going**. `hermes handoff cursor` would try to continue **in Cursor**. Two names make it unmistakable: `hermes handoff cursor claude`.

## Cursor → Antigravity? Antigravity → Cursor?

```bash
hermes handoff cursor antigravity
hermes handoff antigravity cursor
```

Then paste into a new chat in the destination (or `@PROMPT.md`).

## What if I have many Cursor sessions?

```bash
hermes list --agent cursor --here
hermes handoff cursor claude --id <id-from-the-list>
```

Without `--id`, Hermes picks the newest session that looks like this folder.

---

## Does Hermes send my code anywhere?

No. Local disk only. The landing page is a static file served on `127.0.0.1`.

## Why do I still see PROMPT.md? Does Catchup use that?

Catchup never asks you to attach `PROMPT.md`. `catchup fork claude --into codex` **starts Codex** and feeds it the transcript. Hermes now does the same for CLI agents (`hermes handoff claude`). Cursor cannot be started that way as an IDE chat, so Hermes puts the same text on your clipboard and in the pack. You paste once, or `@PROMPT.md`. Catchup’s honest version of that sentence is “save the transcript and paste it.”

## Do I have to install Hermes globally?

```bash
curl -fsSL https://tryhermes.pages.dev/install | sh
```

That puts `hermes` on your PATH. From a checkout, `./bin/hermes` after `go build` is the same program.

## Why did zsh say `no matches found: ./handoff-*`?

zsh expands `*` before the program starts. If no `handoff-*` folder exists, zsh never launches Hermes. Create a pack with `hermes handoff`, then reuse it with `hermes resume`. Do not type the star.

## Do I have to run Hermes inside a git repo?

No. A repo gives you `git/diff.patch` and changed-file names. Outside git, Hermes still writes a pack for this folder and warns you.

## I have not installed Claude / Codex / Cline. Can I still test?

Yes. `hermes handoff --from none` packs git. `examples/fake-session/claude.jsonl` is a stand-in transcript. Real agents are only needed when you want a live session excerpt.

## Can I restore the same Cursor composer?

No. Cursor has no public session import. Claiming otherwise would be dishonest. Hermes continues the *work*.

## Why is `hermes list` empty for an agent I use?

Discovery paths differ by OS and version. Run `hermes doctor`. You can still `handoff --from none` or pass `--session /path/to/file.jsonl`.

## Will this replace Claude `/export`?

Use both. `/export` is a transcript dump. Hermes adds git state, a structured prompt, and a resume recipe for the next window.

## Does `--open` import a session into the next agent?

No. It starts that CLI in this repo if the binary is on PATH. You still paste or `@` `PROMPT.md`. Cursor cannot be opened this way.
