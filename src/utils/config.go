package utils

import (
	"os"
	"log"
	"fmt"
	"encoding/json"
	"sync"
	"time"
)

type baseSetting struct {
	Env string			`json:"env"`
	Services []string	`json:"services"`
	APIs struct {
		RPC struct {
			Active bool	`json:"active"`
			Port int	`json:"port"`
		}				`json:"rpc"`
	}					`json:"apis"`
}

type subsSetting struct {
	Db struct {
		Url string		`json:"url"`
		Name string		`json:"name"`
		Username string	`json:"username"`
		Password string	`json:"password"`
		MaxConn int		`json:"max_conn"`
	}					`json:"db"`
	Redis struct {
		Password string		`json:"password"`
		TimeFormat string	`json:"time_format"`
		ProcessPubKey string`json:"process_pub_key"`
		Clusters []struct {
			Name string		`json:"name"`
			Url string		`json:"url"`
		}					`json:"clusters"`
	}						`json:"redis"`
}

type coinSetting struct {
	Name string						`json:"name"`
	Url string						`json:"url"`
	Decimal int						`json:"decimal"`
	Stable int						`json:"stable"`
	Collect string					`json:"collect"`
	MinCollect float64				`json:"minCollect"`
	CollectInterval time.Duration	`json:"collectInterval"`
	TradePassword string			`json:"tradePassword"`
	UnlockDuration int				`json:"unlockDuration"`
	Withdraw string					`json:"withdraw"`
}

type msgsSetting struct {
	Logs struct {
		Debug bool					`json:"debug"`
	}								`json:"logs"`
	Level map[string]string			`json:"level"`
	Errors map[string]string		`json:"errors"`
	Warnings map[string]string		`json:"warnings"`
	Information map[string]string	`json:"information"`
	Debugs map[string]string		`json:"debugs"`
}

type cmdsSetting struct {
	Unknown string	`json:"unknown"`
	Help string		`json:"help"`
	Version string	`json:"version"`
}

type config struct {
	sync.Once
	base baseSetting
	subs subsSetting
	coin coinSetting
	msgs msgsSetting
	cmds cmdsSetting
}

var _config *config

func GetConfig() *config {
	if _config == nil {
		_config = new(config)
		_config.Once = sync.Once {}
		_config.Once.Do(func() {
			_config.create()
		})
	}
	return _config
}

func (cfg *config) create() error {
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
	if err = cfg.loadJson("command", &cfg.cmds); err != nil {
		panic(err)
	}
	return nil
}

func (cfg *config) loadJson(fileName string, data interface {}) error {
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

func (cfg *config) GetBaseSettings() baseSetting {
	return cfg.base
}

func (cfg *config) GetSubsSettings() subsSetting  {
	return cfg.subs
}

func (cfg *config) GetCoinSettings() coinSetting {
	return cfg.coin
}

func (cfg *config) GetMsgsSettings() msgsSetting {
	return cfg.msgs
}

func (cfg *config) GetCmdsSettings() cmdsSetting {
	return cfg.cmds
}