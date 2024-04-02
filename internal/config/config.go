package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GetConfig()
}

var (
	c *Config
)

type Config struct {
	Mail    Mail
	WebPort WebPort
}

type WebPort struct {
	Port string
}

type Mail struct {
	Host       string
	Username   string
	Password   string
	Encryption string
	Port       string
}

func GetConfig() *Config {
	if c == nil {
		host := os.Getenv("MAIL_HOST")
		if host == "" {
			panic("MAIL_HOST is not set")
		}

		username := os.Getenv("MAIL_USERNAME")
		if username == "" {
			panic("MAIL_USERNAME is not set")
		}

		password := os.Getenv("MAIL_PASSWORD")
		if password == "" {
			panic("MAIL_PASSWORD is not set")
		}

		encryption := os.Getenv("MAIL_ENCRYPTION")
		if encryption == "" {
			panic("MAIL_ENCRYPTION is not set")
		}

		mailPort := os.Getenv("MAIL_PORT")
		if mailPort == "" {
			panic("MAIL_PORT error")
		}

		port := os.Getenv("WEB_PORT")
		if port == "" {
			panic("WebPort is not set")
		}

		c = &Config{
			Mail: Mail{
				Host:       host,
				Username:   username,
				Password:   password,
				Encryption: encryption,
				Port:       mailPort,
			},
			WebPort: WebPort{
				Port: port,
			},
		}
	}

	return c
}
