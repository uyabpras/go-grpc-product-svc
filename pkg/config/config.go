package config

import "github.com/spf13/viper"

type Config struct {
	Port   string `mapstructure:"PORT"`
	DB_url string `mapstructure:"DB_URL"`
}

func LoadConfig() (Config Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&Config)
	return
}
