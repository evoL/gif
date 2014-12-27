package main

import (
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/store"
	"regexp"
	"strings"
)

func typeFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "tag, t",
			Usage: "Enforces searching by tag.",
		},
	}
}

func typeFilter(c *cli.Context) (filter store.Filter) {
	if c.Args().Present() {
		arg := strings.Join(c.Args(), " ")

		if !c.Bool("tag") && regexp.MustCompile("^[0-9a-f]+$").MatchString(arg) {
			filter = store.IdFilter{Id: arg}
		} else {
			filter = store.TagFilter{Tag: arg}
		}
	} else {
		filter = store.NullFilter{}
	}

	return
}
