package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/image"
	"github.com/evoL/gif/store"
	"os"
	"regexp"
	"strings"
)

func TagInterface(s *store.Store, img *image.Image) error {
	fmt.Println("How would you like to tag this image? (Enter new tags separated by commas)")
	if len(img.Tags) > 0 {
		fmt.Println("Current tags:", strings.Join(img.Tags, ", "))
	}

	fmt.Print("=> ")
	bio := bufio.NewReader(os.Stdin)
	tagList, _, _ := bio.ReadLine()

	return s.UpdateTags(img, regexp.MustCompile(`,\s*`).Split(string(tagList[:]), -1))
}

func TagCommand(c *cli.Context) {
	s := getStore()
	defer s.Close()

	if !c.Args().Present() {
		cli.ShowCommandHelp(c, "tag")
		os.Exit(1)
	}

	filter := listFilter(c)

	images, err := s.List(filter)
	if err != nil {
		fmt.Println("Error while fetching: " + err.Error())
		os.Exit(1)
	}

	for _, image := range images {
		image.Print()

		err = TagInterface(s, &image)
		if err != nil {
			fmt.Println("Error while updating tags: " + err.Error())
			os.Exit(1)
		}
	}
}
