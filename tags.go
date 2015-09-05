package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/store"
	"os"
	"text/tabwriter"
)

func TagsCommand(c *cli.Context) {
	s := getStore()
	defer s.Close()

	var filter store.Filter
	if c.Args().Present() {
		filter = store.TagPrefixFilter{Prefix: c.Args().First()}
	} else {
		filter = store.NullFilter{}
	}

	tags, err := s.ListTags(filter)
	if err != nil {
		fmt.Println("Error while fetching: " + err.Error())
		os.Exit(1)
	}

	fmt.Printf("%v tags\n", len(tags))

	if len(tags) > 0 {
		writer := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
		defer writer.Flush()

		fmt.Fprintln(writer, "TAG\tIMAGE COUNT")

		for _, tag := range tags {
			fmt.Fprintf(writer, "%s\t%v\n", tag.Tag, tag.Count)
		}
	}
}
