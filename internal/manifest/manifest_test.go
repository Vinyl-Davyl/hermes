package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteReadValidate(t *testing.T) {
	dir := t.TempDir()
	m := New("fix jwt", "git@github.com:you/repo.git", "main", dir, "claude", "0.1.0")
	if err := Write(dir, m); err != nil {
		t.Fatal(err)
	}
	got, err := Read(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.Message != "fix jwt" || got.Branch != "main" || got.Tool != "hermes" {
		t.Fatalf("unexpected manifest: %+v", got)
	}
}

func TestReadMissing(t *testing.T) {
	if _, err := Read(t.TempDir()); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateUnknownTool(t *testing.T) {
	if err := Validate(Manifest{Version: "1", Tool: "baton"}); err == nil {
		t.Fatal("expected unknown tool error")
	}
}

func TestReadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "MANIFEST.json"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(dir); err == nil {
		t.Fatal("expected invalid json")
	}
}
