package fsutil

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

var skipNames = map[string]bool{
	".env": true, ".DS_Store": true,
}

var skipExt = map[string]bool{
	".exe": true, ".so": true, ".dylib": true, ".bin": true,
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
	".mp4": true, ".zip": true, ".tar": true, ".gz": true,
}

const DefaultMaxBytes int64 = 2 * 1024 * 1024
const DefaultMaxFile int64 = 256 * 1024

func ShouldSkip(rel string) bool {
	base := filepath.Base(rel)
	if skipNames[base] || strings.HasPrefix(base, ".env.") {
		return true
	}
	ext := strings.ToLower(filepath.Ext(base))
	return skipExt[ext]
}

func CopyLimited(src, dst string, remaining *int64, maxFile int64) (copied bool, err error) {
	info, err := os.Stat(src)
	if err != nil {
		return false, err
	}
	if info.IsDir() || info.Size() > maxFile || info.Size() > *remaining {
		return false, nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return false, err
	}
	in, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return false, err
	}
	defer out.Close()
	n, err := io.Copy(out, in)
	if err != nil {
		return false, err
	}
	*remaining -= n
	return true, nil
}

func WriteFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
