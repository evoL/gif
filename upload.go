package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/evoL/gif/config"
	"github.com/evoL/gif/image"
	"github.com/evoL/gif/store"
	"github.com/evoL/gif/upload"
	"os"
	"regexp"
	"strings"
)

func UploadCommand(c *cli.Context) {
	s := getStore()
	defer s.Close()

	uploader, err := makeUploader()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	filter := listFilter(c)
	images, err := s.List(filter)
	if err != nil {
		fmt.Println("Error while fetching: " + err.Error())
		os.Exit(1)
	}

	fmt.Printf("%v images\n", len(images))

	writer := image.DefaultWriter()
	defer writer.Flush()

	for _, img := range images {
		img.PrintTo(writer, false)
	}
	fmt.Println()

	if c.Bool("really") {
		uploadImages(images, uploader, s)
	} else {
		fmt.Print("Do you really want to upload those images? [y/n] ")

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			fmt.Println("IO Error: " + scanner.Err().Error())
			os.Exit(1)
		}

		response := strings.TrimSpace(scanner.Text())
		if regexp.MustCompile(`\A[yY]\z`).MatchString(response) {
			uploadImages(images, uploader, s)
		}
	}
}

func makeUploader() (upload.ImageUploader, error) {
	switch config.Global.Upload.Provider {
	case "imgur":
		{
			clientId, ok := config.Global.Upload.Credentials["ClientId"]
			if ok {
				return upload.ImgurUploader{ClientId: clientId}, nil
			} else {
				return nil, errors.New("A Client ID for the Imgur API is required for uploading.")
			}
		}
	default:
		return nil, errors.New("Uploading is disabled until you set up a provider. Available providers: imgur.")
	}
}

func uploadImages(images []image.Image, uploader upload.ImageUploader, s *store.Store) {
	count := len(images)
	for i, img := range images {
		fmt.Printf("Uploading image %v/%v… ", i+1, count)
		s.Hydrate(&img)

		_, err := upload.UploadImage(&img, uploader)
		if err != nil {
			fmt.Println("✘")
			fmt.Println("Error: " + err.Error())
			continue
		}

		err = s.UpdateUrl(&img, img.Url)
		if err != nil {
			fmt.Println("✘")
			fmt.Println("Error: " + err.Error())
			continue
		}

		fmt.Println("✔")
	}
}
