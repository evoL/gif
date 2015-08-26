package upload

import (
	"errors"
	"github.com/evoL/gif/image"
)

type ImageUploader interface {
	Upload(img *image.Image) (bool, error)
}

func UploadImage(img *image.Image, uploader ImageUploader) (bool, error) {
	if !img.IsHydrated() {
		return false, errors.New("The image is not hydrated.")
	}

	return uploader.Upload(img)
}
