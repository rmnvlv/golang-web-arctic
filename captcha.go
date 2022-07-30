package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const hCaptchaAPIURL = "https://hcaptcha.com/siteverify"

var (
	hCaptcha = HCaptcha{
		SiteKey:   os.Getenv("HCAPTCHA_SITE_KEY"),
		SecretKey: os.Getenv("HCAPTCHA_SECRET_KEY"),
	}
	ErrCaptchaEmpty = errors.New("captcha is empty")
)

type HCaptcha struct {
	SiteKey   string `mapstructure:"HCAPTCHA_SITE_KEY"`
	SecretKey string `mapstructure:"HCAPTCHA_SECRET_KEY"`
}

// Response is the hcaptcha JSON response.
type Response struct {
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
	Success     bool     `json:"success"`
	Credit      bool     `json:"credit,omitempty"`
}

func verifyCaptcha(captcha string) (bool, error) {
	if captcha == "" {
		return false, ErrCaptchaEmpty
	}

	form := url.Values{}
	form.Add("secret", hCaptcha.SecretKey)
	form.Add("response", captcha)
	form.Add("sitekey", hCaptcha.SiteKey)

	resp, err := http.DefaultClient.PostForm(hCaptchaAPIURL, form)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return false, err
	}

	if !response.Success {
		return false, fmt.Errorf("hCaptcha: %v", response.ErrorCodes)
	}

	return true, nil
}
