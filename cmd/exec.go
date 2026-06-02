package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var execTemplate string
var execSuccess string
var execVars []string

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run a command and send a Slack alert on failure",
	Long:  "Wraps a command — sends the failure template if it exits non-zero, optionally sends a success template on exit 0.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("exec: template=%s, command=%v\n", execTemplate, args)
		return nil
	},
}

func init() {
	execCmd.Flags().StringVarP(&execTemplate, "template", "t", "", "failure template name from squawk.yaml")
	execCmd.Flags().StringVar(&execSuccess, "success", "", "success template name (optional)")
	execCmd.Flags().StringSliceVar(&execVars, "var", nil, "additional template variables (key=value)")
	execCmd.MarkFlagRequired("template")
	rootCmd.AddCommand(execCmd)
}
