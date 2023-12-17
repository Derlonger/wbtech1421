package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	ServerConfig ServerConfig
	Postgres     Postgres
	Stan         Stan
}

type ServerConfig struct {
	Host string
	Port string
}

type Postgres struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

type Stan struct {
	ClusterID   string
	ClientID    string
	ChannelName string
}

func init() {
	loadConfig()
}
func loadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func GetConfig() Config {
	return Config{
		ServerConfig: ServerConfig{
			Host: getEnv("HOST", "0.0.0.0"),
			Port: getEnv("PORT", "8080"),
		},
		Postgres: Postgres{
			Username: getEnv("POSTGRES_USERNAME", "admin"),
			Password: getEnv("POSTGRES_PASSWORD", "admin"),
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			DBName:   getEnv("POSTGRES_DBNAME", "wbtech_db"),
			SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
		},
		Stan: Stan{
			ClusterID:   getEnv("STAN_CLUSTER_ID", "my-cluster"),
			ClientID:    getEnv("STAN_CLIENT_ID", "client-1"),
			ChannelName: getEnv("STAN_CHANNEL_NAME", "events"),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
