package agents

import (
	"os"
	"path/filepath"
	"runtime"
)

// EnvRoot returns $key when set, otherwise fallback under the home directory.
func EnvRoot(key string, fallback ...string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return filepath.Join(append([]string{Home()}, fallback...)...)
}

func claudeRoot() string {
	return EnvRoot("CLAUDE_CONFIG_DIR", ".claude")
}

func codexRoot() string {
	return EnvRoot("CODEX_HOME", ".codex")
}

func kimiRoots() []string {
	return uniqueExisting(
		os.Getenv("KIMI_CODE_HOME"),
		filepath.Join(Home(), ".kimi-code"),
		filepath.Join(Home(), ".kimi"),
		filepath.Join(Home(), ".kimi-cli"),
		filepath.Join(Home(), ".moonshot"),
	)
}

func clineRoots() []string {
	home := Home()
	paths := []string{
		os.Getenv("CLINE_DIR"),
		filepath.Join(home, ".cline"),
	}
	switch runtime.GOOS {
	case "darwin":
		paths = append(paths,
			filepath.Join(home, "Library", "Application Support", "Code", "User", "globalStorage"),
			filepath.Join(home, "Library", "Application Support", "Cursor", "User", "globalStorage"),
		)
	case "windows":
		paths = append(paths,
			filepath.Join(os.Getenv("APPDATA"), "Code", "User", "globalStorage"),
			filepath.Join(os.Getenv("APPDATA"), "Cursor", "User", "globalStorage"),
		)
	default:
		paths = append(paths,
			filepath.Join(home, ".config", "Code", "User", "globalStorage"),
			filepath.Join(home, ".config", "Cursor", "User", "globalStorage"),
		)
	}
	return uniqueExisting(paths...)
}

func antigravityRoots() []string {
	home := Home()
	paths := []string{
		filepath.Join(home, ".gemini", "antigravity-cli"),
		filepath.Join(home, ".antigravity"),
	}
	switch runtime.GOOS {
	case "darwin":
		paths = append(paths, filepath.Join(home, "Library", "Application Support", "Antigravity"))
	case "windows":
		paths = append(paths, filepath.Join(os.Getenv("APPDATA"), "Antigravity"))
	default:
		paths = append(paths, filepath.Join(home, ".config", "Antigravity"))
	}
	return uniqueExisting(paths...)
}

func opencodeRoots() []string {
	home := Home()
	xdg := os.Getenv("XDG_DATA_HOME")
	paths := []string{
		filepath.Join(home, ".local", "share", "opencode"),
		filepath.Join(home, ".opencode"),
	}
	if xdg != "" {
		paths = append([]string{filepath.Join(xdg, "opencode")}, paths...)
	}
	return uniqueExisting(paths...)
}

func piRoots() []string {
	return uniqueExisting(
		os.Getenv("PI_CODING_AGENT_DIR"),
		filepath.Join(Home(), ".pi", "agent"),
		filepath.Join(Home(), ".pi"),
		filepath.Join(Home(), ".pi-go"),
		filepath.Join(Home(), ".local", "share", "pi"),
	)
}

func copilotRoot() string {
	return EnvRoot("COPILOT_HOME", ".copilot")
}

func zcodeRoot() string {
	if v := os.Getenv("ZCODE_HOME"); v != "" {
		return v
	}
	return filepath.Join(Home(), ".zcode", "cli", "db")
}

func deepseekRoot() string {
	return EnvRoot("DSH_HOME", ".dsh")
}

func uniqueExisting(paths ...string) []string {
	seen := map[string]bool{}
	var out []string
	for _, p := range paths {
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		if fileExists(p) {
			out = append(out, p)
		}
	}
	return out
}
