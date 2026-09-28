package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Vinyl-Davyl/hermes/internal/version"
)

func TestParseHandoffArgs(t *testing.T) {
	from, to, err := parseHandoffArgs(nil, "", "")
	if err != nil || from != "auto" || to != "cursor" {
		t.Fatalf("defaults: %s %s %v", from, to, err)
	}

	from, to, err = parseHandoffArgs([]string{"codex"}, "", "")
	if err != nil || from != "auto" || to != "codex" {
		t.Fatalf("one arg: %s %s %v", from, to, err)
	}

	from, to, err = parseHandoffArgs([]string{"claude", "cursor"}, "", "")
	if err != nil || from != "claude" || to != "cursor" {
		t.Fatalf("two args: %s %s %v", from, to, err)
	}

	from, to, err = parseHandoffArgs([]string{"cursor"}, "claude", "")
	if err != nil || from != "claude" || to != "cursor" {
		t.Fatalf("flag from + dest: %s %s %v", from, to, err)
	}

	if _, _, err := parseHandoffArgs([]string{"nope"}, "", ""); err == nil {
		t.Fatal("expected unknown dest")
	}
	if _, _, err := parseHandoffArgs([]string{"nope", "cursor"}, "", ""); err == nil {
		t.Fatal("expected unknown source")
	}
}

func TestRootHasHandoff(t *testing.T) {
	root := New()
	found := false
	for _, c := range root.Commands() {
		if c.Name() == "handoff" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing handoff command")
	}
}

func TestRootVersionFlag(t *testing.T) {
	root := New()
	if root.Version != version.String {
		t.Fatalf("root version %q", root.Version)
	}
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--version"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "hermes "+version.String) {
		t.Fatalf("version output %q", got)
	}
}
