package image

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
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
	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("%s %s", response.Proto, response.Status)
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

func FromFile(path string) (*Image, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	image := fromData(data)
	return image, nil
}

func fromData(data []byte) *Image {
	image := &Image{Data: data}
	image.Id = image.generateId()
	return image
}

func (image *Image) generateId() string {
	h := sha1.New()
	h.Write(image.Data)
	return hex.EncodeToString(h.Sum(nil))
}
