
[![Sketch fonts](https://see.fontimg.com/api/rf5/BWWo5/YTZkZTMxNDlhNDEwNDZhZmFiZThhODFhNjA5N2U3NTgub3Rm/c3F1YXdr/typo-draft-demo.png?r=fs&h=250&w=2000&fg=5A3922&bg=FFFFFF&tb=1&s=125)](https://www.fontspace.com/typo-draft-font-f41179)

> A dead-simple CLI for sending Slack alerts from CI pipelines. One red light, not ten green lights.

## Install

```bash
go install github.com/cliwright/squawk@latest
```

Or download a pre-built binary from the [releases page](https://github.com/cliwright/squawk/releases).

## Quick Start

1. **Set your Slack token:**
   ```bash
   export SQUAWK_SLACK_TOKEN="xoxb-your-token"
   ```

2. **Create a config file:**
   ```bash
   squawk init
   ```
   This creates `.squawk/squawk.yaml` with starter templates.

3. **Preview a message without sending:**
   ```bash
   squawk send --template failure --dry-run --var branch=main --var repo=my-app
   ```

## Usage

### `send`
Send a Slack alert using a named template.

```bash
# send a failure alert
squawk send --template deploy-failed

# pipe command output into the message
echo "$ERROR_LOG" | squawk send --template deploy-failed
```

### `exec`
Run a command and send a Slack alert based on the exit code.

```bash
# sends the failure template if the command exits non-zero
squawk exec --on-failure deploy-failed -- ./scripts/deploy.sh

# also send a success message on exit 0
squawk exec --on-failure deploy-failed --on-success deploy-succeeded -- ./scripts/deploy.sh
```

### `send --dry-run`
Render a template without sending it.

```bash
squawk send --template failure --dry-run --var branch=main
```

### `init`
Initialize a `.squawk` directory with starter templates.

```bash
squawk init
```

## Config

Squawk looks for `.squawk/*.yaml` files in the current working directory.

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

| Field | Description |
|---|---|
| `channel` | Slack channel to post to (e.g. `#alerts`) |
| `color` | Sidebar color for the attachment (e.g. `#CC0000`) |
| `mentions` | List of Slack user IDs to mention (optional) |
| `text` | Go `text/template` string |

### Template Variables

Variables are auto-populated from the CI environment. No flags required in most cases.

| Variable | Source |
|---|---|
| `.branch` | `GITHUB_REF_NAME`, `CI_COMMIT_REF_NAME`, `CIRCLE_BRANCH`, `TRAVIS_BRANCH`, `DRONE_BRANCH`, `BUILDKITE_BRANCH`, `BITBUCKET_BRANCH`, `BRANCH_NAME`, `BUILD_SOURCEBRANCHNAME`, or `git branch` |
| `.repo` | `GITHUB_REPOSITORY`, `CI_PROJECT_NAME`, `CIRCLE_PROJECT_REPONAME`, `TRAVIS_REPO_SLUG`, `DRONE_REPO`, `BUILDKITE_PIPELINE_SLUG`, `BITBUCKET_REPO_SLUG`, `JOB_NAME`, `BUILD_REPOSITORY_NAME`, or `git remote` |
| `.run_url` | Platform-specific CI run URL |
| `.input` | Piped stdin or command output (from `exec`) |
| `.mentions` | Auto-generated from `mentions` list in config |

Pass additional variables with `--var key=value`.

## Supported CI Platforms

- GitHub Actions
- GitLab CI
- CircleCI
- Travis CI
- Jenkins
- Azure Pipelines
- Bitbucket Pipelines
- Buildkite
- Drone CI

Works identically in CI and locally. Falls back to `git` commands when no CI env vars are detected.

## Auth

Set `SQUAWK_SLACK_TOKEN` or `SLACK_TOKEN` in the environment. The bot token needs the `chat:write` scope.

## License

Apache 2.0
