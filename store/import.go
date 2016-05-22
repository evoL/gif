package store

import (
	"encoding/json"
	"github.com/evoL/gif/image"
	"io"
)

func ParseMetadata(reader io.Reader) ([]ExportedImage, error) {
	var input ExportFormat
	jsonDecoder := json.NewDecoder(reader)
	if err := jsonDecoder.Decode(&input); err != nil {
		return nil, err
	}

	return input.Images, nil
}

func (s *Store) ImportMetadata(images []ExportedImage) (err error) {
	for _, exported := range images {
		img := image.Image{
			Id:   exported.Id,
			Url:  exported.Url,
			Tags: exported.Tags,
		}

		if exported.AddedAt != "" {
			if err = img.SetAddedAtFromString(exported.AddedAt); err != nil {
				return
			}
		}

		if err = s.Add(&img); err != nil {
			return
		}

		if err = s.UpdateTags(&img, exported.Tags); err != nil {
			return
		}
	}

	return
}
