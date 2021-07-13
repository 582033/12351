package libs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func configFieldGet(config *viper.Viper, taskName string, field string) string {
	return fmt.Sprintf("%s", config.Get(taskName+"."+field))
}

func LoadConf(confFile string) Config {
	config := viper.New()
	config.SetConfigFile(confFile)
	config.SetConfigType("ini") //设置文件的类型
	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		log.Println(err)
	}
	configList := make(Config)
	for c := range config.AllSettings() {
		JSESSIONID := configFieldGet(config, c, "JSESSIONID")
		service_id := configFieldGet(config, c, "service_id")
		token := configFieldGet(config, c, "token")
		receive_id := configFieldGet(config, c, "receive_id")
		phone := configFieldGet(config, c, "token")
		firm_id := configFieldGet(config, c, "firm_id")

		configList[c] = UserInfo{
			JSESSIONID: JSESSIONID,
			service_id: service_id,
			token:      token,
			receive_id: receive_id,
			phone:      phone,
			firm_id:    firm_id,
		}
	}

	return configList
}
