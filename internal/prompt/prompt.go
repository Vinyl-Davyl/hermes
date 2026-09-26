package prompt

import (
	"fmt"
	"strings"

	"github.com/Vinyl-Davyl/hermes/internal/git"
	"github.com/Vinyl-Davyl/hermes/internal/manifest"
)

type Data struct {
	Manifest    manifest.Manifest
	Git         git.Snapshot
	SessionNote string
	Decisions   string
	Target      string
}

func Continuation(d Data) string {
	var b strings.Builder
	b.WriteString("# Hermes handoff\n\n")
	b.WriteString("You are continuing work from another coding agent / window. ")
	b.WriteString("Do **not** restart from scratch. Use the git state and notes below.\n\n")

	if d.Manifest.Message != "" {
		b.WriteString("## Goal\n\n")
		b.WriteString(d.Manifest.Message)
		b.WriteString("\n\n")
	}

	b.WriteString("## Repo\n\n")
	if d.Manifest.Repo != "" {
		fmt.Fprintf(&b, "- Remote: `%s`\n", d.Manifest.Repo)
	}
	if d.Git.Root != "" {
		fmt.Fprintf(&b, "- Path: `%s`\n", d.Git.Root)
	}
	if d.Git.Branch != "" {
		fmt.Fprintf(&b, "- Branch: `%s`\n", d.Git.Branch)
	}
	if d.Manifest.FromAgent != "" {
		fmt.Fprintf(&b, "- Came from: `%s`\n", d.Manifest.FromAgent)
	}
	b.WriteString("\n")

	if len(d.Git.ChangedFiles) > 0 {
		b.WriteString("## Files in play\n\n")
		for _, f := range d.Git.ChangedFiles {
			fmt.Fprintf(&b, "- `%s`\n", f)
		}
		b.WriteString("\nRead these first. The patch is in `git/diff.patch`.\n\n")
	}

	if d.Decisions != "" {
		b.WriteString("## Decisions / do not redo\n\n")
		b.WriteString(d.Decisions)
		b.WriteString("\n\n")
	}

	if d.SessionNote != "" {
		b.WriteString("## Session excerpt\n\n")
		b.WriteString(d.SessionNote)
		b.WriteString("\n\n")
	}

	b.WriteString("## Next\n\n")
	b.WriteString("1. Confirm the goal and current branch.\n")
	b.WriteString("2. Inspect the named files and `git/diff.patch`.\n")
	b.WriteString("3. Continue the unfinished work. Ask only if a decision is missing.\n")

	if extra := targetHint(d.Target); extra != "" {
		b.WriteString("\n")
		b.WriteString(extra)
	}
	return b.String()
}

func ResumeInstructions(target, packDir string, files []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Pack: %s\n\n", packDir)
	patch := has(files, "git/diff.patch")
	switch strings.ToLower(target) {
	case "cursor":
		b.WriteString("Open a new Cursor chat (this window or another). Attach:\n\n")
		b.WriteString("  @PROMPT.md\n")
		if patch {
			b.WriteString("  @git/diff.patch\n")
		}
		b.WriteString("\nCursor cannot import a native session. This pack is the continue path.\n")
	case "claude":
		b.WriteString("If `claude` is on PATH, Hermes starts it and points it at .hermes/seed.md.\n")
		b.WriteString("Otherwise start Claude Code in this repo and paste PROMPT.md.\n")
	case "codex":
		b.WriteString("If `codex` is on PATH, Hermes starts it and points it at .hermes/seed.md.\n")
		b.WriteString("Otherwise start Codex and paste PROMPT.md.\n")
	case "opencode":
		b.WriteString("Start OpenCode in the repo and paste or @ PROMPT.md.\n")
	case "cline":
		b.WriteString("Open Cline in this repo and paste PROMPT.md into a new task.\n")
	case "kimi":
		b.WriteString("Start Kimi in the repo and paste PROMPT.md.\n")
	case "antigravity", "agy":
		b.WriteString("Open Antigravity in this repo and paste PROMPT.md.\n")
	case "pi":
		b.WriteString("Start Pi Agent in the repo and paste PROMPT.md.\n")
	case "copilot":
		b.WriteString("In the repo:\n\n")
		b.WriteString("  copilot\n\n")
		b.WriteString("Paste PROMPT.md as the first message.\n")
	case "zcode":
		b.WriteString("Open ZCode in this repo and paste PROMPT.md.\n")
	case "deepseek":
		b.WriteString("Start DeepSeek Harness in the repo and paste PROMPT.md.\n")
	default:
		b.WriteString("Open the next agent in this repo and paste PROMPT.md.\n")
		if patch {
			b.WriteString("Attach git/diff.patch if the agent accepts files.\n")
		}
	}
	return b.String()
}

func targetHint(target string) string {
	switch strings.ToLower(target) {
	case "cursor":
		return "## Cursor note\n\nThis is a new chat, not the old composer id. Continuity is the work, not the thread id.\n"
	default:
		return ""
	}
}

func has(files []string, name string) bool {
	for _, f := range files {
		if f == name {
			return true
		}
	}
	return false
}
