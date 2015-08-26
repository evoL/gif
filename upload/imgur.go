package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/evoL/gif/image"
	"github.com/evoL/gif/version"
	"mime/multipart"
	"net/http"
)

const endpoint = "https://api.imgur.com/3/"

var extensions map[string]string = map[string]string{
	"image/gif":  ".gif",
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/tiff": ".tif",
}

type ImgurUploader struct {
	ClientId string
}

type imgurBasicModel struct {
	Data       map[string]interface{} `json:"data"`
	Success    bool                   `json:"success"`
	StatusCode int                    `json:"status"`
}

func (u ImgurUploader) Upload(img *image.Image) (bool, error) {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	writer.WriteField("type", "file")

	// TODO: Implement proper file type detection
	mimeType := http.DetectContentType(img.Data)
	extension, ok := extensions[mimeType]
	if !ok {
		extension = ".gif"
	}

	part, err := writer.CreateFormFile("image", img.Id+extension)
	if err != nil {
		return false, err
	}

	_, err = part.Write(img.Data)
	if err != nil {
		return false, err
	}

	err = writer.Close()
	if err != nil {
		return false, err
	}

	request, err := http.NewRequest("POST", endpoint+"image", body)
	if err != nil {
		return false, err
	}

	request.Header.Add("Authorization", "Client-ID "+u.ClientId)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("User-Agent", "gif/"+version.Version)

	// testfile, _ := os.Create("request.txt")
	// defer testfile.Close()
	// request.Write(testfile)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return false, err
	}

	if response.StatusCode >= 500 {
		return false, errors.New("imgur: Internal server error")
	}

	responseJson := &imgurBasicModel{}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(responseJson)
	if err != nil {
		return false, err
	}

	if responseJson.Success {
		img.Url = responseJson.Data["link"].(string)
		return true, nil
	} else {
		return false, errors.New(responseJson.Data["error"].(string))
	}
}
