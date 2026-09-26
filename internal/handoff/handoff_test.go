package handoff

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Vinyl-Davyl/hermes/internal/manifest"
)

func TestRunGitOnly(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := gitRepo(t)
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	res, err := Run(Options{
		CWD:     dir,
		Message: "try hermes",
		From:    "none",
		To:      "cursor",
		Output:  filepath.Join(dir, "handoff-local"),
		NoOpen:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.ToAgent != "cursor" {
		t.Fatalf("to %s", res.ToAgent)
	}
	if !strings.Contains(res.Instructions, "@PROMPT.md") {
		t.Fatalf("instructions:\n%s", res.Instructions)
	}
	if _, err := os.Stat(filepath.Join(res.Dir, "PROMPT.md")); err != nil {
		t.Fatal(err)
	}

	got, err := Resume("", "cursor")
	if err != nil {
		t.Fatal(err)
	}
	if got.PackDir == "" {
		t.Fatal("empty pack dir")
	}
}

func TestResumeMissing(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	if _, err := Resume("", "cursor"); err == nil {
		t.Fatal("expected missing pack")
	}
}

func TestResumeFindsPack(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	packDir := filepath.Join(dir, "handoff-test")
	if err := os.MkdirAll(filepath.Join(packDir, "git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := manifest.Write(packDir, manifest.New("go", "", "main", dir, "claude", "test")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packDir, "git", "changed-files.txt"), []byte("a.go\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Resume("./handoff-*", "cursor")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Prompt, "go") {
		t.Fatalf("prompt:\n%s", res.Prompt)
	}
}

func gitRepo(t *testing.T) string {
	t.Helper()
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
	return dir
}
