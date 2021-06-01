package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var Conf Config

type Config struct {
	Port     string `json:"port"`
	DBport   string `json:"db_port"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBname   string `json:"dbname"`
}

func ParseConf(configPath string) error {
	jsonFile, err := os.Open(configPath)
	defer jsonFile.Close()
	if err != nil {
		return err
	}
	confJson, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(confJson, &Conf)
	return err
}
