package fsutil

import "testing"

func TestShouldSkip(t *testing.T) {
	if !ShouldSkip(".env") || !ShouldSkip("secret.png") {
		t.Fatal("expected skip")
	}
	if ShouldSkip("internal/auth/jwt.go") {
		t.Fatal("should keep go file")
	}
}
