package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/image"
	"github.com/evoL/gif/store"
	"io"
	"os"
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
			img.PrintTo(writer, false)
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

	AddInterface(s, writer, img, true)
	return
}

func AddInterface(s *store.Store, writer image.FlushableWriter, img *image.Image, replaceTags bool) {
	// Check if it already exists and show metadata
	var hit *image.Image
	hit, err := s.Get(img.Id)
	if hit != nil && err == nil {
		// Check if the URL can be updated
		if img.Url != "" && img.Url != hit.Url {
			err = s.UpdateUrl(hit, img.Url)
			if err == nil {
				io.WriteString(writer, "[update]\t")
			} else {
				io.WriteString(writer, "[unchgd]\t")
			}
		} else {
			io.WriteString(writer, "[exists]\t")
		}
		hit.PrintTo(writer, true)
		return
	}

	if err = s.WriteImage(img); err != nil {
		fmt.Println("Cannot save image: " + err.Error())
		os.Exit(1)
	}

	if err = s.Add(img); err != nil {
		fmt.Println("Cannot save image: " + err.Error())
		os.Exit(1)
	}

	if replaceTags {
		err = TagInterface(s, img)
	} else {
		err = s.UpdateTags(img, img.Tags)
	}
	if err != nil {
		fmt.Println("Cannot save tags: " + err.Error())
	}

	io.WriteString(writer, "[added]\t")
	img.PrintTo(writer, true)

	return
}
