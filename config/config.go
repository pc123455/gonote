package config

import (
	"io/ioutil"
	"log"
	"os"
	"yaml"
)

type Config struct {
	Net struct {
		Bind string `ymal:"bind"`
		Port int    `ymal:"port"`
	}
	Log struct {
		File  string `ymal:"file"`
		Level string `ymal:"level"`
	}
}

func readConfig(fileanme string) ([]byte, error) {
	f, err := os.Open(fileanme)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

func ParseConfigFromFile(filename string) *Config {
	data, err := readConfig(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
		//os.Exit(1)
	}
	conf := new(Config)
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		log.Fatalf("error: %v", err)
		//os.Exit(1)
	}
	return conf
}
