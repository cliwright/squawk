package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// buildVars parses key=value strings into a map.
func buildVars(raw []string) map[string]string {
	vars := make(map[string]string)
	for _, v := range raw {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars
}

// tailLines returns the last n lines of s.
func tailLines(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}

// envVars populates template variables from the CI environment.
func envVars() map[string]string {
	vars := make(map[string]string)

	// .branch
	if v := os.Getenv("GITHUB_REF_NAME"); v != "" {
		vars["branch"] = v
	} else if v := os.Getenv("CI_COMMIT_REF_NAME"); v != "" {
		vars["branch"] = v
	} else if out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output(); err == nil {
		vars["branch"] = strings.TrimSpace(string(out))
	}

	// .repo
	if v := os.Getenv("GITHUB_REPOSITORY"); v != "" {
		vars["repo"] = v
	} else if v := os.Getenv("CI_PROJECT_NAME"); v != "" {
		vars["repo"] = v
	} else if out, err := exec.Command("git", "remote", "get-url", "origin").Output(); err == nil {
		vars["repo"] = strings.TrimSpace(string(out))
	}

	// .run_url
	serverURL := os.Getenv("GITHUB_SERVER_URL")
	repo := os.Getenv("GITHUB_REPOSITORY")
	runID := os.Getenv("GITHUB_RUN_ID")
	if serverURL != "" && repo != "" && runID != "" {
		vars["run_url"] = fmt.Sprintf("%s/%s/actions/runs/%s", serverURL, repo, runID)
	} else if v := os.Getenv("CI_JOB_URL"); v != "" {
		vars["run_url"] = v
	}

	return vars
}
