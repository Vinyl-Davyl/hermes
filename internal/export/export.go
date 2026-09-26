package export

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Vinyl-Davyl/hermes/internal/agents"
	"github.com/Vinyl-Davyl/hermes/internal/bundle"
	"github.com/Vinyl-Davyl/hermes/internal/git"
)

type Options struct {
	CWD           string
	Message       string
	From          string
	SessionPath   string
	SessionID     string
	IncludeFiles  bool
	Output        string
	Zip           bool
	Decisions     string
	HermesVersion string
}

type Result struct {
	Dir       string
	Zip       string
	FromAgent string
	Warnings  []string
	Session   *agents.Transcript
	Git       git.Snapshot
}

func Run(opt Options) (Result, error) {
	var r Result
	cwd := opt.CWD
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return r, err
		}
	}

	root, err := git.Root(cwd)
	var snap git.Snapshot
	if err != nil {
		r.Warnings = append(r.Warnings, "not a git repository; packing this folder without a diff (init git, or run from a repo, for a better pack)")
		abs, absErr := filepath.Abs(cwd)
		if absErr != nil {
			abs = cwd
		}
		snap = git.Snapshot{Root: abs}
	} else {
		snap, err = git.SnapshotAt(root)
		if err != nil {
			return r, err
		}
	}
	r.Git = snap

	from := opt.From
	if from == "" {
		from = "auto"
	}

	var tr *agents.Transcript
	if from != "none" {
		t, warn, err := pickTranscript(from, opt.SessionPath, opt.SessionID, snap.Root)
		if err != nil {
			r.Warnings = append(r.Warnings, err.Error())
		} else {
			tr = t
			r.Session = t
			if t != nil {
				r.FromAgent = t.Agent
			}
		}
		r.Warnings = append(r.Warnings, warn...)
	}

	dir, err := bundle.Write(bundle.Options{
		Dir:           opt.Output,
		Message:       opt.Message,
		IncludeFiles:  opt.IncludeFiles,
		FromAgent:     r.FromAgent,
		Transcript:    tr,
		Git:           snap,
		HermesVersion: opt.HermesVersion,
		Decisions:     opt.Decisions,
	})
	if err != nil {
		return r, err
	}
	r.Dir = dir

	if opt.Zip {
		z, err := bundle.Zip(dir)
		if err != nil {
			return r, err
		}
		r.Zip = z
	}
	return r, nil
}

func pickTranscript(from, explicit, sessionID, gitRoot string) (*agents.Transcript, []string, error) {
	var warnings []string
	if explicit != "" {
		ag, ok := inferAgent(from, explicit)
		if !ok {
			ag = agents.Claude{}
		}
		t, err := ag.Read(explicit)
		if err != nil {
			return nil, warnings, fmt.Errorf("read session: %w", err)
		}
		return &t, warnings, nil
	}

	if sessionID != "" {
		t, err := readByID(from, sessionID)
		if err != nil {
			return nil, warnings, err
		}
		return t, warnings, nil
	}

	if from != "auto" {
		ag, ok := agents.ByID(from)
		if !ok {
			return nil, warnings, fmt.Errorf("unknown agent %q (try: %s)", from, joinIDs())
		}
		t, warn := latestFor(ag, gitRoot)
		warnings = append(warnings, warn...)
		return t, warnings, nil
	}

	if claude, ok := agents.ByID("claude"); ok {
		if t, warn := latestFor(claude, gitRoot); t != nil {
			warnings = append(warnings, warn...)
			return t, warnings, nil
		}
	}
	all, _ := agents.DiscoverAll()
	here := agents.Filter(all, "", gitRoot)
	if len(here) == 0 {
		warnings = append(warnings, "no agent session for this folder; pack is git-only (this is fine). Use --id to pick a chat from another project")
		return nil, warnings, nil
	}
	ag, ok := agents.ByID(here[0].Agent)
	if !ok {
		return nil, warnings, nil
	}
	t, err := ag.Read(here[0].Path)
	if err != nil {
		warnings = append(warnings, "could not read latest session: "+err.Error())
		return nil, warnings, nil
	}
	return &t, warnings, nil
}

func latestFor(ag agents.Agent, gitRoot string) (*agents.Transcript, []string) {
	var warnings []string
	found, err := ag.Discover()
	if err != nil || len(found) == 0 {
		warnings = append(warnings, fmt.Sprintf("no %s sessions found", ag.Name()))
		return nil, warnings
	}
	sort.Slice(found, func(i, j int) bool { return found[i].Modified.After(found[j].Modified) })

	here := found
	if gitRoot != "" {
		here = agents.Filter(found, "", gitRoot)
		if len(here) == 0 && ag.ID() == "claude" {
			want := agents.ClaudeProjectDir(gitRoot)
			for _, s := range found {
				if filepathHasPrefix(s.Path, want) {
					here = []agents.Session{s}
					break
				}
			}
		}
	}
	if len(here) == 0 {
		warnings = append(warnings, fmt.Sprintf("no %s session for this folder (other projects were ignored). Pack is git-only, or pass --id", ag.Name()))
		return nil, warnings
	}
	pick := here[0]

	t, err := ag.Read(pick.Path)
	if err != nil {
		warnings = append(warnings, err.Error())
		return nil, warnings
	}
	return &t, warnings
}

func readByID(from, id string) (*agents.Transcript, error) {
	var pool []agents.Session
	if from != "" && from != "auto" && from != "none" {
		ag, ok := agents.ByID(from)
		if !ok {
			return nil, fmt.Errorf("unknown agent %q", from)
		}
		found, err := ag.Discover()
		if err != nil {
			return nil, err
		}
		pool = found
	} else {
		var err error
		pool, err = agents.DiscoverAll()
		if err != nil {
			return nil, err
		}
	}
	s, ok := agents.Find(pool, id)
	if !ok {
		return nil, fmt.Errorf("no session matching %q — run hermes list and copy the id", id)
	}
	ag, ok := agents.ByID(s.Agent)
	if !ok {
		return nil, fmt.Errorf("unknown agent %q", s.Agent)
	}
	t, err := ag.Read(s.Path)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func inferAgent(from, path string) (agents.Agent, bool) {
	if from != "" && from != "auto" {
		return agents.ByID(from)
	}
	return agents.ByID("claude")
}

func joinIDs() string {
	return fmt.Sprintf("%v", agents.IDs())
}

func filepathHasPrefix(path, prefix string) bool {
	if prefix == "" {
		return false
	}
	return len(path) >= len(prefix) && path[:len(prefix)] == prefix
}
