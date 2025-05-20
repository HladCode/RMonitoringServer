package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Host                 string `json:"host"`
	Port                 string `json:"port"`
	BD_connect_parametrs string `json:"BD_URL"`
	KeyJWT               string `json:"KeyJWT"`
}

func MustRead(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening config file: ", err)
		return Config{}
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading config file: ", err)
		return Config{}
	}

	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatal("Error unmarshaling config file: ", err)
		return Config{}
	}

	return config
}
