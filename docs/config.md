# Config Reference

Squawk looks for `.squawk/*.yaml` files in the current working directory.

## Example

```yaml
templates:
  deploy-failed:
    channel: "#alerts"
    color: "#CC0000"
    mentions:
      - "U12345678"
    text: |
      ❌ *{{ .repo }}* failed on `{{ .branch }}`
      {{ .input }}
      <{{ .run_url }}|View run>
      {{ .mentions }}
```

## Fields

| Field | Description |
|---|---|
| `channel` | Slack channel to post to (e.g. `#alerts`) |
| `color` | Sidebar color for the attachment (e.g. `#CC0000`) |
| `mentions` | List of Slack user IDs to mention (optional) |
| `text` | Go `text/template` string |

## Template Variables

Variables are auto-populated from the CI environment. No flags required in most cases.

| Variable | Source |
|---|---|
| `.branch` | `GITHUB_REF_NAME`, `CI_COMMIT_REF_NAME`, `CIRCLE_BRANCH`, `TRAVIS_BRANCH`, `DRONE_BRANCH`, `BUILDKITE_BRANCH`, `BITBUCKET_BRANCH`, `BRANCH_NAME`, `BUILD_SOURCEBRANCHNAME`, or `git branch` |
| `.repo` | `GITHUB_REPOSITORY`, `CI_PROJECT_NAME`, `CIRCLE_PROJECT_REPONAME`, `TRAVIS_REPO_SLUG`, `DRONE_REPO`, `BUILDKITE_PIPELINE_SLUG`, `BITBUCKET_REPO_SLUG`, `JOB_NAME`, `BUILD_REPOSITORY_NAME`, or `git remote` |
| `.run_url` | Platform-specific CI run URL |
| `.input` | Piped stdin or command output (from `exec`) |
| `.mentions` | Auto-generated from `mentions` list in config |

Pass additional variables with `--var key=value`.

## Auth

Set `SQUAWK_SLACK_TOKEN` or `SLACK_TOKEN` in the environment. The bot token needs the `chat:write` scope.
