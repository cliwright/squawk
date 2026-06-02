package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sendTemplate string
var sendVars []string

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a Slack alert using a named template",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("sending template: %s\n", sendTemplate)
		return nil
	},
}

func init() {
	sendCmd.Flags().StringVarP(&sendTemplate, "template", "t", "", "template name from squawk.yaml")
	sendCmd.Flags().StringSliceVar(&sendVars, "var", nil, "additional template variables (key=value)")
	sendCmd.MarkFlagRequired("template")
	rootCmd.AddCommand(sendCmd)
}
