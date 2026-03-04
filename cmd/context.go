package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jonson/ha-ctl/internal/cache"
	"github.com/jonson/ha-ctl/internal/config"
	"github.com/jonson/ha-ctl/internal/haclient"
	"github.com/jonson/ha-ctl/internal/output"
	"github.com/jonson/ha-ctl/internal/util"
)

// controllableDomains are the domains an agent would typically act on.
var controllableDomains = map[string]bool{
	"light":         true,
	"switch":        true,
	"climate":       true,
	"media_player":  true,
	"cover":         true,
	"fan":           true,
	"lock":          true,
	"scene":         true,
	"automation":    true,
	"input_boolean": true,
}

var contextFull bool
var contextDomains string

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Generate home summary for LLM context",
	Long:  "Generate a compact summary of the home state, grouped by domain. Optimized for LLM token efficiency.",
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

		if c == nil || cache.IsStale(c) {
			client := haclient.New(cfg.HAURL, cfg.HAToken)
			c, err = cache.Refresh(context.Background(), client, cfg.CacheTTL)
			if err != nil {
				fmt.Fprintln(os.Stderr, util.FormatError("api", err.Error()))
				return err
			}
		}

		// Build the set of domains to expand
		expandDomains := controllableDomains
		if contextFull {
			// Expand all domains
			expandDomains = nil // nil means expand everything
		} else if contextDomains != "" {
			// Expand only the user-specified domains
			expandDomains = make(map[string]bool)
			for _, d := range strings.Split(contextDomains, ",") {
				d = strings.TrimSpace(d)
				if d != "" {
					expandDomains[d] = true
				}
			}
		}

		// Group entities by domain
		domainMap := make(map[string][]cache.CacheEntity)
		for _, e := range c.Entities {
			domainMap[e.Domain] = append(domainMap[e.Domain], e)
		}

		// Sort domain names
		domainNames := make([]string, 0, len(domainMap))
		for d := range domainMap {
			domainNames = append(domainNames, d)
		}
		sort.Strings(domainNames)

		var expanded []output.DomainSummary
		other := make(map[string]int)
		totalEntities := 0

		for _, d := range domainNames {
			entities := domainMap[d]
			totalEntities += len(entities)

			shouldExpand := expandDomains == nil || expandDomains[d]
			if !shouldExpand {
				other[d] = len(entities)
				continue
			}

			sort.Slice(entities, func(i, j int) bool {
				return entities[i].EntityID < entities[j].EntityID
			})

			briefs := make([]output.EntityBrief, len(entities))
			for i, e := range entities {
				briefs[i] = output.EntityBrief{
					EntityID:      e.EntityID,
					Name:          e.FriendlyName,
					State:         e.State,
					KeyAttributes: e.KeyAttributes,
				}
			}

			expanded = append(expanded, output.DomainSummary{
				Domain:   d,
				Count:    len(entities),
				Entities: briefs,
			})
		}

		out := output.ContextOutput{
			Summary:      fmt.Sprintf("Home has %d entities across %d domains", totalEntities, len(domainNames)),
			Controllable: expanded,
		}
		if len(other) > 0 {
			out.Other = other
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
	contextCmd.Flags().BoolVar(&contextFull, "full", false, "Expand all domains (restore full output)")
	contextCmd.Flags().StringVar(&contextDomains, "domain", "", "Comma-separated domains to expand (e.g. light,sensor)")
	rootCmd.AddCommand(contextCmd)
}
