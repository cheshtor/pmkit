package pkg

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var configLoadLock sync.Mutex

var configs *viper.Viper

func GetConfig() *viper.Viper {
	configLoadLock.Lock()
	defer configLoadLock.Unlock()
	if configs == nil {
		config := viper.New()
		config.SetConfigFile("./configs/application.yaml")
		err := config.ReadInConfig()
		if err != nil {
			fmt.Printf("%s", err.Error())
			panic("系统配置文件异常")
		}
		configs = config
	}
	return configs
}
