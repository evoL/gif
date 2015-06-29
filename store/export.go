package store

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/evoL/gif/image"
	"io"
	"os"
	"time"
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
		gzipWriter := gzip.NewWriter(writer)
		defer gzipWriter.Close()

		tarWriter := tar.NewWriter(gzipWriter)
		defer tarWriter.Close()

		// Create a file with metadata
		metadataBuffer := new(bytes.Buffer)
		err = exportMetadata(exportedImages, metadataBuffer)
		if err != nil {
			return err
		}

		metadataHeader := &tar.Header{
			Name:    "gif.json",
			Mode:    0644,
			Size:    int64(metadataBuffer.Len()),
			ModTime: time.Now(),
		}

		if err = tarWriter.WriteHeader(metadataHeader); err != nil {
			return err
		}

		_, err = tarWriter.Write(metadataBuffer.Bytes())
		if err != nil {
			return err
		}

		// Add files
		for _, img := range images {
			fileInfo, err := os.Stat(s.PathFor(&img))
			if err != nil {
				return err
			}

			fileHeader, err := tar.FileInfoHeader(fileInfo, "")
			if err != nil {
				return err
			}

			if err = tarWriter.WriteHeader(fileHeader); err != nil {
				return err
			}

			file, err := os.Open(s.PathFor(&img))
			if err != nil {
				return err
			}
			defer file.Close()

			bufferedReader := bufio.NewReader(file)
			_, err = bufferedReader.WriteTo(tarWriter)
			if err != nil {
				return err
			}
		}

		if err = tarWriter.Close(); err != nil {
			return err
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

	encoder := json.NewEncoder(writer)

	if err := encoder.Encode(output); err != nil {
		return err
	}

	return nil
}
