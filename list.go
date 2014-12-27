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

	filter := listFilter(c)

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

	filter := orderAndLimit(store.RemoteFilter{Filter: typeFilter(c)}, c)

	images, err := s.List(filter)
	if err != nil {
		fmt.Println("Error while fetching: " + err.Error())
		os.Exit(1)
	}

	for _, image := range images {
		fmt.Println(image.Url)
	}
}

func PathCommand(c *cli.Context) {
	s := getStore()
	defer s.Close()

	filter := orderAndLimit(typeFilter(c), c)

	images, err := s.List(filter)
	if err != nil {
		fmt.Println("Error while fetching: " + err.Error())
		os.Exit(1)
	}

	for _, image := range images {
		fmt.Println(s.PathFor(&image))
	}
}
