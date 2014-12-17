package image

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
)

type Image struct {
	Id   string
	Url  string
	Tags []string
	Data []byte
}

func FromUrl(url string) (*Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	image := fromData(data)
	image.Url = url
	return image, nil
}

func fromData(data []byte) *Image {
	image := &Image{Data: data}
	image.Id = image.GenerateId()
	return image
}

func (image *Image) GenerateId() string {
	h := sha1.New()
	h.Write(image.Data)
	return hex.EncodeToString(h.Sum(nil))
}
