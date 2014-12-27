package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/image"
	"github.com/evoL/gif/store"
	"os"
)

func ListCommand(c *cli.Context) {
	s := getStore()
	defer s.Close()

	var filter store.Filter
	if c.Bool("untagged") {
		filter = store.UntaggedFilter{}
	} else {
		filter = typeFilter(c)
	}
	filter = store.DateOrderer{
		Filter:    filter,
		Direction: store.Descending,
	}

	images, err := s.List(filter)
	if err != nil {
		fmt.Println("Error while fetching: " + err.Error())
		os.Exit(1)
	}

	fmt.Printf("%v images\n", len(images))

	image.PrintAll(images)
}

func UrlCommand(c *cli.Context) {
	s := getStore()
	defer s.Close()

	var filter store.Filter = store.RemoteFilter{Filter: typeFilter(c)}

	switch c.String("order") {
	case "random":
		filter = store.RandomOrderer{Filter: filter}
	case "newest":
		filter = store.DateOrderer{Filter: filter, Direction: store.Descending}
	case "oldest":
		filter = store.DateOrderer{Filter: filter, Direction: store.Ascending}
	default:
		fmt.Println("Invalid order.")
		os.Exit(1)
	}

	if !c.Bool("all") {
		filter = store.Limiter{Filter: filter, Limit: 1}
	}

	images, err := s.List(filter)
	if err != nil {
		fmt.Println("Error while fetching: " + err.Error())
		os.Exit(1)
	}

	for _, image := range images {
		fmt.Println(image.Url)
	}
}
