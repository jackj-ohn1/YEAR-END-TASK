package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	name string
}

func (cfg *Config) Init(name string) {
	cfg.name = name
	cfg.Run()
}

func (cfg *Config) Run() {
	if cfg.name == "" {
		viper.AddConfigPath("./config")
		viper.SetConfigFile("config")
	} else {
		viper.SetConfigFile(cfg.name)
	}
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("读取配置文件失败:", err)
	}
}
