package pack

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Vinyl-Davyl/hermes/internal/manifest"
)

func TestLooksLikeGlob(t *testing.T) {
	if !LooksLikeGlob("./handoff-*") || !LooksLikeGlob("handoff-?") {
		t.Fatal("expected glob")
	}
	if LooksLikeGlob("./handoff-20260926T153022") {
		t.Fatal("real pack name is not a glob")
	}
}

func TestLatestAndResolve(t *testing.T) {
	dir := t.TempDir()
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldwd) })

	if _, err := Latest("."); !os.IsNotExist(err) {
		t.Fatalf("expected missing pack, got %v", err)
	}

	a := filepath.Join(dir, "handoff-aaa")
	b := filepath.Join(dir, "handoff-bbb")
	for _, p := range []string{a, b} {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
		m := manifest.New("x", "", "main", dir, "claude", "test")
		if err := manifest.Write(p, m); err != nil {
			t.Fatal(err)
		}
	}
	// Make b newer.
	if err := os.Chtimes(b, time.Now(), time.Now()); err != nil {
		t.Fatal(err)
	}

	got, err := Latest(".")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != "handoff-bbb" {
		t.Fatalf("got %s", got)
	}

	resolved, err := Resolve("./handoff-*")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(resolved) != "handoff-bbb" {
		t.Fatalf("glob resolve got %s", resolved)
	}
}

func TestResolveRealDir(t *testing.T) {
	dir := t.TempDir()
	if err := manifest.Write(dir, manifest.New("x", "", "main", dir, "", "test")); err != nil {
		t.Fatal(err)
	}
	got, err := Resolve(dir)
	if err != nil {
		t.Fatal(err)
	}
	abs, _ := filepath.Abs(dir)
	if got != abs {
		t.Fatalf("got %s want %s", got, abs)
	}
}
