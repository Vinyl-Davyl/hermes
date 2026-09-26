package agents

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Copilot struct{}

func (Copilot) ID() string   { return "copilot" }
func (Copilot) Name() string { return "Copilot CLI" }

func (Copilot) Discover() ([]Session, error) {
	root := copilotRoot()
	var out []Session
	out = append(out, walkJSONL(filepath.Join(root, "session-state"), "copilot", 100)...)
	out = append(out, walkJSONL(root, "copilot", 100)...)
	return out, nil
}

func (a Copilot) Read(path string) (Transcript, error) {
	t, err := readJSONLAsMarkdown(path, a.ID())
	t.ID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return t, err
}

type ZCode struct{}

func (ZCode) ID() string   { return "zcode" }
func (ZCode) Name() string { return "ZCode" }

func (ZCode) Discover() ([]Session, error) {
	root := zcodeRoot()
	var out []Session
	out = append(out, walkJSONL(root, "zcode", 50)...)
	out = append(out, sqliteFileSession(filepath.Join(root, "db.sqlite"), "zcode", "ZCode database")...)
	return out, nil
}

func (a ZCode) Read(path string) (Transcript, error) {
	if strings.HasSuffix(path, ".db") || strings.HasSuffix(path, ".sqlite") {
		return readSQLiteHint(path, a.ID(), "ZCode keeps chats in SQLite. Hermes does not write that database — start a new session and paste PROMPT.md.")
	}
	t, err := readJSONLAsMarkdown(path, a.ID())
	t.ID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return t, err
}

type DeepSeek struct{}

func (DeepSeek) ID() string   { return "deepseek" }
func (DeepSeek) Name() string { return "DeepSeek Harness" }

func (DeepSeek) Discover() ([]Session, error) {
	root := deepseekRoot()
	var out []Session
	out = append(out, walkJSONL(filepath.Join(root, "sessions"), "deepseek", 100)...)
	out = append(out, walkJSONL(root, "deepseek", 100)...)
	return out, nil
}

func (a DeepSeek) Read(path string) (Transcript, error) {
	t, err := readJSONLAsMarkdown(path, a.ID())
	t.ID = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return t, err
}

func sqliteFileSession(path, agent, title string) []Session {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	if hint := sqliteTitle(path); hint != "" {
		title = hint
	}
	return []Session{{
		Agent:    agent,
		ID:       strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Title:    title,
		Path:     path,
		Modified: info.ModTime(),
	}}
}

func sqliteTitle(db string) string {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		return ""
	}
	out, err := exec.Command("sqlite3", db, `SELECT title FROM session ORDER BY time_updated DESC LIMIT 1;`).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func readSQLiteHint(path, agent, note string) (Transcript, error) {
	title := sqliteTitle(path)
	if title == "" {
		title = filepath.Base(path)
	}
	return Transcript{
		Agent:    agent,
		ID:       strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Title:    title,
		Path:     path,
		Markdown: note,
		Messages: 1,
	}, nil
}
