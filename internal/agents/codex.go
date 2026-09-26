package agents

import (
	"path/filepath"
	"strings"
)

type Codex struct{}

func (Codex) ID() string   { return "codex" }
func (Codex) Name() string { return "Codex CLI" }

func (Codex) Discover() ([]Session, error) {
	root := filepath.Join(codexRoot(), "sessions")
	return walkJSONL(root, "codex", 200), nil
}

func (a Codex) Read(path string) (Transcript, error) {
	t, err := readJSONLAsMarkdown(path, a.ID())
	if err != nil {
		return t, err
	}
	t.ID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return t, nil
}
