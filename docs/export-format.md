# Handoff pack format (v1)

```text
handoff-<timestamp>/
├── MANIFEST.json
├── PROMPT.md              # generated continuation prompt
├── SUMMARY.md
├── DECISIONS.md           # optional
├── git/
│   ├── branch.txt
│   ├── diff.patch
│   ├── log.txt
│   └── changed-files.txt
├── files/                 # optional, --include-files
└── sessions/
    └── <agent>.md         # best-effort excerpt
```

## MANIFEST.json

| Field | Meaning |
|-------|---------|
| `version` | Pack schema (`1`) |
| `tool` | Always `hermes` |
| `created_at` | UTC |
| `repo` | `origin` URL if set |
| `branch` | Current branch |
| `git_root` | Absolute repo path on the exporting machine |
| `message` | `--message` |
| `from_agent` | `claude`, `cursor`, … |
| `sources` | `git`, `manual`, plus agent id |
| `hermes_version` | CLI version |
| `session_path` | Source transcript if any |

Unknown `tool` values fail `hermes validate`.

Create a pack with `hermes handoff`. Reuse it with `hermes resume`. Do not rely on a shell glob such as `./handoff-*`.
