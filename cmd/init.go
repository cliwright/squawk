package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const squawkDir = ".squawk"
const defaultConfigFile = "squawk.yaml"

var defaultConfig = `templates:
  deploy-failed:
    channel: "#alerts"
    text: |
      ❌ *{{ .repo }}* failed on ` + "`{{ .branch }}`" + `
      {{ .input }}
      <{{ .run_url }}|View run>
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a .squawk directory with a starter template file",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := squawkDir
		configPath := filepath.Join(dir, defaultConfigFile)

		if _, err := os.Stat(dir); err == nil {
			return fmt.Errorf("%s already exists", dir)
		}

		if err := os.MkdirAll(dir, 0o700); err != nil {
			return fmt.Errorf("creating %s: %w", dir, err)
		}

		if err := os.WriteFile(configPath, []byte(defaultConfig), 0o600); err != nil {
			return fmt.Errorf("writing %s: %w", configPath, err)
		}

		fmt.Printf("created %s\n", configPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
