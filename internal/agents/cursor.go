package agents

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Cursor struct{}

func (Cursor) ID() string   { return "cursor" }
func (Cursor) Name() string { return "Cursor" }

func cursorRoots() []string {
	home := Home()
	var paths []string
	if v := os.Getenv("CURSOR_CONFIG_DIR"); v != "" {
		paths = append(paths, v)
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		paths = append(paths, filepath.Join(xdg, "cursor"))
	}
	paths = append(paths, filepath.Join(home, ".cursor"))
	switch runtime.GOOS {
	case "darwin":
		paths = append(paths, filepath.Join(home, "Library", "Application Support", "Cursor", "User"))
	case "windows":
		paths = append(paths, filepath.Join(os.Getenv("APPDATA"), "Cursor", "User"))
	default:
		paths = append(paths, filepath.Join(home, ".config", "Cursor", "User"))
	}
	return uniqueExisting(paths...)
}

func (Cursor) Discover() ([]Session, error) {
	var out []Session
	for _, root := range cursorRoots() {
		// Prefer real composer DBs over random JSON in ~/.cursor
		for _, db := range findFiles(root, "state.vscdb") {
			info, err := os.Stat(db)
			if err != nil {
				continue
			}
			out = append(out, Session{
				Agent:    "cursor",
				ID:       filepath.Base(filepath.Dir(db)),
				Title:    cursorDBTitle(db),
				Path:     db,
				Modified: info.ModTime(),
			})
		}
	}
	return out, nil
}

func (a Cursor) Read(path string) (Transcript, error) {
	if strings.HasSuffix(path, ".vscdb") {
		md, n := readCursorDB(path)
		title := cursorDBTitle(path)
		return Transcript{
			Agent:    a.ID(),
			ID:       filepath.Base(filepath.Dir(path)),
			Title:    title,
			Path:     path,
			Markdown: md,
			Messages: n,
		}, nil
	}
	return readJSONLAsMarkdown(path, a.ID())
}

func cursorDBTitle(db string) string {
	ws := filepath.Join(filepath.Dir(db), "workspace.json")
	data, err := os.ReadFile(ws)
	if err == nil {
		if folder := extractJSONString(string(data), "folder"); folder != "" {
			folder = strings.TrimPrefix(folder, "file://")
			return filepath.Base(folder)
		}
	}
	return "Cursor workspace " + filepath.Base(filepath.Dir(db))
}

func readCursorDB(db string) (string, int) {
	if _, err := exec.LookPath("sqlite3"); err != nil {
		return "_Cursor session is in SQLite (`state.vscdb`). Install `sqlite3` or export the chat from Cursor, then re-run `hermes export --from cursor`._", 0
	}
	// Read-only; copy WAL siblings are left to sqlite3.
	query := `SELECT key FROM ItemTable WHERE key LIKE '%composer%' OR key LIKE '%aichat%' OR key LIKE '%aiService%' LIMIT 20;`
	out, err := exec.Command("sqlite3", db, query).Output()
	if err != nil {
		return "_Could not read Cursor database (is Cursor closed, or is this an older schema?). Git state is still in this pack._", 0
	}
	keys := strings.TrimSpace(string(out))
	if keys == "" {
		return "_No composer/chat keys found in this Cursor DB. Hermes still packed git state so you can continue in a new window._", 0
	}
	var b strings.Builder
	b.WriteString("Cursor stores chats inside the editor database. Hermes listed these keys:\n\n```\n")
	b.WriteString(keys)
	b.WriteString("\n```\n\n")
	b.WriteString("Hermes does **not** write back into Cursor. Open a new chat and `@PROMPT.md`.\n")
	return b.String(), strings.Count(keys, "\n") + 1
}

func findFiles(root, name string) []string {
	var out []string
	if !fileExists(root) {
		return out
	}
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && d.Name() == name {
			out = append(out, path)
		}
		if d.IsDir() && (d.Name() == "node_modules" || d.Name() == "CachedData") {
			return filepath.SkipDir
		}
		return nil
	})
	return out
}
