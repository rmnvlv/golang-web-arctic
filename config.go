package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPAddress     string
	HTTPAddressUnix string
	AdminPassword   string
	Domain          string
	DatabaseURL     string
	Captcha         HCaptchaConfig
	SMTP            SMTPConfig
	DiskPath        string
}

type HCaptchaConfig struct {
	Enable    bool
	SecretKey string
	SiteKey   string
}

type SMTPConfig struct {
	User     string
	Password string
	Host     string
	Port     int
}

func NewConfig() (Config, error) {
	c := Config{}
	var ok bool

	a, aok := os.LookupEnv("HTTP_ADDRESS")
	au, auok := os.LookupEnv("HTTP_ADDRESS_UNIX")

	if !aok && !auok {
		return Config{}, fmt.Errorf("HTTP_ADDRESS or HTTP_ADDRESS_UNIX must be set")
	}

	if aok {
		c.HTTPAddress = a
	}

	if auok {
		c.HTTPAddressUnix = au
	}

	if captcha, ok := os.LookupEnv("ENABLE_CAPTCHA"); !ok {
		if captcha == "true" {
			c.Captcha.Enable = true

			c.Captcha.SecretKey, ok = os.LookupEnv("HCAPTCHA_SECRET_KEY")
			if !ok {
				return Config{}, fmt.Errorf("HCAPTCHA_SECRET_KEY not set")
			}
			c.Captcha.SiteKey, ok = os.LookupEnv("HCAPTCHA_SITE_KEY")
			if !ok {
				return Config{}, fmt.Errorf("HCAPTCHA_SITE_KEY not set")
			}

		}
	}

	c.SMTP.User, ok = os.LookupEnv("SMTP_USER")
	if !ok {
		return Config{}, fmt.Errorf("SMTP_USER not set")
	}

	c.SMTP.Password, ok = os.LookupEnv("SMTP_PASSWORD")
	if !ok {
		return Config{}, fmt.Errorf("SMTP_PASSWORD not set")
	}

	c.SMTP.Host, ok = os.LookupEnv("SMTP_HOST")
	if !ok {
		return Config{}, fmt.Errorf("SMTP_HOST not set")
	}

	if SMTPPort, ok := os.LookupEnv("SMTP_PORT"); !ok {
		return Config{}, fmt.Errorf("SMTP_PORT not set")
	} else {
		var err error
		c.SMTP.Port, err = strconv.Atoi(SMTPPort)
		if err != nil {
			return Config{}, fmt.Errorf("SMTP_PORT not set: %w", err)

		}
	}

	c.DiskPath, ok = os.LookupEnv("DISK_PATH")
	if !ok {
		return Config{}, fmt.Errorf("DISK_PATH not set")
	}

	c.DatabaseURL, ok = os.LookupEnv("DATABASE_URL")
	if !ok {
		return Config{}, fmt.Errorf("DATABASE_URL not set")
	}

	c.AdminPassword, ok = os.LookupEnv("ADMIN_PASSWORD")
	if !ok {
		return Config{}, fmt.Errorf("ADMIN_PASSWORD not set")
	}

	return c, nil

}

func (c *Config) String() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(b)

	return buf.String()
}
