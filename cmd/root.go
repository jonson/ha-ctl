package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jonson/ha-ctl/internal/util"
)

var formatFlag string

var rootCmd = &cobra.Command{
	Use:   "ha-ctl",
	Short: "Home Assistant CLI for LLM agent skills",
	Long:  "ha-ctl wraps the Home Assistant REST API for use as an OpenClaw/Moltbot AgentSkill.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&formatFlag, "format", "json", "Output format: json or text")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, util.FormatError("cli", err.Error()))
		os.Exit(1)
	}
}

// SetVersion sets the version string for the version command.
func SetVersion(v string) {
	rootCmd.Version = v
}

// GetFormat returns the current output format flag value.
func GetFormat() string {
	return formatFlag
}
