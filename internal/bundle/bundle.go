package bundle

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Vinyl-Davyl/hermes/internal/agents"
	"github.com/Vinyl-Davyl/hermes/internal/fsutil"
	"github.com/Vinyl-Davyl/hermes/internal/git"
	"github.com/Vinyl-Davyl/hermes/internal/manifest"
	"github.com/Vinyl-Davyl/hermes/internal/prompt"
)

type Options struct {
	Dir           string
	Message       string
	IncludeFiles  bool
	FromAgent     string
	Transcript    *agents.Transcript
	Git           git.Snapshot
	HermesVersion string
	Decisions     string
}

func Write(opt Options) (string, error) {
	dir := opt.Dir
	if dir == "" {
		stamp := time.Now().UTC().Format("20060102T150405")
		dir = filepath.Join(".", "handoff-"+stamp)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(dir, "git"), 0o755); err != nil {
		return "", err
	}

	m := manifest.New(opt.Message, opt.Git.Remote, opt.Git.Branch, opt.Git.Root, opt.FromAgent, opt.HermesVersion)
	if opt.Transcript != nil {
		m.SessionPath = opt.Transcript.Path
		m.SessionTitle = opt.Transcript.Title
	}
	if err := manifest.Write(dir, m); err != nil {
		return "", err
	}

	_ = fsutil.WriteFile(filepath.Join(dir, "git", "branch.txt"), opt.Git.Branch+"\n")
	_ = fsutil.WriteFile(filepath.Join(dir, "git", "diff.patch"), opt.Git.Diff)
	_ = fsutil.WriteFile(filepath.Join(dir, "git", "log.txt"), opt.Git.Log+"\n")
	_ = fsutil.WriteFile(filepath.Join(dir, "git", "changed-files.txt"), strings.Join(opt.Git.ChangedFiles, "\n")+"\n")

	if opt.IncludeFiles && opt.Git.Root != "" {
		remaining := fsutil.DefaultMaxBytes
		for _, rel := range opt.Git.ChangedFiles {
			if fsutil.ShouldSkip(rel) {
				continue
			}
			src := filepath.Join(opt.Git.Root, rel)
			dst := filepath.Join(dir, "files", rel)
			_, _ = fsutil.CopyLimited(src, dst, &remaining, fsutil.DefaultMaxFile)
		}
	}

	sessionNote := ""
	if opt.Transcript != nil && opt.Transcript.Markdown != "" {
		_ = os.MkdirAll(filepath.Join(dir, "sessions"), 0o755)
		name := opt.Transcript.Agent
		if name == "" {
			name = "session"
		}
		_ = fsutil.WriteFile(filepath.Join(dir, "sessions", name+".md"), opt.Transcript.Markdown+"\n")
		sessionNote = fmt.Sprintf("Excerpt from %s (%d messages). Full excerpt in sessions/%s.md.\n",
			opt.Transcript.Agent, opt.Transcript.Messages, name)
	}

	if opt.Decisions != "" {
		_ = fsutil.WriteFile(filepath.Join(dir, "DECISIONS.md"), opt.Decisions+"\n")
	}

	_ = fsutil.WriteFile(filepath.Join(dir, "SUMMARY.md"), summary(opt.Message, opt.Git))

	p := prompt.Continuation(prompt.Data{
		Manifest:    m,
		Git:         opt.Git,
		SessionNote: sessionNote,
		Decisions:   opt.Decisions,
	})
	if err := fsutil.WriteFile(filepath.Join(dir, "PROMPT.md"), p+"\n"); err != nil {
		return "", err
	}
	return dir, nil
}

func Zip(dir string) (string, error) {
	zipPath := strings.TrimRight(dir, string(os.PathSeparator)) + ".zip"
	f, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(filepath.Dir(dir), path)
		if err != nil {
			return err
		}
		w, err := zw.Create(rel)
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(w, in)
		in.Close()
		return copyErr
	})
	if err != nil {
		zw.Close()
		return "", err
	}
	return zipPath, zw.Close()
}

func summary(message string, g git.Snapshot) string {
	var b strings.Builder
	b.WriteString("# Summary\n\n")
	if message != "" {
		b.WriteString(message)
		b.WriteString("\n\n")
	} else {
		b.WriteString("(no --message given)\n\n")
	}
	fmt.Fprintf(&b, "Branch: `%s`\n\n", g.Branch)
	if len(g.ChangedFiles) > 0 {
		b.WriteString("Changed files:\n\n")
		for _, f := range g.ChangedFiles {
			fmt.Fprintf(&b, "- %s\n", f)
		}
	}
	return b.String()
}
