package main

import "testing"

func TestNormalizeRepoPath(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"plain", "architecture/003.md", "architecture/003.md", false},
		{"dot prefix", "./architecture/003.md", "architecture/003.md", false},
		{"repeated dot prefix", ".//architecture/003.md", "architecture/003.md", false},
		{"trailing slash", "architecture/003.md/", "architecture/003.md", false},
		{"duplicate slash", "architecture//003.md", "architecture/003.md", false},
		{"backslash", `architecture\003.md`, "architecture/003.md", false},
		{"dotdot inside root", "docs/specs/../../standards/001.md", "standards/001.md", false},
		{"dotdot stays inside", "docs/../standards/001.md", "standards/001.md", false},
		{"dotdot escapes", "../standards/001.md", "", true},
		{"dotdot escapes deep", "docs/../../standards/001.md", "", true},
		{"absolute", "/etc/passwd", "", true},
		{"drive letter", `C:\docs\001.md`, "", true},
		{"url", "https://example.com/001.md", "", true},
		{"empty", "", "", true},
		{"dot only", ".", "", true},
		{"case preserved", "Architecture/003.md", "Architecture/003.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeRepoPath(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("NormalizeRepoPath(%q) = %q, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("NormalizeRepoPath(%q) returned error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("NormalizeRepoPath(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
