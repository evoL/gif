package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/image"
	"github.com/evoL/gif/store"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gif"
	app.Usage = "a stupid gif manager"
	app.Commands = []cli.Command{
		{
			Name:   "add",
			Usage:  "Adds an image",
			Action: AddCommand,
		},
	}

	app.Run(os.Args)
}

func AddCommand(c *cli.Context) {
	url := c.Args().First()
	if url == "" {
		fmt.Println("add: No image specified")
		os.Exit(1)
	}

	store, err := store.Default()
	if err != nil {
		fmt.Println("Cannot create store: " + err.Error())
		os.Exit(1)
	}

	image, err := image.FromUrl(url)
	if err != nil {
		fmt.Println("Cannot load image: " + err.Error())
		os.Exit(1)
	}

	if store.Contains(image) {
		fmt.Println("Image already exists: " + store.PathFor(image))
		return
	}

	if err := store.Save(image); err != nil {
		fmt.Println("Cannot save image: " + err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Saved image: " + store.PathFor(image))
	}
}
