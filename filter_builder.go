package main

import (
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/store"
	"regexp"
	"strings"
)

func filterFlags() []cli.Flag {
	return []cli.Flag{
		// cli.BoolFlag{
		// 	Name:  "all, a",
		// 	Usage: "Takes all images into consideration. Do not use with --one.",
		// },
		// cli.BoolFlag{
		// 	Name:  "one, 1",
		// 	Usage: "Takes only one image into consideration. Do not use with --all.",
		// },
		cli.BoolFlag{
			Name:  "sample",
			Usage: "Gets a random picture.",
		},
		cli.StringFlag{
			Name:  "sort, order, s",
			Value: "[asc, desc]",
			Usage: "Enables ordering by addition date.",
		},
		cli.BoolFlag{
			Name:  "tag, t",
			Usage: "Enforces searching by tag.",
		},
	}
}

func buildFilter(c *cli.Context) (filter store.Filter, tagFilter store.Filter) {
	// Detect the type
	tagFilter = store.NullFilter{}

	if c.Args().Present() {
		arg := strings.Join(c.Args(), " ")

		if !c.Bool("tag") && regexp.MustCompile("^[0-9a-f]+$").MatchString(arg) {
			filter = store.IdFilter{Id: arg}
		} else {
			filter = store.NullFilter{}
			tagFilter = store.TagFilter{Tag: arg}
		}
	} else {
		filter = store.NullFilter{}
	}

	// Additional filters

	if c.Bool("sample") {
		filter = store.Limiter{
			Filter: store.RandomOrderer{Filter: filter},
			Limit:  1,
		}
	} else if c.String("sort") == "asc" {
		filter = store.DateOrderer{
			Filter:    filter,
			Direction: store.Ascending,
		}
	} else if c.String("sort") == "desc" {
		filter = store.DateOrderer{
			Filter:    filter,
			Direction: store.Descending,
		}
	}

	return
}
