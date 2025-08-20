package config

import (
	"flag"
	"os"
	"strconv"
)

// Config содержит настройки приложения
type Config struct {
	Port int
}

// Load загружает конфигурацию из переменных окружения и флагов
func Load() *Config {
	var port int
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	// Проверяем переменную окружения
	if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		if parsedPort, err := strconv.Atoi(envPort); err == nil {
			port = parsedPort
		}
	}

	return &Config{
		Port: port,
	}
}
