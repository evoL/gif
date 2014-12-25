package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/image"
	"github.com/evoL/gif/store"
	"os"
)

func ListCommand(c *cli.Context) {
	s, err := store.Default()
	if err != nil {
		fmt.Println("Cannot create store: " + err.Error())
		os.Exit(1)
	}
	defer s.Close()

	filter := store.DateOrderer{
		Filter:    store.NullFilter{},
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
