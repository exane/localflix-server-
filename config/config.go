package config

import (
	"encoding/json"
	"io/ioutil"
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
	}
	TMDb struct {
		ApiKey string
	}
}

var config *Config

func LoadConfig() *Config {
	if config != nil {
		return config
	}
	config = &Config{}
	js, _ := ioutil.ReadFile("./config.json")
	json.Unmarshal(js, config)
	return config
}
