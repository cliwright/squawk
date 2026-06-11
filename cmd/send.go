package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/cliwright/squawk/config"
	"github.com/cliwright/squawk/slack"
	"github.com/spf13/cobra"
)

var sendTemplate string
var sendVars []string
var sendDryRun bool

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a Slack alert using a named template",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(squawkDir)
		if err != nil {
			return err
		}

		tmpl, ok := cfg.Templates[sendTemplate]
		if !ok {
			return fmt.Errorf("template %q not found", sendTemplate)
		}

		vars := envVars()
		for k, v := range buildVars(sendVars) {
			vars[k] = v
		}

		if stat, err := os.Stdin.Stat(); err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
			stdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
			if len(stdin) > 0 {
				vars["input"] = string(stdin)
			}
		}

		text, err := tmpl.Render(vars)
		if err != nil {
			return err
		}

		if sendDryRun {
			fmt.Println("--- channel:", tmpl.Channel, "---")
			fmt.Println(text)
			return nil
		}

		token, err := slack.Token()
		if err != nil {
			return err
		}

		return slack.Send(token, slack.NewMessage(tmpl.Channel, text, tmpl.Color))
	},
}

func init() {
	sendCmd.Flags().StringVarP(&sendTemplate, "template", "t", "", "template name from squawk.yaml")
	sendCmd.Flags().StringSliceVar(&sendVars, "var", nil, "additional template variables (key=value)")
	sendCmd.Flags().BoolVar(&sendDryRun, "dry-run", false, "render the message without sending")
	sendCmd.MarkFlagRequired("template")
	rootCmd.AddCommand(sendCmd)
}
