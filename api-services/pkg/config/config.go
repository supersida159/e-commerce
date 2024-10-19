package config

import (
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
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
	Environment   string `env:"ENVIRONMENT"`
	HttpPort      int    `env:"HTTP_PORT"`
	GrpcPort      int    `env:"GRPC_PORT"`
	AuthSecret    string `env:"AUTH_SECRET"`
	DatabaseURI   string `env:"DATABASE_URI"`
	RedisURI      string `env:"REDIS_URI"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`
	SecretKey     string `env:"SECRET_KEY"`
	S3BucketName  string `env:"S3_BUCKET_NAME"`
	S3Region      string `env:"S3_REGION"`
	S3APIKey      string `env:"S3_API_KEY"`
	S3SecretKey   string `env:"S3_SECRET_KEY"`
	S3Domain      string `env:"S3_DOMAIN"`
}

var (
	cfg Schema
)

func LoadConfig() *Schema {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	err := godotenv.Load(filepath.Join(currentDir, "config.yaml"))
	if err != nil {
		log.Printf("Error on load configuration file, error: %v", err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error on parsing configuration file, error: %v", err)
	}

	return &cfg
}

func GetConfig() *Schema {
	return &cfg
}
