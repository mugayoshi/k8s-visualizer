// pkg/config/config.go
package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	Server     ServerConfig
	Kubernetes KubernetesConfig
	WebSocket  WebSocketConfig
	Logging    LoggingConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Mode         string // "debug" or "release"
}

// KubernetesConfig holds Kubernetes client configuration
type KubernetesConfig struct {
	KubeConfig string
	InCluster  bool
	Namespace  string // Default namespace
	QPS        float32
	Burst      int
}

// WebSocketConfig holds WebSocket configuration
type WebSocketConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
	PingPeriod      time.Duration
	PongWait        time.Duration
	WriteWait       time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string // "debug", "info", "warn", "error"
	Format     string // "json" or "text"
	OutputPath string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			Mode:         getEnv("GIN_MODE", "debug"),
		},
		Kubernetes: KubernetesConfig{
			KubeConfig: getEnv("KUBECONFIG", ""),
			InCluster:  getBoolEnv("K8S_IN_CLUSTER", false),
			Namespace:  getEnv("K8S_DEFAULT_NAMESPACE", "default"),
			QPS:        getFloat32Env("K8S_QPS", 50.0),
			Burst:      getIntEnv("K8S_BURST", 100),
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:  getIntEnv("WS_READ_BUFFER_SIZE", 1024),
			WriteBufferSize: getIntEnv("WS_WRITE_BUFFER_SIZE", 1024),
			PingPeriod:      getDurationEnv("WS_PING_PERIOD", 30*time.Second),
			PongWait:        getDurationEnv("WS_PONG_WAIT", 60*time.Second),
			WriteWait:       getDurationEnv("WS_WRITE_WAIT", 10*time.Second),
		},
		Logging: LoggingConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "text"),
			OutputPath: getEnv("LOG_OUTPUT_PATH", "stdout"),
		},
	}
}

// Helper functions to get environment variables with defaults

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
		return defaultValue
	}
	return intValue
}

func getFloat32Env(key string, defaultValue float32) float32 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	floatValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		log.Printf("Invalid float value for %s: %s, using default: %f", key, value, defaultValue)
		return defaultValue
	}
	return float32(floatValue)
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Invalid boolean value for %s: %s, using default: %t", key, value, defaultValue)
		return defaultValue
	}
	return boolValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Invalid duration value for %s: %s, using default: %s", key, value, defaultValue)
		return defaultValue
	}
	return duration
}

// Print logs the configuration (sanitized)
func (c *Config) Print() {
	log.Println("=== Configuration ===")
	log.Printf("Server Host: %s", c.Server.Host)
	log.Printf("Server Port: %s", c.Server.Port)
	log.Printf("Server Mode: %s", c.Server.Mode)
	log.Printf("K8s In-Cluster: %t", c.Kubernetes.InCluster)
	log.Printf("K8s Default Namespace: %s", c.Kubernetes.Namespace)
	log.Printf("Log Level: %s", c.Logging.Level)
	log.Println("====================")
}
