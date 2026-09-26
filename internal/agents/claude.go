package agents

import (
	"os"
	"path/filepath"
	"strings"
)

type Claude struct{}

func (Claude) ID() string   { return "claude" }
func (Claude) Name() string { return "Claude Code" }

func (Claude) Discover() ([]Session, error) {
	root := filepath.Join(claudeRoot(), "projects")
	return walkJSONL(root, "claude", 200), nil
}

func (a Claude) Read(path string) (Transcript, error) {
	t, err := readJSONLAsMarkdown(path, a.ID())
	if err != nil {
		return t, err
	}
	t.ID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return t, nil
}

func EncodeClaudeProject(absPath string) string {
	absPath = filepath.Clean(absPath)
	enc := strings.ReplaceAll(absPath, string(os.PathSeparator), "-")
	enc = strings.TrimPrefix(enc, "-")
	return enc
}

func ClaudeProjectDir(cwd string) string {
	return filepath.Join(claudeRoot(), "projects", "-"+EncodeClaudeProject(cwd))
}
