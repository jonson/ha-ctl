package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/jonson/ha-ctl/internal/cache"
	"github.com/jonson/ha-ctl/internal/config"
	"github.com/jonson/ha-ctl/internal/haclient"
	"github.com/jonson/ha-ctl/internal/output"
	"github.com/jonson/ha-ctl/internal/util"
)

var findDomain string

var findCmd = &cobra.Command{
	Use:   "find <query>",
	Short: "Search entities by name",
	Long:  "Case-insensitive substring search against entity_id and friendly_name. Much faster than dumping full context.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, util.FormatError("config", err.Error()))
			return err
		}

		c, err := cache.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, util.FormatError("cache", err.Error()))
			return err
		}

		if c == nil || cache.IsStale(c) {
			client := haclient.New(cfg.HAURL, cfg.HAToken)
			c, err = cache.Refresh(context.Background(), client, cfg.CacheTTL)
			if err != nil {
				fmt.Fprintln(os.Stderr, util.FormatError("api", err.Error()))
				return err
			}
		}

		entities := cache.Search(c, query, findDomain)

		sort.Slice(entities, func(i, j int) bool {
			return entities[i].EntityID < entities[j].EntityID
		})

		items := make([]output.EntityItem, len(entities))
		for i, e := range entities {
			items[i] = output.EntityItem{
				EntityID:      e.EntityID,
				State:         e.State,
				FriendlyName:  e.FriendlyName,
				Domain:        e.Domain,
				KeyAttributes: e.KeyAttributes,
			}
		}

		out := output.EntityListOutput{
			Entities: items,
			Count:    len(items),
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
	findCmd.Flags().StringVar(&findDomain, "domain", "", "Filter by domain (e.g. light, switch)")
	rootCmd.AddCommand(findCmd)
}
