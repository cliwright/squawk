package cmd

import (
	"fmt"

	"github.com/cliwright/squawk/config"
	"github.com/cliwright/squawk/slack"
	"github.com/spf13/cobra"
)

var testTemplate string
var testVars []string

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Render and send a template to Slack for verification",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(squawkDir)
		if err != nil {
			return err
		}

		tmpl, ok := cfg.Templates[testTemplate]
		if !ok {
			return fmt.Errorf("template %q not found", testTemplate)
		}

		vars := buildVars(testVars)

		text, err := tmpl.Render(vars)
		if err != nil {
			return err
		}

		fmt.Println("--- rendered message ---")
		fmt.Println(text)
		fmt.Println("--- sending to", tmpl.Channel, "---")

		token, err := slack.Token()
		if err != nil {
			return err
		}

		return slack.Send(token, slack.Message{
			Channel: tmpl.Channel,
			Text:    text,
		})
	},
}

func init() {
	testCmd.Flags().StringVarP(&testTemplate, "template", "t", "", "template name from .squawk/")
	testCmd.Flags().StringSliceVar(&testVars, "var", nil, "additional template variables (key=value)")
	testCmd.MarkFlagRequired("template")
	rootCmd.AddCommand(testCmd)
}
