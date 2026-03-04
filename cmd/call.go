package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jonson/ha-ctl/internal/config"
	"github.com/jonson/ha-ctl/internal/haclient"
	"github.com/jonson/ha-ctl/internal/output"
	"github.com/jonson/ha-ctl/internal/util"
)

var callEntity string
var callData string

var callCmd = &cobra.Command{
	Use:   "call <domain> <service>",
	Short: "Call a Home Assistant service",
	Long:  "Call a service on Home Assistant. Example: ha-ctl call light turn_on --entity light.kitchen",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		service := args[1]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, util.FormatError("config", err.Error()))
			return err
		}

		data := make(map[string]any)
		if callEntity != "" {
			data["entity_id"] = callEntity
		}

		if callData != "" {
			var extra map[string]any
			if err := json.Unmarshal([]byte(callData), &extra); err != nil {
				fmt.Fprintln(os.Stderr, util.FormatError("validation", "invalid --data JSON: "+err.Error()))
				return err
			}
			for k, v := range extra {
				data[k] = v
			}
		}

		client := haclient.New(cfg.HAURL, cfg.HAToken)
		_, err = client.CallService(context.Background(), domain, service, data)
		if err != nil {
			fmt.Fprintln(os.Stderr, util.FormatError("api", err.Error()))
			return err
		}

		out := output.ServiceCallOutput{
			Success:  true,
			Domain:   domain,
			Service:  service,
			EntityID: callEntity,
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
	callCmd.Flags().StringVar(&callEntity, "entity", "", "Entity ID to target")
	callCmd.Flags().StringVar(&callData, "data", "", "Additional service data as JSON")
	rootCmd.AddCommand(callCmd)
}
