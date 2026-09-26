package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestSnapshotAt(t *testing.T) {
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
			t.Fatalf("git %v: %s", args, out)
		}
	}
	run("init")
	run("checkout", "-b", "feat")
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "init")
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\nfunc X() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	root, err := Root(dir)
	if err != nil {
		t.Fatal(err)
	}
	snap, err := SnapshotAt(root)
	if err != nil {
		t.Fatal(err)
	}
	if snap.Branch != "feat" {
		t.Fatalf("branch %q", snap.Branch)
	}
	if !snap.Dirty || len(snap.ChangedFiles) == 0 {
		t.Fatalf("expected dirty tree, got %+v", snap)
	}
}

func TestRootNotGit(t *testing.T) {
	if _, err := Root(t.TempDir()); err == nil {
		t.Fatal("expected error")
	}
}
