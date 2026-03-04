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

var entitiesDomain string
var entitiesState string
var entitiesRefresh bool

var entitiesCmd = &cobra.Command{
	Use:   "entities",
	Short: "List cached entities",
	Long:  "List entities from the local cache. Use --refresh to force a cache update.",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		if c == nil || cache.IsStale(c) || entitiesRefresh {
			client := haclient.New(cfg.HAURL, cfg.HAToken)
			c, err = cache.Refresh(context.Background(), client, cfg.CacheTTL)
			if err != nil {
				fmt.Fprintln(os.Stderr, util.FormatError("api", err.Error()))
				return err
			}
		}

		var entities []cache.CacheEntity
		if entitiesDomain != "" {
			entities = cache.FilterByDomain(c, entitiesDomain)
		} else {
			entities = cache.EntityList(c)
		}

		if entitiesState != "" {
			entities = cache.FilterByState(entities, entitiesState)
		}

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
	entitiesCmd.Flags().StringVar(&entitiesDomain, "domain", "", "Filter by domain (e.g. light, switch)")
	entitiesCmd.Flags().StringVar(&entitiesState, "state", "", "Filter by state (e.g. on, off, unavailable)")
	entitiesCmd.Flags().BoolVar(&entitiesRefresh, "refresh", false, "Force cache refresh")
	rootCmd.AddCommand(entitiesCmd)
}
