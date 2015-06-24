package store

import (
	"bytes"
	"encoding/json"
	"github.com/evoL/gif/image"
	"io"
	// "zip"
)

type exportedImage struct {
	Id   string
	Url  string
	Tags []string
}

type metadata struct {
	Creator string
	Images  []exportedImage
}

func (s *Store) Export(writer io.Writer, filter Filter, exportFiles bool) error {
	images, err := s.List(filter)
	if err != nil {
		return err
	}

	exportedImages := prepareImages(images)

	if exportFiles {
		buffer := new(bytes.Buffer)

		err = exportMetadata(exportedImages, buffer)
		if err != nil {
			return err
		}

		// TODO: make a zip
	} else {
		err = exportMetadata(exportedImages, writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func prepareImages(images []image.Image) (exportedImages []exportedImage) {
	exportedImages = make([]exportedImage, len(images))

	for i, img := range images {
		exportedImages[i].Id = img.Id
		exportedImages[i].Url = img.Url
		exportedImages[i].Tags = img.Tags
	}

	return
}

func exportMetadata(images []exportedImage, writer io.Writer) error {
	output := metadata{
		Creator: "gif",
		Images:  images,
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)
	return err
}
