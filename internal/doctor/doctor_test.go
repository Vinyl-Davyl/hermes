package doctor

import (
	"strings"
	"testing"

	"github.com/Vinyl-Davyl/hermes/internal/version"
)

func TestFormatIncludesVersion(t *testing.T) {
	out := Format(nil)
	if !strings.Contains(out, "hermes doctor  "+version.String) {
		t.Fatalf("expected version in doctor header, got %q", out)
	}
}
