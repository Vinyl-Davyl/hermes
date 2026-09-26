package importpkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Vinyl-Davyl/hermes/internal/git"
	"github.com/Vinyl-Davyl/hermes/internal/manifest"
	"github.com/Vinyl-Davyl/hermes/internal/prompt"
)

type Result struct {
	Manifest     manifest.Manifest
	Prompt       string
	Instructions string
	PackDir      string
}

func Run(packDir, target string) (Result, error) {
	var r Result
	if packDir == "" {
		return r, fmt.Errorf("missing pack directory — run hermes handoff, or hermes resume")
	}
	abs, err := filepath.Abs(packDir)
	if err != nil {
		return r, err
	}
	if _, err := os.Stat(abs); err != nil {
		return r, fmt.Errorf("no pack at %s — run hermes handoff first (do not use ./handoff-* in zsh; that glob fails when nothing matches)", packDir)
	}
	r.PackDir = abs

	m, err := manifest.Read(abs)
	if err != nil {
		return r, err
	}
	r.Manifest = m

	dec, _ := os.ReadFile(filepath.Join(abs, "DECISIONS.md"))
	sessionNote := ""
	if entries, err := os.ReadDir(filepath.Join(abs, "sessions")); err == nil && len(entries) > 0 {
		sessionNote = "See sessions/ in this pack for an excerpt of the previous agent chat."
	}

	snap := git.Snapshot{
		Root:   m.GitRoot,
		Branch: m.Branch,
		Remote: m.Repo,
	}
	if data, err := os.ReadFile(filepath.Join(abs, "git", "changed-files.txt")); err == nil {
		for _, line := range splitNonEmpty(string(data)) {
			snap.ChangedFiles = append(snap.ChangedFiles, line)
		}
	}

	r.Prompt = prompt.Continuation(prompt.Data{
		Manifest:    m,
		Git:         snap,
		SessionNote: sessionNote,
		Decisions:   string(dec),
		Target:      target,
	})
	_ = os.WriteFile(filepath.Join(abs, "PROMPT.md"), []byte(r.Prompt+"\n"), 0o644)

	files := []string{"PROMPT.md"}
	if _, err := os.Stat(filepath.Join(abs, "git", "diff.patch")); err == nil {
		files = append(files, "git/diff.patch")
	}
	if target == "" {
		target = "cursor"
	}
	r.Instructions = prompt.ResumeInstructions(target, abs, files)
	return r, nil
}

func splitNonEmpty(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			if line != "" {
				out = append(out, line)
			}
			start = i + 1
		}
	}
	return out
}
