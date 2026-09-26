package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Vinyl-Davyl/hermes/internal/agents"
	"github.com/Vinyl-Davyl/hermes/internal/git"
)

type Check struct {
	Name   string
	OK     bool
	Detail string
}

func Run(cwd string) []Check {
	var checks []Check

	if _, err := exec.LookPath("git"); err != nil {
		checks = append(checks, Check{"git binary", false, "git not on PATH"})
	} else {
		checks = append(checks, Check{"git binary", true, "ok"})
	}

	if cwd == "" {
		cwd, _ = os.Getwd()
	}
	if root, err := git.Root(cwd); err != nil {
		checks = append(checks, Check{"git repo", false, "optional — not a git repo; handoff still packs this folder"})
	} else {
		snap, _ := git.SnapshotAt(root)
		detail := fmt.Sprintf("%s (branch %s, %d changed files)", root, snap.Branch, len(snap.ChangedFiles))
		checks = append(checks, Check{"git repo", true, detail})
	}

	if _, err := exec.LookPath("sqlite3"); err != nil {
		checks = append(checks, Check{"sqlite3 (Cursor DBs)", false, "optional — Cursor excerpts need sqlite3"})
	} else {
		checks = append(checks, Check{"sqlite3 (Cursor DBs)", true, "ok"})
	}

	for _, a := range agents.All() {
		found, err := a.Discover()
		if err != nil {
			checks = append(checks, Check{a.Name(), false, err.Error()})
			continue
		}
		checks = append(checks, Check{
			Name:   a.Name(),
			OK:     len(found) > 0,
			Detail: fmt.Sprintf("%d session(s)", len(found)),
		})
	}

	home, _ := os.UserHomeDir()
	if home != "" {
		checks = append(checks, Check{"home", true, home})
	}

	return checks
}

func Format(checks []Check) string {
	var b strings.Builder
	b.WriteString("hermes doctor\n\n")
	for _, c := range checks {
		mark := "✗"
		if c.OK {
			mark = "✓"
		}
		fmt.Fprintf(&b, "  %s  %-22s  %s\n", mark, c.Name, c.Detail)
	}
	b.WriteString("\nHermes never writes into Cursor/Claude databases.\n")
	b.WriteString("Missing sessions is OK — hermes handoff --from none still packs git.\n")
	b.WriteString("You do not need every agent installed to test Hermes.\n")
	return b.String()
}
