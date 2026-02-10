package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort     string
	GetchipsURL    string
	GetchipsToken  string
	EfindURL       string
	EfindToken     string
	PromelecURL    string
	PromelecLogin  string
	PromelecPass   string
	RedisAddr      string
	RabbitMQURL    string
	ChunkSize      int
	WorkerPoolSize int
}

func LoadConfig() Config {
	cfg := Config{
		ServerPort:     getEnv("PORT", "5004"),
		GetchipsURL:    getEnv("GETCHIPS_URL", "https://api.client-service.getchips.ru/client/api/gh/v1/search/partnumber"),
		GetchipsToken:  getEnv("GETCHIPS_TOKEN", ""),
		EfindURL:       getEnv("EFIND_URL", "https://efind.ru/api/search"),
		EfindToken:     getEnv("EFIND_TOKEN", ""),
		PromelecURL:    getEnv("PROMELEC_URL", "https://aaa.na4u.ru/rpc"),
		PromelecLogin:  getEnv("PROMELEC_LOGIN", ""),
		PromelecPass:   getEnv("PROMELEC_PASS", ""),
		RedisAddr:      getEnv("REDIS_ADDR", ""),
		RabbitMQURL:    getEnv("RABBITMQ_URL", ""),
		ChunkSize:      getEnvAsInt("CHUNK_SIZE", 50),
		WorkerPoolSize: getEnvAsInt("WORKER_POOL_SIZE", 20),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
