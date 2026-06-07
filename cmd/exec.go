package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/cliwright/squawk/config"
	"github.com/cliwright/squawk/slack"
	"github.com/spf13/cobra"
)

var onFailure string
var onSuccess string
var execVars []string
var execTail int

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run a command and send a Slack alert on failure",
	Long:  "Wraps a command — sends the failure template if it exits non-zero, optionally sends a success template on exit 0.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(squawkDir)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		c := exec.Command(args[0], args[1:]...)
		c.Stdout = io.MultiWriter(os.Stdout, &buf)
		c.Stderr = io.MultiWriter(os.Stderr, &buf)

		runErr := c.Run()

		vars := envVars()
		for k, v := range buildVars(execVars) {
			vars[k] = v
		}
		vars["input"] = tailLines(buf.String(), execTail)

		var templateName string
		if runErr != nil {
			templateName = onFailure
		} else {
			if onSuccess == "" {
				return nil
			}
			templateName = onSuccess
		}

		tmpl, ok := cfg.Templates[templateName]
		if !ok {
			return fmt.Errorf("template %q not found — add it to .squawk/", templateName)
		}

		text, err := tmpl.Render(vars)
		if err != nil {
			return err
		}

		token, err := slack.Token()
		if err != nil {
			return err
		}

		if err := slack.Send(token, slack.NewMessage(tmpl.Channel, text, tmpl.Color)); err != nil {
			return err
		}

		if runErr != nil {
			if exitErr, ok := runErr.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			return runErr
		}

		return nil
	},
}

func init() {
	execCmd.Flags().StringVar(&onFailure, "on-failure", "failure", "failure template name from .squawk/")
	execCmd.Flags().StringVar(&onSuccess, "on-success", "", "success template name from .squawk/ (optional)")
	execCmd.Flags().StringSliceVar(&execVars, "var", nil, "additional template variables (key=value)")
	execCmd.Flags().IntVar(&execTail, "tail", 20, "number of output lines to include in the Slack message")
	rootCmd.AddCommand(execCmd)
}
