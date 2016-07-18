package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

type Config struct {
	Database struct {
		Root     string
		Password string
		Db       string
	}
	Server struct {
		URL  string
		Port string
	}
	Fileserver struct {
		URL           string
		RootDirectory string
		Port          string
		FetchOut      string
	}
	TMDb struct {
		ApiKey string
	}
}

type ConfigData struct {
	Test        *Config
	Development *Config
}

var config *Config

func LoadConfig() *Config {
	println("env: " + os.Getenv("ENV"))
	_, filename, _, _ := runtime.Caller(1)
	filepath := path.Join(path.Dir(filename), "../config.json")
	if config != nil {
		return config
	}
	config = &Config{}
	configData := &ConfigData{}
	js, _ := ioutil.ReadFile(filepath)
	json.Unmarshal(js, configData)

	if os.Getenv("ENV") == "test" {
		config = configData.Test
	} else if os.Getenv("ENV") == "development" {
		config = configData.Development
	} else {
		config = configData.Development
	}

	return config
}
