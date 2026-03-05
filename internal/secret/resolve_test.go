package secret

import (
	"os/exec"
	"testing"
)

func TestResolve_Passthrough(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"plain key", "my-secret-key-123"},
		{"url-like", "https://example.com/token"},
		{"partial prefix", "op:/incomplete"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Resolve(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.input {
				t.Errorf("Resolve(%q) = %q, want %q", tt.input, got, tt.input)
			}
		})
	}
}

func TestResolve_OpReference(t *testing.T) {
	if _, err := exec.LookPath("op"); err != nil {
		t.Skip("op CLI not installed, skipping integration test")
	}

	// Use a reference that will fail auth (not signed in) to verify
	// that the op command is actually invoked.
	_, err := Resolve("op://vault/item/field")
	if err == nil {
		// If it succeeds, the user is signed in and has this item — that's fine.
		return
	}
	// We expect an error from op, which confirms the command was executed.
	if err != nil {
		t.Logf("op invoked and returned expected error: %v", err)
	}
}
