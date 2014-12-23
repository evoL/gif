package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/image"
	"github.com/evoL/gif/store"
	"io"
	"net/url"
	"os"
	"text/tabwriter"
)

type locationType int

const (
	invalidLocation locationType = iota
	pathLocation
	urlLocation
)

func AddCommand(c *cli.Context) {
	location := c.Args().First()

	ltype, err := parseLocation(location)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	store, err := store.Default()
	if err != nil {
		fmt.Println("Cannot create store: " + err.Error())
		os.Exit(1)
	}
	defer store.Close()

	var img *image.Image
	switch ltype {
	case urlLocation:
		img, err = image.FromUrl(location)
	case pathLocation:
		img, err = image.FromFile(location)
	}
	if err != nil {
		fmt.Println("Cannot load image: " + err.Error())
		os.Exit(1)
	}

	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
	defer writer.Flush()

	if store.Contains(img) {
		io.WriteString(writer, "[exists]\t")
		img.PrintTo(writer)
		return
	}

	if err := store.Add(img); err != nil {
		fmt.Println("Cannot save image: " + err.Error())
		os.Exit(1)
	} else {
		io.WriteString(writer, "[added]\t")
		img.PrintTo(writer)
		return
	}
}

func parseLocation(location string) (locationType, error) {
	if location == "" {
		return invalidLocation, errors.New("No image specified")
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
	_, err = os.Stat(location)
	if err == nil {
		return pathLocation, nil
	}

	return invalidLocation, errors.New("Invalid location")
}
