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
)

type locationType int

const (
	invalidLocation locationType = iota
	fileLocation
	directoryLocation
	urlLocation
)

func AddCommand(c *cli.Context) {
	location := c.Args().First()

	ltype, err := parseLocation(location)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	s := getStore()
	defer s.Close()

	writer := image.DefaultWriter()
	defer writer.Flush()

	var img *image.Image
	switch ltype {
	case urlLocation:
		// Check duplicates by URL
		img, _ = s.Find(store.UrlFilter{Url: location})
		if img != nil {
			io.WriteString(writer, "[exists]\t")
			img.PrintTo(writer)
			return
		}

		img, err = image.FromUrl(location)
	case fileLocation:
		img, err = image.FromFile(location)
	case directoryLocation:
		fmt.Println("Cannot add a directory.")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("Cannot load image: " + err.Error())
		os.Exit(1)
	}

	// Check if it already exists and show saved metadata
	var hit *image.Image
	hit, err = s.Get(img.Id)
	if hit != nil && err == nil {
		io.WriteString(writer, "[exists]\t")
		hit.PrintTo(writer)
		return
	}

	if err = s.Add(img); err != nil {
		fmt.Println("Cannot save image: " + err.Error())
		os.Exit(1)
	}

	err = TagInterface(s, img)
	if err != nil {
		fmt.Println("Cannot save tags: " + err.Error())
	}

	io.WriteString(writer, "[added]\t")
	img.PrintTo(writer)
	return
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
	fileInfo, err := os.Stat(location)
	if err == nil {
		if fileInfo.IsDir() {
			return directoryLocation, nil
		}
		return fileLocation, nil
	}

	return invalidLocation, errors.New("Invalid location")
}
