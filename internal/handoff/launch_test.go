package handoff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteSeed(t *testing.T) {
	dir := t.TempDir()
	rel, err := writeSeed(dir, "# hello\n")
	if err != nil {
		t.Fatal(err)
	}
	if rel != ".hermes/seed.md" {
		t.Fatalf("rel %s", rel)
	}
	data, err := os.ReadFile(filepath.Join(dir, rel))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "hello") {
		t.Fatalf("body %s", data)
	}
	if _, err := os.Stat(filepath.Join(dir, ".hermes", ".gitignore")); err != nil {
		t.Fatal(err)
	}
}

func TestCanLaunch(t *testing.T) {
	if !CanLaunch("claude") || !CanLaunch("codex") {
		t.Fatal("expected CLI agents to launch")
	}
	if CanLaunch("zcode") || CanLaunch("deepseek") || CanLaunch("kimi") {
		t.Fatal("those should stay paste-only")
	}
}

func TestLaunchArgs(t *testing.T) {
	name, args := launchArgs("opencode", "go")
	if name != "opencode" || len(args) != 2 || args[0] != "--prompt" {
		t.Fatalf("%s %v", name, args)
	}
	name, args = launchArgs("claude", "go")
	if name != "claude" || len(args) != 1 {
		t.Fatalf("%s %v", name, args)
	}
}
