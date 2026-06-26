package agent

import (
	"runtime"
	"testing"
)

func TestIsAuthFailure(t *testing.T) {
	tests := []struct {
		name     string
		agent    string
		exitCode int
		stderr   string
		want     bool
	}{
		{"gemini exit 41", "gemini", 41, "", true},
		{"gemini exit 41 with spawn EPERM follows platform behavior", "gemini", 41, "spawn EPERM", runtime.GOOS != "windows"},
		{"gemini exit 41 with relaunch failure follows platform behavior", "gemini", 41, "Failed to relaunch the CLI process", runtime.GOOS != "windows"},
		{"gemini exit 41 with auth stderr remains auth failure", "gemini", 41, "authentication required", true},
		{"gemini exit 0 is never auth failure", "gemini", 0, "", false},
		{"gemini other exit code", "gemini", 1, "", false},
		{"unknown agent with auth stderr", "unknown", 1, "api_key not set", true},
		{"unknown agent no auth signal", "unknown", 1, "something failed", false},
		{"stderr unauthorized", "codex", 1, "Error: Unauthorized", true},
		{"stderr 401", "claude", 1, "HTTP 401 response", true},
		{"stderr authentication required", "gemini", 1, "authentication required", true},
		{"stderr invalid credentials", "codex", 1, "invalid credentials", true},
		{"stderr bare credentials is not auth failure", "codex", 1, "credential helper error", false},
		{"exit 0 ignores auth stderr", "codex", 0, "api_key not set", false},
		{"case insensitive stderr", "claude", 1, "UNAUTHORIZED access", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAuthFailure(tt.agent, tt.exitCode, tt.stderr)
			if got != tt.want {
				t.Errorf("IsAuthFailure(%q, %d, %q) = %v, want %v",
					tt.agent, tt.exitCode, tt.stderr, got, tt.want)
			}
		})
	}
}

func TestIsAuthFailureStdout(t *testing.T) {
	tests := []struct {
		name     string
		agent    string
		exitCode int
		stderr   string
		stdout   string
		want     bool
	}{
		{
			"structured JSON 401",
			"agy", 1, "",
			`{"is_error":true,"api_error_status":401,"result":"authentication required"}`,
			true,
		},
		{
			"structured JSON 403",
			"agy", 1, "",
			`{"is_error":true,"api_error_status":403,"error":"forbidden"}`,
			true,
		},
		{
			"structured JSON 401 but is_error false",
			"agy", 1, "",
			`{"is_error":false,"api_error_status":401,"result":"ok"}`,
			false,
		},
		{
			"short text not authenticated",
			"agy", 1, "",
			"not authenticated",
			true,
		},
		{
			"short text login required",
			"codex", 1, "",
			"login required",
			true,
		},
		{
			"stdout prefix api error 401",
			"claude", 1, "",
			"api error: 401 Unauthorized",
			true,
		},
		{
			"long normal output is not auth failure",
			"codex", 1, "",
			"This is a long review output with many lines.\nLine 2\nLine 3\nLine 4\nLine 5",
			false,
		},
		{
			"exit 0 ignores stdout auth",
			"agy", 0, "",
			`{"is_error":true,"api_error_status":401}`,
			false,
		},
		{
			"stderr auth takes precedence over empty stdout",
			"codex", 1, "unauthorized", "",
			true,
		},
		{
			"empty stdout is not auth failure",
			"agy", 1, "", "",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAuthFailure(tt.agent, tt.exitCode, tt.stderr, tt.stdout)
			if got != tt.want {
				t.Errorf("IsAuthFailure(%q, %d, stderr=%q, stdout=%q) = %v, want %v",
					tt.agent, tt.exitCode, tt.stderr, tt.stdout, got, tt.want)
			}
		})
	}
}

func TestAuthHintAgy(t *testing.T) {
	hint := AuthHint("agy")
	if hint == "" {
		t.Error("AuthHint(\"agy\") returned empty string")
	}
}

func TestAuthHint(t *testing.T) {
	agents := []string{"gemini", "claude", "codex", "unknown"}
	for _, name := range agents {
		t.Run(name, func(t *testing.T) {
			hint := AuthHint(name)
			if hint == "" {
				t.Errorf("AuthHint(%q) returned empty string", name)
			}
		})
	}
}
