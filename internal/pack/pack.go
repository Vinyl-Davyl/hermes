package pack

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Vinyl-Davyl/hermes/internal/manifest"
)

// LooksLikeGlob reports whether s is a shell glob, not a real path.
// zsh errors with "no matches found" before Hermes even starts when
// ./handoff-* matches nothing — callers should use Latest instead.
func LooksLikeGlob(s string) bool {
	return strings.ContainsAny(s, "*?[")
}

// Latest returns the newest hermes pack under dir (directories named handoff-*).
func Latest(dir string) (string, error) {
	if dir == "" {
		dir = "."
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	type hit struct {
		path string
		mod  time.Time
	}
	var found []hit
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), "handoff-") {
			continue
		}
		p := filepath.Join(dir, e.Name())
		if _, err := manifest.Read(p); err != nil {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		found = append(found, hit{path: p, mod: info.ModTime()})
	}
	if len(found) == 0 {
		return "", os.ErrNotExist
	}
	sort.Slice(found, func(i, j int) bool { return found[i].mod.After(found[j].mod) })
	return found[0].path, nil
}

// Resolve turns a user path into a pack directory.
// Empty, "latest", ".", or a glob like ./handoff-* all mean "newest pack here".
func Resolve(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" || path == "latest" || path == "." || LooksLikeGlob(path) {
		return Latest(".")
	}
	if st, err := os.Stat(path); err == nil && st.IsDir() {
		if _, err := manifest.Read(path); err != nil {
			return "", err
		}
		return filepath.Abs(path)
	}
	return "", os.ErrNotExist
}
