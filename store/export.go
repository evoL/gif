package store

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/evoL/gif/image"
	"io"
	"os"
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
		zipWriter := zip.NewWriter(writer)
		defer zipWriter.Close()

		// Create a file with metadata
		metadataFile, err := zipWriter.Create("gif.json")
		if err != nil {
			return err
		}

		err = exportMetadata(exportedImages, metadataFile)
		if err != nil {
			return err
		}

		// Add files
		for _, img := range images {
			fileName := fmt.Sprintf("%v.gif", img.Id)
			imageFile, err := zipWriter.Create(fileName)
			if err != nil {
				return err
			}

			file, err := os.Open(s.PathFor(&img))
			if err != nil {
				return err
			}
			defer file.Close()

			bufferedReader := bufio.NewReader(file)
			_, err = bufferedReader.WriteTo(imageFile)
			if err != nil {
				return err
			}
		}

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
