package frame

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)


var Config globalConfig


type globalConfig struct {
  Server struct {
    Host string
    Port int
  }
  Mongo struct {
    Host string
    Port int
    UserName string
    Password string
    DBName string
    Params []string
  }
  Redis struct {
    Host string
    Port int
    DB int
  }
  Email struct {
    Type string
    Username string
    Auth string
  }
}


func init()  {
  wd, err := os.Getwd()
  if err != nil {
    panic(err)
  }
  configPath := filepath.Join(wd, "config.yml")
  viper.SetConfigFile(configPath)
  if err := viper.ReadInConfig(); err != nil {
    panic(err)
  }
  Config = globalConfig{}
  if err := viper.Unmarshal(&Config); err != nil {
    panic(err)
  }

  // log.Println(Config)
}
