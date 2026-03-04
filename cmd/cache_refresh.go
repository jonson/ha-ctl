package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/jonson/ha-ctl/internal/cache"
	"github.com/jonson/ha-ctl/internal/config"
	"github.com/jonson/ha-ctl/internal/haclient"
	"github.com/jonson/ha-ctl/internal/output"
	"github.com/jonson/ha-ctl/internal/util"
)

var cacheRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Force a full cache refresh from Home Assistant",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, util.FormatError("config", err.Error()))
			return err
		}

		client := haclient.New(cfg.HAURL, cfg.HAToken)

		fmt.Fprintf(os.Stderr, "Fetching entities from %s...\n", cfg.HAURL)

		c, err := cache.Refresh(context.Background(), client, cfg.CacheTTL)
		if err != nil {
			fmt.Fprintln(os.Stderr, util.FormatError("api", err.Error()))
			return err
		}

		out := output.CacheStatsOutput{
			Success:     true,
			EntityCount: len(c.Entities),
			Timestamp:   c.Timestamp.Format(time.RFC3339),
			CachePath:   cache.Path(),
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
	cacheCmd.AddCommand(cacheRefreshCmd)
}
