package importpkg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Vinyl-Davyl/hermes/internal/manifest"
)

func TestImportRefreshesPrompt(t *testing.T) {
	dir := t.TempDir()
	m := manifest.New("continue jwt", "origin", "main", dir, "claude", "0.1.0")
	if err := manifest.Write(dir, m); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "git", "changed-files.txt"), []byte("jwt.go\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Run(dir, "cursor")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Prompt, "continue jwt") {
		t.Fatalf("prompt:\n%s", res.Prompt)
	}
	if !strings.Contains(res.Instructions, "@PROMPT.md") {
		t.Fatalf("instructions:\n%s", res.Instructions)
	}
}
