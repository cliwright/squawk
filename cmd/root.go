package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "squawk",
	Version: version,
	Short:   "A dead-simple CLI for sending Slack alerts from CI pipelines",
	Long:    "squawk sends Slack alerts when things break in CI. One red light, not ten green lights.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
