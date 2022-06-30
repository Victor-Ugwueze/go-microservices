package config

import "github.com/spf13/viper"

type Config struct {
	Port          string `mapstructure:"PORT"`
	AccountSvcUrl string `mapstructure:"ACCOUNT_SVC_URL"`
	ProductSvcUrl string `mapstructure:"PRODUCT_SVC_URL"`
	OrderSvcUrl   string `mapstructure:"ORDER_SVC_URL"`
}


func LoadConfig() (c Config, err error) {

	viper.AddConfigPath("./")
	viper.SetConfigName(".env.dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()


	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return 
}