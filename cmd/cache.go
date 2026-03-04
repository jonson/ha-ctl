package cmd

import (
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Cache management commands",
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}
