package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/config"
	"github.com/evoL/gif/store"
	"os"
	"regexp"
	"strings"
)

func main() {
	typeFlags := []cli.Flag{
		cli.BoolFlag{
			Name:  "tag, t",
			Usage: "Enforces searching by tag.",
		},
	}
	listFlags := append(
		typeFlags,
		cli.BoolFlag{
			Name:  "untagged",
			Usage: "Lists only images that have no tag.",
		},
	)
	getFlags := append(
		typeFlags,
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "Gets all matching images",
		},
		cli.StringFlag{
			Name:  "order, sort, s",
			Usage: "Specifies the order of images. Must be one of: random, newest, oldest.",
			Value: "random",
		},
	)

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
			Name:   "config",
			Usage:  "Prints the current configuration",
			Action: ConfigCommand,
		},
		{
			Name:   "list",
			Usage:  "Lists stored images",
			Action: ListCommand,
			Flags:  listFlags,
		},
		{
			Name:   "url",
			Usage:  "Lists URLs of images",
			Action: UrlCommand,
			Flags:  getFlags,
		},
	}
	app.Before = func(c *cli.Context) (err error) {
		err = loadConfig(c.String("config"))
		return
	}

	app.Run(os.Args)
}

func ConfigCommand(c *cli.Context) {
	config.Global.Print()
}

///////////////////////////////////////////////////////////////////////////////

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

func getStore() *store.Store {
	s, err := store.Default()
	if err != nil {
		fmt.Println("Cannot create store: " + err.Error())
		os.Exit(1)
	}
	return s
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
