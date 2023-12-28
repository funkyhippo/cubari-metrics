package internal

import (
	"log"
	"os"
)

type Config interface {
	GetRemoteWriteURL() string
	GetRemoteWriteUsername() string
	GetRemoteWritePassword() string
	GetPort() string
}

type EnvConfig struct {
	remoteWriteURL      string
	remoteWriteUsername string
	remoteWritePassword string
	port                string
}

func NewConfigFromEnv() EnvConfig {
	remoteWriteURL := os.Getenv("REMOTE_WRITE_URL")
	if remoteWriteURL == "" {
		log.Fatal("REMOTE_WRITE_URL was empty")
	}
	remoteWriteUsername := os.Getenv("REMOTE_WRITE_USERNAME")
	if remoteWriteUsername == "" {
		log.Fatal("REMOTE_WRITE_USERNAME was empty")
	}
	remoteWritePassword := os.Getenv("REMOTE_WRITE_PASSWORD")
	if remoteWritePassword == "" {
		log.Fatal("REMOTE_WRITE_PASSWORD was empty")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT was empty")
	}

	return EnvConfig{
		remoteWriteURL:      remoteWriteURL,
		remoteWriteUsername: remoteWriteUsername,
		remoteWritePassword: remoteWritePassword,
		port:                port,
	}
}

func (c EnvConfig) GetRemoteWriteURL() string {
	return c.remoteWriteURL
}

func (c EnvConfig) GetRemoteWriteUsername() string {
	return c.remoteWriteUsername
}

func (c EnvConfig) GetRemoteWritePassword() string {
	return c.remoteWritePassword
}

func (c EnvConfig) GetPort() string {
	return c.port
}
