package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort    string `mapstructure:"PORT"`
	GetchipsURL   string `mapstructure:"GETCHIPS_URL"`
	GetchipsToken string `mapstructure:"GETCHIPS_TOKEN"`
	RedisAddr     string `mapstructure:"REDIS_ADDRESS"`
	RabbitMQURL   string `mapstructure:"BROKER_URL"`
	ChunkSize     int    `mapstructure:"CHUNK_SIZE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.SetDefault("PORT", "5004")
	viper.SetDefault("GETCHIPS_URL", "https://api.client-service.getchips.ru/client/api/gh/v1/search/partnumber")
	viper.SetDefault("CHUNK_SIZE", 50)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
