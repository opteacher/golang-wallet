package utils

import (
	"os"
	"log"
	"fmt"
	"encoding/json"
)

type Config struct {
	isCreated bool
	setting struct {
		Db struct {
			Url string		`json:url`
			Name string		`json:name`
			Username string	`json:username`
			Password string	`json:password`
		}					`json:db`
	}
}

func (cfg *Config) Create() error {
	var settings struct {
		Env string	`json:env`
	}
	var err error
	if err = cfg.loadJson("settings", &settings); err != nil {
		panic(err)
	}
	if err = cfg.loadJson(settings.Env, &cfg.setting); err != nil {
		panic(err)
	}
	cfg.isCreated = true
	return nil
}

func (cfg *Config) loadJson(fileName string, data interface {}) error {
	file, err := os.Open(fmt.Sprintf("../config/%s.json", fileName))
	if err != nil {
		log.Println(err)
		return err
	}

	chunks := make([]byte, 1024, 1024)
	bufData := []byte {}
	totalLen := 0
	for {
		n, err := file.Read(chunks)
		if n == 0 { break }
		totalLen += n
		if err != nil {
			log.Fatal(err)
			return err
		}
		bufData = append(bufData, chunks...)
	}

	if err = json.Unmarshal(bufData[:totalLen], &data); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (cfg *Config) IsCreate() bool {
	return cfg.isCreated
}

func NewConfig() *Config {
	return new(Config)
}