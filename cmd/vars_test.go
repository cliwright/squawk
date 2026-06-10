package cmd

import (
	"testing"
)

func TestEnvVars(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		expected   map[string]string
		unexpected []string
	}{
		{
			name: "github actions full",
			env: map[string]string{
				"GITHUB_REF_NAME":   "main",
				"GITHUB_REPOSITORY": "cliwright/squawk",
				"GITHUB_SERVER_URL": "https://github.com",
				"GITHUB_RUN_ID":     "12345",
			},
			expected: map[string]string{
				"branch":  "main",
				"repo":    "cliwright/squawk",
				"run_url": "https://github.com/cliwright/squawk/actions/runs/12345",
			},
		},
		{
			name: "gitlab ci full",
			env: map[string]string{
				"CI_COMMIT_REF_NAME": "feature/x",
				"CI_PROJECT_NAME":    "squawk",
				"CI_JOB_URL":         "https://gitlab.com/cliwright/squawk/-/jobs/987",
			},
			expected: map[string]string{
				"branch":  "feature/x",
				"repo":    "squawk",
				"run_url": "https://gitlab.com/cliwright/squawk/-/jobs/987",
			},
		},
		{
			name: "circleci full",
			env: map[string]string{
				"CIRCLE_BRANCH":           "main",
				"CIRCLE_PROJECT_REPONAME": "squawk",
				"CIRCLE_BUILD_URL":        "https://circleci.com/gh/cliwright/squawk/123",
			},
			expected: map[string]string{
				"branch":  "main",
				"repo":    "squawk",
				"run_url": "https://circleci.com/gh/cliwright/squawk/123",
			},
		},
		{
			name: "travis ci full",
			env: map[string]string{
				"TRAVIS_BRANCH":        "main",
				"TRAVIS_REPO_SLUG":     "cliwright/squawk",
				"TRAVIS_BUILD_WEB_URL": "https://travis-ci.org/cliwright/squawk/builds/123",
			},
			expected: map[string]string{
				"branch":  "main",
				"repo":    "cliwright/squawk",
				"run_url": "https://travis-ci.org/cliwright/squawk/builds/123",
			},
		},
		{
			name: "jenkins full",
			env: map[string]string{
				"BRANCH_NAME": "main",
				"JOB_NAME":    "squawk",
				"BUILD_URL":   "https://jenkins.example.com/job/squawk/123/",
			},
			expected: map[string]string{
				"branch":  "main",
				"repo":    "squawk",
				"run_url": "https://jenkins.example.com/job/squawk/123/",
			},
		},
		{
			name: "azure pipelines full",
			env: map[string]string{
				"BUILD_SOURCEBRANCHNAME":         "main",
				"BUILD_REPOSITORY_NAME":          "squawk",
				"SYSTEM_TEAMFOUNDATIONSERVERURI": "https://dev.azure.com/cliwright/",
				"SYSTEM_TEAMPROJECT":             "squawk",
				"BUILD_BUILDID":                  "123",
			},
			expected: map[string]string{
				"branch":  "main",
				"repo":    "squawk",
				"run_url": "https://dev.azure.com/cliwright/squawk/_build/results?buildId=123",
			},
		},
		{
			name: "bitbucket pipelines full",
			env: map[string]string{
				"BITBUCKET_BRANCH":       "main",
				"BITBUCKET_REPO_SLUG":    "squawk",
				"BITBUCKET_WORKSPACE":    "cliwright",
				"BITBUCKET_BUILD_NUMBER": "123",
			},
			expected: map[string]string{
				"branch":  "main",
				"repo":    "squawk",
				"run_url": "https://bitbucket.org/cliwright/squawk/pipelines/results/123",
			},
		},
		{
			name: "buildkite full",
			env: map[string]string{
				"BUILDKITE_BRANCH":        "main",
				"BUILDKITE_PIPELINE_SLUG": "squawk",
				"BUILDKITE_BUILD_URL":     "https://buildkite.com/cliwright/squawk/builds/123",
			},
			expected: map[string]string{
				"branch":  "main",
				"repo":    "squawk",
				"run_url": "https://buildkite.com/cliwright/squawk/builds/123",
			},
		},
		{
			name: "drone ci full",
			env: map[string]string{
				"DRONE_BRANCH":     "main",
				"DRONE_REPO":       "cliwright/squawk",
				"DRONE_BUILD_LINK": "https://drone.example.com/cliwright/squawk/123",
			},
			expected: map[string]string{
				"branch":  "main",
				"repo":    "cliwright/squawk",
				"run_url": "https://drone.example.com/cliwright/squawk/123",
			},
		},
		{
			name: "github takes precedence over gitlab",
			env: map[string]string{
				"GITHUB_REF_NAME":    "gh-branch",
				"CI_COMMIT_REF_NAME": "gl-branch",
				"GITHUB_REPOSITORY":  "gh-repo",
				"CI_PROJECT_NAME":    "gl-repo",
				"GITHUB_SERVER_URL":  "https://github.com",
				"GITHUB_RUN_ID":      "1",
			},
			expected: map[string]string{
				"branch":  "gh-branch",
				"repo":    "gh-repo",
				"run_url": "https://github.com/gh-repo/actions/runs/1",
			},
		},
		{
			name: "partial github env falls back to gitlab for run_url",
			env: map[string]string{
				"GITHUB_REF_NAME":    "main",
				"CI_COMMIT_REF_NAME": "develop",
				"CI_JOB_URL":         "https://gitlab.com/jobs/1",
			},
			expected: map[string]string{
				"branch":  "main",
				"run_url": "https://gitlab.com/jobs/1",
			},
		},
		{
			name: "missing run_url env vars",
			env: map[string]string{
				"GITHUB_REF_NAME":   "main",
				"GITHUB_REPOSITORY": "cliwright/squawk",
			},
			expected: map[string]string{
				"branch": "main",
				"repo":   "cliwright/squawk",
			},
			unexpected: []string{"run_url"},
		},
		{
			name:       "empty ci environment",
			env:        map[string]string{},
			unexpected: []string{"run_url"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all CI env vars to avoid leakage from the host.
			for _, key := range []string{
				"GITHUB_REF_NAME",
				"GITHUB_REPOSITORY",
				"GITHUB_SERVER_URL",
				"GITHUB_RUN_ID",
				"CI_COMMIT_REF_NAME",
				"CI_PROJECT_NAME",
				"CI_JOB_URL",
				"CIRCLE_BRANCH",
				"CIRCLE_PROJECT_REPONAME",
				"CIRCLE_BUILD_URL",
				"TRAVIS_BRANCH",
				"TRAVIS_REPO_SLUG",
				"TRAVIS_BUILD_WEB_URL",
				"DRONE_BRANCH",
				"DRONE_REPO",
				"DRONE_BUILD_LINK",
				"BUILDKITE_BRANCH",
				"BUILDKITE_PIPELINE_SLUG",
				"BUILDKITE_BUILD_URL",
				"BITBUCKET_BRANCH",
				"BITBUCKET_REPO_SLUG",
				"BITBUCKET_WORKSPACE",
				"BITBUCKET_BUILD_NUMBER",
				"BRANCH_NAME",
				"JOB_NAME",
				"BUILD_URL",
				"BUILD_SOURCEBRANCHNAME",
				"BUILD_REPOSITORY_NAME",
				"SYSTEM_TEAMFOUNDATIONSERVERURI",
				"SYSTEM_TEAMPROJECT",
				"BUILD_BUILDID",
			} {
				t.Setenv(key, "")
			}

			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got := envVars()

			for k, want := range tt.expected {
				if got[k] != want {
					t.Errorf("envVars()[%q] = %q, want %q", k, got[k], want)
				}
			}

			for _, k := range tt.unexpected {
				if _, ok := got[k]; ok {
					t.Errorf("envVars() unexpected key %q = %q", k, got[k])
				}
			}
		})
	}
}
