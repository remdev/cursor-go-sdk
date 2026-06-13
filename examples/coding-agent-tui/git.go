package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type cloudRepository struct {
	URL         string
	StartingRef string
}

func detectCloudRepository(cwd string) (cloudRepository, error) {
	remote, err := runGit(cwd, "config", "--get", "remote.origin.url")
	if err != nil || remote == "" {
		return cloudRepository{}, fmt.Errorf("cloud mode requires a git repository with remote.origin.url set")
	}
	url, ok := normalizeGitHubRemote(remote)
	if !ok {
		return cloudRepository{}, fmt.Errorf("cloud mode currently expects remote.origin.url to point at GitHub")
	}
	branch, err := runGit(cwd, "rev-parse", "--abbrev-ref", "HEAD")
	if err == nil && branch != "" && branch != "HEAD" {
		return cloudRepository{URL: url, StartingRef: branch}, nil
	}
	return cloudRepository{URL: url}, nil
}

func formatCloudRepository(r cloudRepository) string {
	if r.StartingRef != "" {
		return r.URL + "#" + r.StartingRef
	}
	return r.URL
}

func runGit(cwd string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", cwd}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

var (
	githubHTTPS = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/.]+)(?:\.git)?/?$`)
	githubSSH   = regexp.MustCompile(`^git@github\.com:([^/]+)/([^/.]+)(?:\.git)?$`)
)

func normalizeGitHubRemote(remote string) (string, bool) {
	remote = strings.TrimSpace(remote)
	if m := githubHTTPS.FindStringSubmatch(remote); len(m) == 3 {
		return fmt.Sprintf("https://github.com/%s/%s", m[1], m[2]), true
	}
	if m := githubSSH.FindStringSubmatch(remote); len(m) == 3 {
		return fmt.Sprintf("https://github.com/%s/%s", m[1], m[2]), true
	}
	return "", false
}
