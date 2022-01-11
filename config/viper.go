package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var ViperConfig Configuration

func init() {
	runtimeViper := viper.New()
	runtimeViper.AddConfigPath(".")
	runtimeViper.SetConfigName("config")
	runtimeViper.SetConfigType("json")
	err := runtimeViper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	runtimeViper.Unmarshal(&ViperConfig)

	// 监听配置文件变更
	runtimeViper.WatchConfig()
	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		runtimeViper.Unmarshal(&ViperConfig)
		//ViperConfig.LocaleBundle.MustLoadMessageFile(ViperConfig.App.Locale + "/active." + ViperConfig.App.Language + ".json")
	})
}
