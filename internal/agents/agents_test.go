package agents

import (
	"os"
	"path/filepath"
	"testing"
)

func TestByIDAliases(t *testing.T) {
	cases := map[string]string{
		"claude-code": "claude",
		"agy":         "antigravity",
		"dsh":         "deepseek",
		"copilot-cli": "copilot",
		"zai":         "zcode",
	}
	for in, want := range cases {
		a, ok := ByID(in)
		if !ok || a.ID() != want {
			t.Fatalf("%s alias failed: ok=%v id=%v", in, ok, a)
		}
	}
	if _, ok := ByID("nope"); ok {
		t.Fatal("unknown should fail")
	}
}

func TestAllIncludesNewAgents(t *testing.T) {
	want := map[string]bool{"copilot": true, "zcode": true, "deepseek": true}
	for _, a := range All() {
		delete(want, a.ID())
	}
	if len(want) > 0 {
		t.Fatalf("missing agents: %v", want)
	}
}

func TestSameProjectIgnoresOtherRepos(t *testing.T) {
	cwd := "/Users/mac/Documents/GitHub/Hermes"
	foreign := Session{
		Agent: "claude",
		ID:    "deum",
		Title: "Use the claude_design MCP",
		Path:  "/Users/mac/.claude/projects/-Users-mac-Desktop-credo-portal-deum-credo-client/83068cb9.jsonl",
	}
	here := Session{
		Agent: "claude",
		ID:    "local",
		Title: "work on hermes",
		Path:  "/Users/mac/.claude/projects/-Users-mac-Documents-GitHub-Hermes/abc.jsonl",
	}
	cursorHere := Session{Agent: "cursor", ID: "ws", Title: "Hermes", Path: "/tmp/state.vscdb"}
	if SameProject(foreign, cwd) {
		t.Fatal("must not take a session from another project")
	}
	if !SameProject(here, cwd) {
		t.Fatal("must take the Claude project for this folder")
	}
	if !SameProject(cursorHere, cwd) {
		t.Fatal("cursor title Hermes should match this folder")
	}
}

func TestFilter(t *testing.T) {
	sessions := []Session{
		{Agent: "claude", ID: "1", Title: "Fix JWT", Path: "/Users/mac/.claude/projects/-Users-mac-src/1.jsonl"},
		{Agent: "codex", ID: "2", Title: "other", Path: "/tmp/other.jsonl"},
	}
	got := Filter(sessions, "jwt", "")
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("query: %+v", got)
	}
	got = Filter(sessions, "", "/Users/mac/src")
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("here: %+v", got)
	}
}

func TestFind(t *testing.T) {
	sessions := []Session{
		{Agent: "cursor", ID: "e188497675f3e6600dbe616bde83cacb", Title: "Hermes"},
		{Agent: "claude", ID: "abc-999", Title: "other"},
	}
	s, ok := Find(sessions, "e1884976")
	if !ok || s.Title != "Hermes" {
		t.Fatalf("prefix: %+v ok=%v", s, ok)
	}
	s, ok = Find(sessions, "Hermes")
	if !ok || s.ID != "e188497675f3e6600dbe616bde83cacb" {
		t.Fatalf("title: %+v ok=%v", s, ok)
	}
	if _, ok := Find(sessions, "nope"); ok {
		t.Fatal("expected miss")
	}
}

func TestEncodeClaudeProject(t *testing.T) {
	got := EncodeClaudeProject("/Users/mac/Documents/GitHub/pulse-net-go")
	if got != "Users-mac-Documents-GitHub-pulse-net-go" {
		t.Fatalf("got %q", got)
	}
}

func TestReadJSONL(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "s.jsonl")
	body := `{"type":"user","message":{"role":"user","content":[{"type":"text","text":"Fix the JWT"}]}}` + "\n" +
		`{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Looking at jwt.go"}]}}` + "\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	tr, err := Claude{}.Read(p)
	if err != nil {
		t.Fatal(err)
	}
	if tr.Messages != 2 {
		t.Fatalf("messages %d", tr.Messages)
	}
	if tr.Title != "Fix the JWT" {
		t.Fatalf("title %q", tr.Title)
	}
}

func TestShouldSkipSessionFile(t *testing.T) {
	if !shouldSkipSessionFile("package.json", "/x/package.json") {
		t.Fatal("should skip package.json")
	}
	if shouldSkipSessionFile("abc.jsonl", "/Users/mac/.claude/projects/x/abc.jsonl") {
		t.Fatal("should keep jsonl")
	}
	if !shouldSkipSessionFile("storage.json", "/x/globalStorage/storage.json") {
		t.Fatal("should skip storage.json")
	}
}
