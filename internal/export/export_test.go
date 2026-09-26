package export

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Vinyl-Davyl/hermes/internal/manifest"
)

func TestExportGitOnly(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s", out)
		}
	}
	run("init")
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "init")

	out := filepath.Join(dir, "pack")
	res, err := Run(Options{
		CWD:           dir,
		Message:       "ship hermes",
		From:          "none",
		Output:        out,
		HermesVersion: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Read(res.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if m.Message != "ship hermes" {
		t.Fatalf("message %q", m.Message)
	}
	if _, err := os.Stat(filepath.Join(res.Dir, "PROMPT.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(res.Dir, "git", "diff.patch")); err != nil {
		t.Fatal(err)
	}
}
