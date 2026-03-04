package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jonson/ha-ctl/internal/config"
	"github.com/jonson/ha-ctl/internal/haclient"
	"github.com/jonson/ha-ctl/internal/output"
	"github.com/jonson/ha-ctl/internal/util"
)

var stateCmd = &cobra.Command{
	Use:   "state <entity_id>",
	Short: "Get live state of an entity",
	Long:  "Fetch the current state of a specific entity directly from Home Assistant (bypasses cache).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entityID := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, util.FormatError("config", err.Error()))
			return err
		}

		client := haclient.New(cfg.HAURL, cfg.HAToken)
		entity, err := client.GetState(context.Background(), entityID)
		if err != nil {
			fmt.Fprintln(os.Stderr, util.FormatError("api", err.Error()))
			return err
		}

		out := output.EntityStateOutput{
			EntityID:    entity.EntityID,
			State:       entity.State,
			Attributes:  entity.Attributes,
			LastChanged: entity.LastChanged,
			LastUpdated: entity.LastUpdated,
		}

		formatter := output.New(GetFormat())
		result, err := formatter.Format(out)
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stateCmd)
}
