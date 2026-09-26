package agents

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Session struct {
	Agent    string
	ID       string
	Title    string
	Path     string
	Modified time.Time
}

type Transcript struct {
	Agent    string
	ID       string
	Title    string
	Path     string
	Markdown string
	Messages int
}

type Agent interface {
	ID() string
	Name() string
	Discover() ([]Session, error)
	Read(path string) (Transcript, error)
}

func All() []Agent {
	return []Agent{
		Claude{},
		Codex{},
		Cursor{},
		Cline{},
		Kimi{},
		Antigravity{},
		OpenCode{},
		Pi{},
		Copilot{},
		ZCode{},
		DeepSeek{},
	}
}

func ByID(id string) (Agent, bool) {
	id = strings.ToLower(strings.TrimSpace(id))
	aliases := map[string]string{
		"claude-code":      "claude",
		"claude_code":      "claude",
		"codex-cli":        "codex",
		"pi-agent":         "pi",
		"pi-go":            "pi",
		"kimi-cli":         "kimi",
		"kimi-code":        "kimi",
		"open-code":        "opencode",
		"anti":             "antigravity",
		"agy":              "antigravity",
		"copilot-cli":      "copilot",
		"github-copilot":   "copilot",
		"z.ai":             "zcode",
		"zai":              "zcode",
		"dsh":              "deepseek",
		"deepseek-harness": "deepseek",
	}
	if mapped, ok := aliases[id]; ok {
		id = mapped
	}
	for _, a := range All() {
		if a.ID() == id {
			return a, true
		}
	}
	return nil, false
}

func IDs() []string {
	var ids []string
	for _, a := range All() {
		ids = append(ids, a.ID())
	}
	return ids
}

func DiscoverAll() ([]Session, error) {
	var all []Session
	for _, a := range All() {
		found, err := a.Discover()
		if err != nil {
			continue
		}
		all = append(all, found...)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Modified.After(all[j].Modified)
	})
	return all, nil
}

func Home() string {
	h, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return h
}

func walkJSONL(root, agent string, max int) []Session {
	var out []Session
	if root == "" {
		return out
	}
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := strings.ToLower(d.Name())
		if !strings.HasSuffix(name, ".jsonl") && !strings.HasSuffix(name, ".json") {
			return nil
		}
		if strings.Contains(path, string(os.PathSeparator)+"subagents"+string(os.PathSeparator)) {
			return nil
		}
		if shouldSkipSessionFile(name, path) {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		out = append(out, Session{
			Agent:    agent,
			ID:       strings.TrimSuffix(d.Name(), filepath.Ext(d.Name())),
			Title:    firstLineHint(path),
			Path:     path,
			Modified: info.ModTime(),
		})
		if max > 0 && len(out) >= max {
			return filepath.SkipAll
		}
		return nil
	})
	sort.Slice(out, func(i, j int) bool {
		return out[i].Modified.After(out[j].Modified)
	})
	return out
}

func shouldSkipSessionFile(name, path string) bool {
	skip := []string{
		"package.json", "tsconfig.json", "package-lock.json",
		"extensions.json", "hooks.json", "settings.json",
		"schema.json", "manifest.json", "storage.json", "metadata.json",
		"state.json", "workspace.json",
	}
	for _, s := range skip {
		if name == s {
			return true
		}
	}
	low := strings.ToLower(path)
	for _, part := range []string{
		"/node_modules/", "/extensions/", "/skills-cursor/", "/canvases/",
		"/anysphere.cursor-commits/", "/checkpoints/", "/CachedData/",
		"/Cache/", "/logs/",
	} {
		if strings.Contains(low, part) {
			return true
		}
	}
	return false
}

func firstLineHint(path string) string {
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return ""
	}
	text := string(data)
	if len(text) > 4000 {
		text = text[:4000]
	}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Pull a user-ish snippet out of JSONL/JSON.
		if i := strings.Index(line, `"text"`); i >= 0 {
			if hint := extractJSONString(line, "text"); hint != "" {
				return truncate(hint, 80)
			}
		}
		if hint := extractJSONString(line, "content"); hint != "" && !strings.HasPrefix(hint, "{") {
			return truncate(hint, 80)
		}
		if strings.HasPrefix(line, "#") || (!strings.HasPrefix(line, "{") && !strings.HasPrefix(line, "[")) {
			return truncate(line, 80)
		}
	}
	return filepath.Base(path)
}

func extractJSONString(line, key string) string {
	needle := `"` + key + `"`
	i := strings.Index(line, needle)
	if i < 0 {
		return ""
	}
	rest := line[i+len(needle):]
	colon := strings.Index(rest, ":")
	if colon < 0 {
		return ""
	}
	rest = strings.TrimSpace(rest[colon+1:])
	if !strings.HasPrefix(rest, `"`) {
		return ""
	}
	rest = rest[1:]
	var b strings.Builder
	escaped := false
	for _, r := range rest {
		if escaped {
			b.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '"' {
			break
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

func truncate(s string, n int) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Filter keeps sessions matching a keyword and/or the current working tree.
func Filter(sessions []Session, query, cwd string) []Session {
	query = strings.ToLower(strings.TrimSpace(query))
	cwd = strings.TrimSpace(cwd)
	var out []Session
	for _, s := range sessions {
		blob := strings.ToLower(s.Title + " " + s.ID + " " + s.Path)
		if query != "" && !strings.Contains(blob, query) {
			continue
		}
		if cwd != "" && !SameProject(s, cwd) {
			continue
		}
		out = append(out, s)
	}
	return out
}

// Find returns the first session whose id, title, or path matches q.
// A prefix of the id is enough (hermes list prints the full id).
func Find(sessions []Session, q string) (Session, bool) {
	q = strings.TrimSpace(q)
	if q == "" {
		return Session{}, false
	}
	low := strings.ToLower(q)
	for _, s := range sessions {
		if s.ID == q || strings.HasPrefix(s.ID, q) {
			return s, true
		}
		blob := strings.ToLower(s.ID + " " + s.Title + " " + s.Path)
		if strings.Contains(blob, low) {
			return s, true
		}
	}
	return Session{}, false
}

// SameProject reports whether a session belongs to this working tree.
// A chat from another repo on the same machine must not match.
func SameProject(s Session, cwd string) bool {
	cwd = filepath.Clean(strings.TrimSpace(cwd))
	if cwd == "" || cwd == "." || cwd == string(filepath.Separator) {
		return false
	}
	lowPath := strings.ToLower(filepath.ToSlash(s.Path))
	lowCwd := strings.ToLower(filepath.ToSlash(cwd))
	enc := strings.ToLower(EncodeClaudeProject(cwd))
	if enc != "" && strings.Contains(lowPath, strings.ToLower("-"+enc)) {
		return true
	}
	if enc != "" && strings.Contains(lowPath, enc) {
		return true
	}
	if lowCwd != "" && strings.Contains(lowPath, lowCwd) {
		return true
	}
	if s.Agent == "cursor" && strings.EqualFold(strings.TrimSpace(s.Title), filepath.Base(cwd)) {
		return true
	}
	return false
}
