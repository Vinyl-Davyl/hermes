package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Vinyl-Davyl/hermes/internal/manifest"
)

func TestExportWithoutGit(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "pack")
	res, err := Run(Options{
		CWD:           dir,
		Message:       "no git here",
		From:          "none",
		Output:        out,
		HermesVersion: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) == 0 {
		t.Fatal("expected a not-a-git-repo warning")
	}
	m, err := manifest.Read(res.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if m.Message != "no git here" {
		t.Fatalf("message %q", m.Message)
	}
	if _, err := os.Stat(filepath.Join(res.Dir, "PROMPT.md")); err != nil {
		t.Fatal(err)
	}
}
