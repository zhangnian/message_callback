package g

import (
	"encoding/json"
	"log"
	"message_callback/utils"
	"sync"
)

type GlobalCondfig struct {
	AMQP     string `json:"amqp"`
	Retry    int    `json:"retry"`
	RetryTTL int    `json:"retry_ttl"`
}

type QueueConfig struct {
	Name        string `json:"name"`
	CallbackUrl string `json:"callback_url"`
}

type Config struct {
	Global GlobalCondfig `json:"global"`
	Queues []QueueConfig `json:"queues"`
}

var (
	config *Config
	lock   sync.RWMutex
)

func GetConfig() *Config {
	lock.RLock()
	defer lock.RUnlock()

	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !utils.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not exist.")
	}

	configContent, err := utils.ReadFile(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c Config
	err = json.Unmarshal(configContent, &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c
}
