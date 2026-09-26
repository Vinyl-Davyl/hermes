package prompt

import (
	"strings"
	"testing"

	"github.com/Vinyl-Davyl/hermes/internal/git"
	"github.com/Vinyl-Davyl/hermes/internal/manifest"
)

func TestContinuationIncludesGoalAndFiles(t *testing.T) {
	out := Continuation(Data{
		Manifest: manifest.Manifest{Message: "Fix JWT expiry", FromAgent: "claude"},
		Git: git.Snapshot{
			Branch:       "main",
			ChangedFiles: []string{"internal/auth/jwt.go"},
		},
	})
	for _, need := range []string{"Fix JWT expiry", "internal/auth/jwt.go", "main", "claude"} {
		if !strings.Contains(out, need) {
			t.Fatalf("missing %q in:\n%s", need, out)
		}
	}
}

func TestResumeCursor(t *testing.T) {
	out := ResumeInstructions("cursor", "/tmp/pack", []string{"PROMPT.md", "git/diff.patch"})
	if !strings.Contains(out, "@PROMPT.md") || !strings.Contains(out, "new Cursor chat") {
		t.Fatalf("unexpected instructions:\n%s", out)
	}
}

func TestResumeClaudeAndCopilot(t *testing.T) {
	claude := ResumeInstructions("claude", "/tmp/pack", []string{"PROMPT.md"})
	if !strings.Contains(claude, "claude") {
		t.Fatalf("claude:\n%s", claude)
	}
	copilot := ResumeInstructions("copilot", "/tmp/pack", nil)
	if !strings.Contains(copilot, "copilot") {
		t.Fatalf("copilot:\n%s", copilot)
	}
}
