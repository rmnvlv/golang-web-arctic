package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	YandexDiskAPIURL = "https://cloud-api.yandex.net/v1/disk"
	ArticlesFolder   = "Articles"
	AbstractsFolder  = "Abstracts"
)

var (
	YandexClientId    = os.Getenv("YANDEX_CLIENT_ID")
	YandexCallbackURL = os.Getenv("YANDEX_CALLBACK_URL")
	YandexOAuthToken  = os.Getenv("YANDEX_OAUTH_TOKEN")
)

type UploadURLResponse struct {
	OperationId string `json:"operation_id"`
	URL         string `json:"href"`
	Method      string `json:"method"`
}

func saveToYandexDisk(file io.Reader, remotePath string) error {
	logger := log.Default()

	header := http.Header{}
	header.Add("Authorization", "OAuth "+YandexOAuthToken)

	client := http.Client{}

	getUploadURLRequest, err := http.NewRequest("GET",
		fmt.Sprintf("%s/resources/upload?path=%s&overwrite=true", YandexDiskAPIURL, remotePath), nil)
	if err != nil {
		logger.Println(err)
		return err
	}

	getUploadURLRequest.Header = header

	response, err := client.Do(getUploadURLRequest)
	if err != nil {
		logger.Println(err)
		return err
	}

	defer response.Body.Close()

	var getUploadURLResponse UploadURLResponse
	if err := json.NewDecoder(response.Body).Decode(&getUploadURLResponse); err != nil {
		logger.Println(err)
		fmt.Println(getUploadURLResponse)
		return err
	}

	logger.Println(getUploadURLResponse)
	//PUT
	uploadFileRequest, err := http.NewRequest(getUploadURLResponse.Method, getUploadURLResponse.URL, file)
	if err != nil {
		fmt.Println(uploadFileRequest)
		return err
	}

	uploadFileRequest.Header = header

	response, err = client.Do(uploadFileRequest)
	if err != nil {
		fmt.Println(response)
		return err
	}

	defer response.Body.Close()

	return nil
}

// func main() {
// 	f, err := os.Open("go.mod")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = saveToYandexDisk(f, "Test/go.mod")
// 	fmt.Println(err)
// }
