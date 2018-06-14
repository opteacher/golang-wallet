package utils

import (
	"os"
	"log"
	"fmt"
	"encoding/json"
	"sync"
)

type baseSetting struct {
	Env string	`json:env`
}

type subsSetting struct {
	Db struct {
		Url string		`json:url`
		Name string		`json:name`
		Username string	`json:username`
		Password string	`json:password`
	}					`json:db`
}

type coinSetting struct {
	Name string	`json:name`
	Url string	`json:url`
	Decimal int	`json:decimal`
	Stable int	`json:stable`
}

type msgsSetting struct {
	Level map[string]string			`json:level`
	Errors map[string]string		`json:errors`
	Warnings map[string]string		`json:warnings`
	Information map[string]string	`json:information`
	Debugs map[string]string		`json:debugs`
}

type Config struct {
	sync.Once
	base baseSetting
	subs subsSetting
	coin coinSetting
	msgs msgsSetting
}

var _self *Config

func GetConfig() *Config {
	if _self == nil {
		_self = new(Config)
		_self.Once = sync.Once {}
		_self.Once.Do(func() {
			_self.create()
		})
	}
	return _self
}

func (cfg *Config) create() error {
	var err error
	if err = cfg.loadJson("settings", &cfg.base); err != nil {
		panic(err)
	}
	if err = cfg.loadJson(cfg.base.Env, &cfg.subs); err != nil {
		panic(err)
	}
	if err = cfg.loadJson("coin", &cfg.coin); err != nil {
		panic(err)
	}
	if err = cfg.loadJson("message", &cfg.msgs); err != nil {
		panic(err)
	}
	return nil
}

func (cfg *Config) loadJson(fileName string, data interface {}) error {
	file, err := os.Open(fmt.Sprintf("config/%s.json", fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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

func (cfg *Config) GetBaseSettings() baseSetting {
	return cfg.base
}

func (cfg *Config) GetSubsSettings() subsSetting  {
	return cfg.subs
}

func (cfg *Config) GetCoinSettings() coinSetting {
	return cfg.coin
}

func (cfg *Config) GetMsgsSettings() msgsSetting {
	return cfg.msgs
}