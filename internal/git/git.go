package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Snapshot struct {
	Root         string
	Branch       string
	Remote       string
	Diff         string
	Log          string
	ChangedFiles []string
	Dirty        bool
}

func Root(cwd string) (string, error) {
	out, err := run(cwd, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not a git repository: %w", err)
	}
	return filepath.Clean(out), nil
}

func SnapshotAt(root string) (Snapshot, error) {
	s := Snapshot{Root: root}

	branch, err := run(root, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return s, err
	}
	s.Branch = branch

	if remote, err := run(root, "config", "--get", "remote.origin.url"); err == nil {
		s.Remote = remote
	}

	diff, err := runRaw(root, "diff", "HEAD")
	if err != nil {
		diff, _ = runRaw(root, "diff")
	}
	s.Diff = string(diff)

	log, err := run(root, "log", "-10", "--oneline", "--decorate")
	if err == nil {
		s.Log = log
	}

	changed, err := run(root, "diff", "--name-only", "HEAD")
	if err != nil {
		changed, _ = run(root, "diff", "--name-only")
	}
	if changed != "" {
		s.ChangedFiles = splitLines(changed)
	}

	untracked, err := run(root, "ls-files", "--others", "--exclude-standard")
	if err == nil && untracked != "" {
		for _, f := range splitLines(untracked) {
			if !contains(s.ChangedFiles, f) {
				s.ChangedFiles = append(s.ChangedFiles, f)
			}
		}
	}

	s.Dirty = s.Diff != "" || len(s.ChangedFiles) > 0
	return s, nil
}

func run(dir string, args ...string) (string, error) {
	out, err := runRaw(dir, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func runRaw(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
	}
	return out, nil
}

func splitLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
