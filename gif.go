package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/config"
	"github.com/evoL/gif/image"
	"github.com/evoL/gif/store"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gif"
	app.Usage = "a stupid gif manager"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "default",
			Usage: "Path to the configuration file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "add",
			Usage:  "Adds an image",
			Action: AddCommand,
		},
		{
			Name:   "purge",
			Usage:  "Clears the image store. This cannot be reversed!",
			Action: PurgeCommand,
		},
		{
			Name:   "config",
			Usage:  "Prints the current configuration",
			Action: ConfigCommand,
		},
	}
	app.Before = func(c *cli.Context) (err error) {
		err = loadConfig(c.String("config"))
		return
	}

	app.Run(os.Args)
}

func loadConfig(arg string) (err error) {
	if arg == "default" {
		err = config.Default()
	} else {
		err = config.Load(arg)
	}

	if err != nil {
		fmt.Println("Error while loading the configuration file: " + err.Error())
	}
	return
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
	defer store.Close()

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

func PurgeCommand(c *cli.Context) {
	if store, err := store.Default(); err != nil {
		// It can't create the store, so why bother?
	} else {
		store.Purge()
	}
}

func ConfigCommand(c *cli.Context) {
	fmt.Printf("%s %v\n", "storePath", config.StorePath())
	fmt.Printf("%s %v\n", "db.dataSource", config.DbDataSource())
	fmt.Printf("%s %v\n", "db.driver", config.DbDriver())
}
