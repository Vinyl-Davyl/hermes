package agents

import (
	"path/filepath"
	"strings"
)

type Cline struct{}

func (Cline) ID() string   { return "cline" }
func (Cline) Name() string { return "Cline / Roo" }

func (Cline) Discover() ([]Session, error) {
	var out []Session
	for _, root := range clineRoots() {
		base := filepath.Base(root)
		// ~/.cline (or $CLINE_DIR) is a real session tree. VS Code
		// globalStorage is only useful for the known extension folders.
		if base == "globalStorage" {
			for _, dir := range []string{"saoudrizwan.claude-dev", "rooveterinaryinc.roo-cline"} {
				out = append(out, walkJSONL(filepath.Join(root, dir), "cline", 100)...)
			}
			continue
		}
		out = append(out, walkJSONL(root, "cline", 100)...)
		out = append(out, walkJSONL(filepath.Join(root, "data", "sessions"), "cline", 100)...)
	}
	return out, nil
}

func (a Cline) Read(path string) (Transcript, error) {
	t, err := readJSONLAsMarkdown(path, a.ID())
	t.ID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return t, err
}

type Kimi struct{}

func (Kimi) ID() string   { return "kimi" }
func (Kimi) Name() string { return "Kimi" }

func (Kimi) Discover() ([]Session, error) {
	var out []Session
	for _, root := range kimiRoots() {
		out = append(out, walkJSONL(root, "kimi", 100)...)
		out = append(out, walkJSONL(filepath.Join(root, "sessions"), "kimi", 100)...)
	}
	return out, nil
}

func (a Kimi) Read(path string) (Transcript, error) {
	t, err := readJSONLAsMarkdown(path, a.ID())
	t.ID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return t, err
}

type Antigravity struct{}

func (Antigravity) ID() string   { return "antigravity" }
func (Antigravity) Name() string { return "Antigravity" }

func (Antigravity) Discover() ([]Session, error) {
	var out []Session
	for _, root := range antigravityRoots() {
		out = append(out, walkJSONL(root, "antigravity", 100)...)
		out = append(out, walkJSONL(filepath.Join(root, "brain"), "antigravity", 100)...)
	}
	return out, nil
}

func (a Antigravity) Read(path string) (Transcript, error) {
	t, err := readJSONLAsMarkdown(path, a.ID())
	t.ID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return t, err
}

type OpenCode struct{}

func (OpenCode) ID() string   { return "opencode" }
func (OpenCode) Name() string { return "OpenCode" }

func (OpenCode) Discover() ([]Session, error) {
	var out []Session
	for _, root := range opencodeRoots() {
		out = append(out, walkJSONL(root, "opencode", 100)...)
		out = append(out, sqliteFileSession(filepath.Join(root, "opencode.db"), "opencode", "OpenCode database")...)
	}
	return out, nil
}

func (a OpenCode) Read(path string) (Transcript, error) {
	if strings.HasSuffix(path, ".db") || strings.HasSuffix(path, ".sqlite") {
		return readSQLiteHint(path, a.ID(), "OpenCode stores chats in SQLite. Hermes listed the database; paste PROMPT.md in a new OpenCode session.")
	}
	t, err := readJSONLAsMarkdown(path, a.ID())
	t.ID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return t, err
}

type Pi struct{}

func (Pi) ID() string   { return "pi" }
func (Pi) Name() string { return "Pi Agent" }

func (Pi) Discover() ([]Session, error) {
	var out []Session
	for _, root := range piRoots() {
		out = append(out, walkJSONL(root, "pi", 100)...)
		out = append(out, walkJSONL(filepath.Join(root, "sessions"), "pi", 100)...)
	}
	return out, nil
}

func (a Pi) Read(path string) (Transcript, error) {
	t, err := readJSONLAsMarkdown(path, a.ID())
	t.ID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return t, err
}
