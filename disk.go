package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	YandexDiskAPIURL = "https://cloud-api.yandex.net/v1/disk"
	ArticlesFolder   = "Articles"
	AbstractsFolder  = "Abstracts"
)

var (
	YandexOAuthToken = Cfg["YANDEX_OAUTH_TOKEN"]
)

type UploadURLResponse struct {
	OperationId string `json:"operation_id"`
	URL         string `json:"href"`
	Method      string `json:"method"`
}

func (a *App) saveToYandexDisk(file io.Reader, remotePath string) error {
	header := http.Header{}
	header.Add("Authorization", "OAuth "+YandexOAuthToken)

	client := http.Client{}

	getUploadURLRequest, err := http.NewRequest("GET",
		fmt.Sprintf("%s/resources/upload?path=%s&overwrite=true", YandexDiskAPIURL, remotePath), nil)
	if err != nil {
		a.log.Error(err)
		return err
	}

	getUploadURLRequest.Header = header

	response, err := client.Do(getUploadURLRequest)
	if err != nil {
		a.log.Error(err)
		return err
	}

	defer response.Body.Close()

	var getUploadURLResponse UploadURLResponse
	if err := json.NewDecoder(response.Body).Decode(&getUploadURLResponse); err != nil {
		a.log.Error(err)
		a.log.Debug(getUploadURLResponse)
		return err
	}

	a.log.Debug(getUploadURLResponse)

	uploadFileRequest, err := http.NewRequest(getUploadURLResponse.Method, getUploadURLResponse.URL, file)
	if err != nil {
		a.log.Error(uploadFileRequest)
		return err
	}

	uploadFileRequest.Header = header

	response, err = client.Do(uploadFileRequest)
	if err != nil {
		a.log.Error(err)
		return err
	}
	a.log.Debug(response)

	defer response.Body.Close()

	return nil
}
