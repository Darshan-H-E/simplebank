package util

import "github.com/spf13/viper"

// stores all configuration values
// values are read by viper from config file or env variables
type Config struct {
  DBDriver string `mapstructure:"DB_DRIVER"`
  DBSource string `mapstructure:"DB_SOURCE"`
  ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
  viper.AddConfigPath(path)
  viper.SetConfigName("app")
  viper.SetConfigType("env")

  viper.AutomaticEnv()

  err = viper.ReadInConfig()
  if err != nil {
    return
  }

  viper.Unmarshal(&config)
  return
}
