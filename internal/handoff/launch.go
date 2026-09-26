package handoff

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const seedDirName = ".hermes"
const seedFileName = "seed.md"

// CanLaunch reports whether Hermes can start this agent as a CLI
// and pass it the pack (Catchup-style). Editors without a seedable
// CLI still get a pack + clipboard / @PROMPT.md.
func CanLaunch(target string) bool {
	return launchBin(target) != ""
}

func SeedAndOpen(target, cwd, prompt string) error {
	bin := launchBin(target)
	if bin == "" {
		return fmt.Errorf("%s has no CLI Hermes can start — open it and attach @PROMPT.md", target)
	}
	if _, err := exec.LookPath(bin); err != nil {
		return fmt.Errorf("%s is not on PATH. Install it, or paste PROMPT.md by hand", bin)
	}
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	rel, err := writeSeed(cwd, prompt)
	if err != nil {
		return err
	}
	lead := "Continue the work from the Hermes pack. Read " + rel + " first, then pick up where it left off. Do not restart from scratch."
	name, args := launchArgs(target, lead)
	cmd := exec.Command(name, args...)
	cmd.Dir = cwd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeSeed(dir, body string) (string, error) {
	out := filepath.Join(dir, seedDirName)
	if err := os.MkdirAll(out, 0o755); err != nil {
		return "", err
	}
	_ = os.WriteFile(filepath.Join(out, ".gitignore"), []byte("*\n"), 0o644)
	path := filepath.Join(out, seedFileName)
	if err := os.WriteFile(path, []byte(strings.TrimSpace(body)+"\n"), 0o644); err != nil {
		return "", err
	}
	return filepath.Join(seedDirName, seedFileName), nil
}

func launchArgs(target, lead string) (name string, args []string) {
	name = launchBin(target)
	switch strings.ToLower(target) {
	case "opencode", "open-code":
		return name, []string{"--prompt", lead}
	case "copilot", "copilot-cli":
		return name, []string{"-i", lead}
	case "cline":
		return name, []string{"-i", lead}
	case "antigravity", "agy":
		return name, []string{"-i", lead}
	default:
		return name, []string{lead}
	}
}

func launchBin(target string) string {
	switch strings.ToLower(target) {
	case "claude", "claude-code":
		return "claude"
	case "codex", "codex-cli":
		return "codex"
	case "opencode", "open-code":
		return "opencode"
	case "kimi", "kimi-cli":
		return "" // no interactive seed prompt
	case "pi", "pi-agent":
		return "pi"
	case "copilot", "copilot-cli":
		return "copilot"
	case "cline":
		return "cline"
	case "cursor":
		return "cursor-agent"
	case "antigravity", "agy":
		return "agy"
	default:
		return ""
	}
}

func needsPaste(target string) bool {
	switch strings.ToLower(target) {
	case "cursor", "antigravity", "agy", "zcode", "deepseek", "kimi":
		return true
	default:
		return launchBin(target) == ""
	}
}

// autoLaunch is true for CLI destinations that Hermes can start
// without being asked. Cursor stays paste/@PROMPT.md unless --open.
func autoLaunch(target string) bool {
	if strings.EqualFold(target, "cursor") {
		return false
	}
	return binOnPATH(target)
}

func binOnPATH(target string) bool {
	bin := launchBin(target)
	if bin == "" {
		return false
	}
	_, err := exec.LookPath(bin)
	return err == nil
}
