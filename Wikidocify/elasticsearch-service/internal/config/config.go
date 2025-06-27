// internal/config/config.go
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port         string        `json:"port"`
		ReadTimeout  time.Duration `json:"read_timeout"`
		WriteTimeout time.Duration `json:"write_timeout"`
		IdleTimeout  time.Duration `json:"idle_timeout"`
	} `json:"server"`

	Elasticsearch struct {
		URL      string `json:"url"`
		Username string `json:"username"`
		Password string `json:"password"`
		Index    string `json:"index"`
	} `json:"elasticsearch"`

	DocService struct {
		BaseURL string        `json:"base_url"`
		Timeout time.Duration `json:"timeout"`
		APIKey  string        `json:"api_key"`
	} `json:"doc_service"`

	Sync struct {
		BatchSize    int           `json:"batch_size"`
		SyncInterval time.Duration `json:"sync_interval"`
		EnableSync   bool          `json:"enable_sync"`
	} `json:"sync"`
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{}

	// Server config
	cfg.Server.Port = getEnv("PORT", "8080")
	cfg.Server.ReadTimeout = getDurationEnv("READ_TIMEOUT", 10*time.Second)
	cfg.Server.WriteTimeout = getDurationEnv("WRITE_TIMEOUT", 10*time.Second)
	cfg.Server.IdleTimeout = getDurationEnv("IDLE_TIMEOUT", 60*time.Second)

	// Elasticsearch config
	cfg.Elasticsearch.URL = getEnv("ELASTICSEARCH_URL", "http://localhost:9200")
	cfg.Elasticsearch.Username = getEnv("ELASTICSEARCH_USERNAME", "")
	cfg.Elasticsearch.Password = getEnv("ELASTICSEARCH_PASSWORD", "")
	cfg.Elasticsearch.Index = getEnv("ELASTICSEARCH_INDEX", "wikidocify_documents")

	// Doc service config
	cfg.DocService.BaseURL = getEnv("DOC_SERVICE_URL", "http://file-upload-service:8081")
	cfg.DocService.Timeout = getDurationEnv("DOC_SERVICE_TIMEOUT", 30*time.Second)
	cfg.DocService.APIKey = getEnv("DOC_SERVICE_API_KEY", "")

	// Sync config
	cfg.Sync.BatchSize = getIntEnv("SYNC_BATCH_SIZE", 100)
	cfg.Sync.SyncInterval = getDurationEnv("SYNC_INTERVAL", 5*time.Minute)
	cfg.Sync.EnableSync = getBoolEnv("ENABLE_SYNC", true)

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}