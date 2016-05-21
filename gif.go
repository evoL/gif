package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/config"
	"github.com/evoL/gif/store"
	"github.com/evoL/gif/version"
	"net/url"
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
		cli.BoolFlag{
			Name:  "untagged",
			Usage: "Lists only images that have no tag.",
		},
		cli.BoolFlag{
			Name:  "local",
			Usage: "Lists only images that are local, that is not available remotely.",
		},
	}
	getFlags := append(
		typeFlags,
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "Gets all matching images.",
		},
		cli.StringFlag{
			Name:  "order, sort, s",
			Usage: "Specifies the order of images. Must be one of: random, newest, oldest.",
			Value: "random",
		},
	)
	removeFlags := append(
		typeFlags,
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "Removes all matching images.",
		},
		cli.BoolFlag{
			Name:  "really",
			Usage: "Doesn't ask for confirmation.",
		},
	)
	exportFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Target output file. Set to '-' for stdout.",
			Value: "-",
		},
		cli.BoolFlag{
			Name:  "bundle",
			Usage: "Export a bundle containing all images and metadata.",
		},
	}
	importFlags := []cli.Flag{
		cli.BoolFlag{
			Name:  "recursive, r",
			Usage: "When importing directories, do it recursively.",
		},
	}
	uploadFlags := append(
		typeFlags,
		cli.BoolFlag{
			Name:  "really",
			Usage: "Doesn't ask for confirmation.",
		},
	)
	recreateFlags := []cli.Flag{
		cli.BoolFlag{
			Name:  "really",
			Usage: "Doesn't ask for confirmation.",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Writes more information to the output.",
		},
	}

	app := cli.NewApp()
	app.Name = "gif"
	app.Usage = "a stupid gif manager"
	app.Author = "Rafał Hirsz"
	app.Email = "rafal@hirsz.co"
	app.Version = version.Version
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
			Name:   "export",
			Usage:  "Exports the database",
			Action: ExportCommand,
			Flags:  exportFlags,
		},
		{
			Name:   "import",
			Usage:  "Imports multiple images into the database",
			Action: ImportCommand,
			Flags:  importFlags,
		},
		{
			Name:   "list",
			Usage:  "Lists stored images",
			Action: ListCommand,
			Flags:  typeFlags,
		},
		{
			Name:   "path",
			Usage:  "Lists paths to images",
			Action: PathCommand,
			Flags:  getFlags,
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "Removes images",
			Action:  RemoveCommand,
			Flags:   removeFlags,
		},
		{
			Name:   "tag",
			Usage:  "Enables to change tags for images",
			Action: TagCommand,
			Flags:  typeFlags,
		},
		{
			Name:   "tags",
			Usage:  "Lists tags available in the database along with their image count",
			Action: TagsCommand,
		},
		{
			Name:   "url",
			Usage:  "Lists URLs of images",
			Action: UrlCommand,
			Flags:  getFlags,
		},
		{
			Name:   "upload",
			Usage:  "Uploads images to a server and saves the URLs for later use",
			Action: UploadCommand,
			Flags:  uploadFlags,
		},
		{
			Name:   "recreate",
			Usage:  "Recreates the database from scratch",
			Action: RecreateCommand,
			Flags:  recreateFlags,
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

type locationType int

const (
	invalidLocation locationType = iota
	fileLocation
	directoryLocation
	urlLocation
)

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
			filter = store.IdOrTagFilter{Id: arg}
		} else {
			filter = store.TagFilter{Tag: arg}
		}
	} else if c.Bool("untagged") {
		filter = store.UntaggedFilter{}
	} else {
		filter = store.NullFilter{}
	}

	if c.Bool("local") {
		filter = store.LocalFilter{Filter: filter}
	}

	return
}

func listFilter(c *cli.Context) store.Filter {
	return store.DateOrderer{
		Filter:    typeFilter(c),
		Direction: store.Descending,
	}
}

func orderAndLimit(input store.Filter, c *cli.Context) (filter store.Filter) {
	switch c.String("order") {
	case "random":
		filter = store.RandomOrderer{Filter: input}
	case "newest":
		filter = store.DateOrderer{Filter: input, Direction: store.Descending}
	case "oldest":
		filter = store.DateOrderer{Filter: input, Direction: store.Ascending}
	default:
		fmt.Println("Invalid order.")
		os.Exit(1)
	}

	if !c.Bool("all") {
		filter = store.Limiter{Filter: filter, Limit: 1}
	}

	return
}

func parseLocation(location string) (locationType, error) {
	if location == "" {
		return invalidLocation, errors.New("No location specified")
	}

	// Check for URL
	u, err := url.Parse(location)
	if err == nil {
		if u.Scheme == "http" || u.Scheme == "https" {
			return urlLocation, nil
		} else if u.Scheme != "" {
			return urlLocation, errors.New("Only HTTP and HTTPS URLs are supported")
		}
	}

	// Check for path
	fileInfo, err := os.Stat(location)
	if err == nil {
		if fileInfo.IsDir() {
			return directoryLocation, nil
		}
		return fileLocation, nil
	}

	return invalidLocation, errors.New("Invalid location")
}

func tempCommand(c *cli.Context) {
	s := getStore()
	defer s.Close()

	version, err := s.Version()
	if err != nil {
		fmt.Println("Error while fetching store version: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("Version", version)
}
