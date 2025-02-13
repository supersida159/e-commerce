package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

const (
	ProductionEnv = "production"

	DatabaseTimeout    = 5 * time.Second
	ProductCachingTime = 1 * time.Minute
)

var AuthIgnoreMethods = []string{
	"/user.UserService/Login",
	"/user.UserService/Register",
}

type Schema struct {
	Environment   string `yaml:"ENVIRONMENT"`
	HttpPort      int    `yaml:"HTTP_PORT"`
	GrpcPort      int    `yaml:"GRPC_PORT"`
	AuthSecret    string `yaml:"AUTH_SECRET"`
	DatabaseURI   string `yaml:"DATABASE_URI"`
	RedisURI      string `yaml:"REDIS_URI"`
	RedisPassword string `yaml:"REDIS_PASSWORD"`
	RedisDB       int    `yaml:"REDIS_DB"`
	ExpireTime    int    `yaml:"EXPIRY_TIME"`
	Kafka         struct {
		Broker               []string `yaml:"KAFKA_BROKERS"`
		Retry                int      `yaml:"KAFKA_RETRY"`
		ConsumerOffsetReset  string   `yaml:"KAFKA_CONSUMER_OFFSET_RESET"`
		ProducerRequiredAcks int      `yaml:"KAFKA_PRODUCER_REQUIRED_ACKS"`
		EnableTLS            bool     `yaml:"KAFKA_ENABLE_TLS"`
		KafkaVersion         string   `yaml:"KAFKA_VERSION"`
		Timeout              int      `yaml:"KAFKA_TIMEOUT"`
	}
	AWSS3 struct {
		AccessKeyID     string `yaml:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `yaml:"AWS_SECRET_ACCESS_KEY"`
		Region          string `yaml:"AWS_REGION"`
		EndPoint        string `yaml:"AWS_ENDPOINT"`
		Bucket          string `yaml:"AWS_BUCKET_NAME"`
	}
	OAuth struct {
		ClientID     string `yaml:"OAUTH_CLIENT_ID"`
		ClientSecret string `yaml:"OAUTH_CLIENT_SECRET"`
	}
}

var cfg Schema

func LoadConfig() *Schema {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	// Load YAML config file
	filePath := filepath.Join(currentDir, "config.sample.yaml")
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	// Decode YAML into the Schema struct
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Error decoding YAML config file: %v", err)
	}

	// Parse environment variables (optional, if you want them to override YAML values)
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing environment variables: %v", err)
	}
	if err := env.Parse(&cfg.Kafka); err != nil {
		log.Fatalf("Error parsing Kafka environment variables: %v", err)
	}

	return &cfg
}

func GetConfig() *Schema {
	return &cfg
}
