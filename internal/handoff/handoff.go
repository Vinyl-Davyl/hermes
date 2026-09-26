package handoff

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Vinyl-Davyl/hermes/internal/export"
	"github.com/Vinyl-Davyl/hermes/internal/importpkg"
	"github.com/Vinyl-Davyl/hermes/internal/pack"
	"github.com/Vinyl-Davyl/hermes/internal/version"
)

type Options struct {
	CWD          string
	Message      string
	From         string
	To           string
	SessionPath  string
	SessionID    string
	IncludeFiles bool
	Output       string
	Zip          bool
	Decisions    string
	Copy         bool
	Open         bool
	NoOpen       bool
}

type Result struct {
	Dir          string
	Zip          string
	FromAgent    string
	ToAgent      string
	Warnings     []string
	Instructions string
	Prompt       string
	Copied       bool
	Opened       bool
	SessionTitle string
	SessionID    string
}

func Run(opt Options) (Result, error) {
	var r Result
	if opt.To == "" {
		opt.To = "cursor"
	}
	if opt.From == "" {
		opt.From = "auto"
	}
	r.ToAgent = opt.To

	exp, err := export.Run(export.Options{
		CWD:           opt.CWD,
		Message:       opt.Message,
		From:          opt.From,
		SessionPath:   opt.SessionPath,
		SessionID:     opt.SessionID,
		IncludeFiles:  opt.IncludeFiles,
		Output:        opt.Output,
		Zip:           opt.Zip,
		Decisions:     opt.Decisions,
		HermesVersion: version.String,
	})
	if err != nil {
		return r, err
	}
	r.Dir = exp.Dir
	r.Zip = exp.Zip
	r.FromAgent = exp.FromAgent
	if exp.Session != nil {
		r.SessionTitle = exp.Session.Title
		r.SessionID = exp.Session.ID
	}
	r.Warnings = append(r.Warnings, exp.Warnings...)

	imp, err := importpkg.Run(exp.Dir, opt.To)
	if err != nil {
		return r, err
	}
	r.Prompt = imp.Prompt
	r.Instructions = imp.Instructions

	if opt.Copy || (needsPaste(opt.To) && !opt.NoOpen) {
		if err := copyClipboard(imp.Prompt); err != nil {
			r.Warnings = append(r.Warnings, "could not copy to clipboard: "+err.Error())
		} else {
			r.Copied = true
		}
	}

	shouldOpen := !opt.NoOpen && (opt.Open || autoLaunch(opt.To))
	if shouldOpen && CanLaunch(opt.To) {
		if err := SeedAndOpen(opt.To, opt.CWD, imp.Prompt); err != nil {
			r.Warnings = append(r.Warnings, err.Error())
		} else {
			r.Opened = true
		}
	}
	return r, nil
}

func Resume(packDir, target string) (importpkg.Result, error) {
	if target == "" {
		target = "cursor"
	}
	resolved, err := pack.Resolve(packDir)
	if err != nil {
		if os.IsNotExist(err) {
			return importpkg.Result{}, fmt.Errorf("no handoff pack here yet. Run: hermes handoff %s", target)
		}
		return importpkg.Result{}, err
	}
	return importpkg.Run(resolved, target)
}

func copyClipboard(s string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	default:
		if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd = exec.Command("wl-copy")
		} else {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		}
	}
	cmd.Stdin = strings.NewReader(s)
	return cmd.Run()
}
