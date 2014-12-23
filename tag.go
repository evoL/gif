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
	store, err := store.Default()
	if err != nil {
		fmt.Println("Cannot create store: " + err.Error())
		os.Exit(1)
	}
	defer store.Close()
}
