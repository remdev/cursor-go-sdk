package main

import "testing"

func TestNormalizeGitHubRemote(t *testing.T) {
	tests := []struct {
		name    string
		remote  string
		wantURL string
		wantOK  bool
	}{
		{
			name:    "https without git suffix",
			remote:  "https://github.com/cursor/cookbook",
			wantURL: "https://github.com/cursor/cookbook",
			wantOK:  true,
		},
		{
			name:    "https with git suffix",
			remote:  "https://github.com/cursor/cookbook.git",
			wantURL: "https://github.com/cursor/cookbook",
			wantOK:  true,
		},
		{
			name:    "ssh",
			remote:  "git@github.com:cursor/cookbook.git",
			wantURL: "https://github.com/cursor/cookbook",
			wantOK:  true,
		},
		{
			name:   "gitlab",
			remote: "git@gitlab.com:group/repo.git",
			wantOK: false,
		},
		{
			name:   "empty",
			remote: "  ",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, gotOK := normalizeGitHubRemote(tt.remote)
			if gotOK != tt.wantOK {
				t.Fatalf("ok: got %v want %v", gotOK, tt.wantOK)
			}
			if gotURL != tt.wantURL {
				t.Fatalf("url: got %q want %q", gotURL, tt.wantURL)
			}
		})
	}
}
