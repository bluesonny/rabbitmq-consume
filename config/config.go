package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Mq struct {
	User     string
	Password string
	Host     string
	Vhost    string
	Queue    string
}

type Redis struct {
	Password string
	Host     string
}

type Database struct {
	Driver   string
	Address  string
	Database string
	User     string
	Password string
}

type Configuration struct {
	Mq    Mq
	Db    Database
	Redis Redis
}

var config *Configuration
var once sync.Once

// 通过单例模式初始化全局配置
func LoadConfig() *Configuration {
	once.Do(func() {
		file, err := os.Open("config.json")
		if err != nil {
			log.Fatalln("Cannot open config file", err)
		}
		decoder := json.NewDecoder(file)
		config = &Configuration{}
		err = decoder.Decode(config)
		if err != nil {
			log.Fatalln("Cannot get configuration from file", err)
		}
	})
	return config
}
