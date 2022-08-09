package main

import (
	"io"
	"os"

	"github.com/google/uuid"
)

type Disk interface {
	Save(file io.Reader, fileName string) error
}

type OsDisk struct {
	Path string
}

func NewOsDisk(path string) (*OsDisk, error) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return nil, err
	}

	return &OsDisk{Path: path}, nil
}

func (d *OsDisk) Save(file io.Reader, fileName string) error {
	f, err := os.Create(d.Path + "/" + fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) saveToDisk(file io.Reader, extention string) error {
	fileName := uuid.New().String()

	return a.disk.Save(file, fileName+"."+extention)
}

// const (
// 	YandexDiskAPIURL = "https://cloud-api.yandex.net/v1/disk"
// 	ArticlesFolder   = "Articles"
// 	AbstractsFolder  = "Abstracts"
// )

// var (
// 	YandexClientId    = os.Getenv("YANDEX_CLIENT_ID")
// 	YandexCallbackURL = os.Getenv("YANDEX_CALLBACK_URL")
// 	YandexOAuthToken  = os.Getenv("YANDEX_OAUTH_TOKEN")
// )

// type UploadURLResponse struct {
// 	OperationId string `json:"operation_id"`
// 	URL         string `json:"href"`
// 	Method      string `json:"method"`
// }

// func (a *App) saveToYandexDisk(file io.Reader, remotePath string) error {
// 	header := http.Header{}
// 	header.Add("Authorization", "OAuth "+YandexOAuthToken)

// 	client := http.Client{}

// 	getUploadURLRequest, err := http.NewRequest("GET",
// 		fmt.Sprintf("%s/resources/upload?path=%s&overwrite=true", YandexDiskAPIURL, remotePath), nil)
// 	if err != nil {
// 		a.log.Error(err)
// 		return err
// 	}

// 	getUploadURLRequest.Header = header

// 	response, err := client.Do(getUploadURLRequest)
// 	if err != nil {
// 		a.log.Error(err)
// 		return err
// 	}

// 	defer response.Body.Close()

// 	var getUploadURLResponse UploadURLResponse
// 	if err := json.NewDecoder(response.Body).Decode(&getUploadURLResponse); err != nil {
// 		a.log.Error(err)
// 		a.log.Debug(getUploadURLResponse)
// 		return err
// 	}

// 	a.log.Debug(getUploadURLResponse)

// 	uploadFileRequest, err := http.NewRequest(getUploadURLResponse.Method, getUploadURLResponse.URL, file)
// 	if err != nil {
// 		a.log.Error(uploadFileRequest)
// 		return err
// 	}

// 	uploadFileRequest.Header = header

// 	response, err = client.Do(uploadFileRequest)
// 	if err != nil {
// 		a.log.Error(err)
// 		return err
// 	}
// 	a.log.Debug(response)

// 	defer response.Body.Close()

// 	return nil
// }
