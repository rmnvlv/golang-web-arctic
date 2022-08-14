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

func (c *Config) LoadEnv() error {
	var ok bool

	if captcha, ok := os.LookupEnv("ENABLE_CAPTCHA"); !ok {
		if captcha == "true" {
			c.Captcha.Enable = true

			c.Captcha.SecretKey, ok = os.LookupEnv("HCAPTCHA_SECRET_KEY")
			if !ok {
				return fmt.Errorf("HCAPTCHA_SECRET_KEY not set")
			}
			c.Captcha.SiteKey, ok = os.LookupEnv("HCAPTCHA_SITE_KEY")
			if !ok {
				return fmt.Errorf("HCAPTCHA_SITE_KEY not set")
			}

		}
	}

	c.SMTP.User, ok = os.LookupEnv("SMTP_USER")
	if !ok {
		return fmt.Errorf("SMTP_USER not set")
	}

	c.SMTP.Password, ok = os.LookupEnv("SMTP_PASSWORD")
	if !ok {
		return fmt.Errorf("SMTP_PASSWORD not set")
	}

	c.SMTP.Host, ok = os.LookupEnv("SMTP_HOST")
	if !ok {
		return fmt.Errorf("SMTP_HOST not set")
	}

	if SMTPPort, ok := os.LookupEnv("SMTP_PORT"); !ok {
		return fmt.Errorf("SMTP_PORT not set")
	} else {
		var err error
		c.SMTP.Port, err = strconv.Atoi(SMTPPort)
		if err != nil {
			return fmt.Errorf("SMTP_PORT not set: %w", err)

		}
	}

	if c.DiskPath == "" {
		c.DiskPath, ok = os.LookupEnv("DISK_PATH")
		if !ok {
			return fmt.Errorf("DISK_PATH not set")
		}
	}

	if c.DatabaseURL == "" {
		c.DatabaseURL, ok = os.LookupEnv("DATABASE_URL")
		if !ok {
			return fmt.Errorf("DATABASE_URL not set")
		}
	}

	c.AdminPassword, ok = os.LookupEnv("ADMIN_PASSWORD")
	if !ok {
		return fmt.Errorf("ADMIN_PASSWORD not set")
	}

	return nil
}

func (c *Config) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")

	buf := bytes.NewBuffer(b)

	return buf.String()
}
