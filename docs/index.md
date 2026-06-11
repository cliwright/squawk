[![Sketch fonts](https://see.fontimg.com/api/rf5/BWWo5/YTZkZTMxNDlhNDEwNDZhZmFiZThhODFhNjA5N2U3NTgub3Rm/c3F1YXdr/typo-draft-demo.png?r=fs&h=250&w=2000&fg=5A3922&bg=FFFFFF&tb=1&s=125)](https://www.fontspace.com/typo-draft-font-f41179)


A dead-simple CLI for sending Slack alerts from CI pipelines.

## Philosophy
The `squawk` philosophy is to squawk when something goes wrong. I don't need my slack channel to 
light up like a Christmas tree every time a a CI pipeline run succeeds.

> One red light, not ten green lights.

We don't need 10 green lights to tell us everything is fine. I just one one actionable red light to 
signal when something needs attention.

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

3. **Preview a message without sending:**
   ```bash
   squawk send --template failure --dry-run --var branch=main --var repo=my-app
   ```

## Commands

### `send`
Send a Slack alert using a named template.

```bash
squawk send --template deploy-failed
```

Pipe command output into the message:
```bash
echo "$ERROR_LOG" | squawk send --template deploy-failed
```

### `exec`
Run a command and send a Slack alert based on the exit code.

```bash
squawk exec --on-failure deploy-failed -- ./scripts/deploy.sh
```

Also send a success message on exit 0:
```bash
squawk exec --on-failure deploy-failed --on-success deploy-succeeded -- ./scripts/deploy.sh
```

### `init`
Initialize a `.squawk` directory with starter templates.

```bash
squawk init
```

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

Works identically in CI and locally, falling back to `git` commands when no CI env vars are detected.

## License

Apache 2.0
