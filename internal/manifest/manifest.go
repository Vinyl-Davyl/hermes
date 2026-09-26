package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const Version = "1"

type Manifest struct {
	Version      string    `json:"version"`
	Tool         string    `json:"tool"`
	CreatedAt    time.Time `json:"created_at"`
	Repo         string    `json:"repo,omitempty"`
	Branch       string    `json:"branch,omitempty"`
	GitRoot      string    `json:"git_root,omitempty"`
	Message      string    `json:"message,omitempty"`
	FromAgent    string    `json:"from_agent,omitempty"`
	Sources      []string  `json:"sources"`
	HermesVersion string   `json:"hermes_version"`
	SessionPath  string    `json:"session_path,omitempty"`
	SessionTitle string    `json:"session_title,omitempty"`
}

func New(message, repo, branch, gitRoot, fromAgent, hermesVersion string) Manifest {
	sources := []string{"git", "manual"}
	if fromAgent != "" && fromAgent != "none" {
		sources = append(sources, fromAgent)
	}
	return Manifest{
		Version:       Version,
		Tool:          "hermes",
		CreatedAt:     time.Now().UTC(),
		Repo:          repo,
		Branch:        branch,
		GitRoot:       gitRoot,
		Message:       message,
		FromAgent:     fromAgent,
		Sources:       sources,
		HermesVersion: hermesVersion,
	}
}

func Write(dir string, m Manifest) error {
	path := filepath.Join(dir, "MANIFEST.json")
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func Read(dir string) (Manifest, error) {
	path := filepath.Join(dir, "MANIFEST.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("not a hermes pack (missing MANIFEST.json): %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("invalid MANIFEST.json: %w", err)
	}
	if err := Validate(m); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

func Validate(m Manifest) error {
	if m.Version == "" {
		return fmt.Errorf("manifest: version is required")
	}
	if m.Tool != "" && m.Tool != "hermes" {
		return fmt.Errorf("manifest: unknown tool %q", m.Tool)
	}
	return nil
}
