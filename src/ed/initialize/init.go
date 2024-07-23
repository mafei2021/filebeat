package initialize

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"sync"

	"ed/gen"
	"github.com/rs/zerolog/log"
)

func Init(confPath string) {
	initConf(confPath)
	initLog()
}

var (
	once = sync.Once{}
)

func initConf(configPath string) {
	var configContent []byte
	var err error
	if configPath != "" {
		configContent, err = os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Error reading configuration file:", err)
			return
		}
	} else {
		configContent, err = gen.Asset("resources/application.yaml")
		if err != nil {
			log.Fatal().Msgf("读取配置文件出现错误 ! %s", err)
		}
	}
	log.Debug().Msgf("configContent:", string(configContent))
	// configs 用于初始化系统的配置
	once.Do(func() {
		viper.SetConfigType("yaml")
		err := viper.ReadConfig(bytes.NewReader(configContent))
		if err != nil {
			log.Fatal().Msgf("读取配置文件出现错误! %s", err.Error())
		}
	})
}
