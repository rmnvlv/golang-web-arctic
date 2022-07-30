package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var resultJson struct {
	Href string `json:"href"`
}

// func main() {
// 	err := uploadArticleYandex("./uploadedFiles/myFile.docx", "Articles/newFile.docx")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

func initEnv() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

//Need local path of the file and path of ydisc like docs/pasport.jpg (docs - folder)
//Two remote paths : Articles/name.doc Thusiss/name.doc
func uploadArticleYandex(localPath io.Reader, remotePath string) error {
	initEnv()

	var yandexURL, _ = os.LookupEnv("YANDEX_URL")
	var token, _ = os.LookupEnv("YANDEX_TOKEN")

	// fmt.Println(token, yandexURL)

	//Read file
	// data, err := os.Open(localPath)
	// fmt.Println("End of open data")
	// if err != nil {
	// 	return err
	// }

	//Get url to upload file
	href, err := getRemoteUrl(remotePath, yandexURL, token)
	// fmt.Println("End of getRemoteUrl")
	if err != nil {
		return err
	}

	// defer data.Close()

	//Upload file with uploading url
	request, err := http.NewRequest("PUT", href, localPath)
	// fmt.Println("End of getting request")
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("OAuth %s", token))

	client := http.Client{}
	response, err := client.Do(request)
	// fmt.Println("End of request")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}

func apiRequest(path, method, yandexURL, token string) (*http.Response, error) {
	client := http.Client{}
	url := fmt.Sprintf("%s/%s", yandexURL, path)
	request, _ := http.NewRequest(method, url, nil)
	request.Header.Add("Authorization", fmt.Sprintf("OAuth %s", token))
	return client.Do(request)
}

//Get url of reader/writer of Ydisk example:https://uploader34g.disk.yandex.net:443...
func getRemoteUrl(path, yandexURL, token string) (string, error) {
	// overwrite - ?
	response, err := apiRequest(fmt.Sprintf("resources/upload?path=%s&overwrite=true", path), "GET", yandexURL, token)
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
