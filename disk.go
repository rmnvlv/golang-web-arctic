package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var yandexURL = "https://cloud-api.yandex.net/v1/disk"
var token = "AQAAAABjUKKxAAhFwmPyMwXHlU-3kzwFGX67I9o"

var resultJson struct {
	Href string `json:"href"`
}

// func main() {
// 	err := uploadArticleYandex("./uploadedFiles/myFile.docx", "Articles/newFile.docx")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

//Need local path of the file and path of ydisc like docs/pasport.jpg (docs - folder)
//Two remote paths : Articles/name.doc Thusiss/name.doc
func uploadArticleYandex(localPath, remotePath string) error {

	//Read file
	data, err := os.Open(localPath)
	fmt.Println("End of open data")
	if err != nil {
		return err
	}

	//Get url to upload file
	href, err := getRemoteUrl(remotePath)
	fmt.Println("End of getRemoteUrl")
	if err != nil {
		return err
	}

	defer data.Close()

	//Upload file with uploading url
	request, err := http.NewRequest("PUT", href, data)
	fmt.Println("End of getting request")
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("OAuth %s", token))

	client := http.Client{}
	response, err := client.Do(request)
	fmt.Println("End of request")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func apiRequest(path, method string) (*http.Response, error) {
	client := http.Client{}
	url := fmt.Sprintf("%s/%s", yandexURL, path)
	request, _ := http.NewRequest(method, url, nil)
	// fmt.Println(request)
	request.Header.Add("Authorization", fmt.Sprintf("OAuth %s", token))
	return client.Do(request)
}

//Get url of reader/writer of Ydisk example:https://uploader34g.disk.yandex.net:443...
func getRemoteUrl(path string) (string, error) {
	// overwrite - ?
	response, err := apiRequest(fmt.Sprintf("resources/upload?path=%s&overwrite=true", path), "GET")
	if err != nil {
		return "", err
	}

	err = json.NewDecoder(response.Body).Decode(&resultJson)
	if err != nil {
		return "", err
	}

	// fmt.Println(resultJson.Href)

	return resultJson.Href, nil
}
