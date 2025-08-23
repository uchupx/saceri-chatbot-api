package config

import (
	"log"
	"os"

	"github.com/subosito/gotenv"
	"github.com/uchupx/saceri-chatbot-api/pkg/helper"
)

var configPath = []string{
	"./",
	"./config/",
}

type Config struct {
	App struct {
		Env      string
		Port     string
		GRPCPort string
		Log      string
		Secret   string
		Name     string
		Version  string
	}
	Database struct {
		URL string
	}

	Redis struct {
		Host         string
		Port         string
		Password     string
		Database     int
		PoolSize     int
		MinIdleConns int
	}

	RSA struct {
		Private string
		Public  string
	}

	Service struct {
		AuthServiceUrl string
	}
}

var config *Config

// NewConfig is a constructor for Config
func new() *Config {
	var err error
	for _, c := range configPath {
		err = gotenv.Load(c + ".env")
		if err != nil {
			log.Printf("failed to laod config from path %s.env", c)
			continue
		}
		break
	}
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config = &Config{}

	config.Database.URL = os.Getenv("DATABASE_URL")

	config.Redis.Host = os.Getenv("REDIS_HOST")
	config.Redis.Port = os.Getenv("REDIS_PORT")
	config.Redis.Password = os.Getenv("REDIS_PASSWORD")
	config.Redis.Database = int(helper.StringToUint(os.Getenv("REDIS_DATABASE")))
	config.Redis.PoolSize = 10
	config.Redis.MinIdleConns = 5

	config.RSA.Private = os.Getenv("RSA_PRIVATE_KEY")
	config.RSA.Public = os.Getenv("RSA_PUBLIC_KEY")

	config.App.Env = os.Getenv("APP_ENV")
	config.App.Port = os.Getenv("APP_PORT")
	// config.App.GRPCPort = os.Getenv("APP_GRPC_PORT")
	config.App.Log = os.Getenv("APP_LOG")
	config.App.Secret = os.Getenv("APP_SECRET")
	config.App.Name = os.Getenv("APP_NAME")
	config.App.Version = os.Getenv("APP_VERSION")

	config.Service.AuthServiceUrl = os.Getenv("AUTH_SERVICE_URL")
	return config
}

func GetConfig() *Config {
	if config == nil {
		config = new()
	}

	return config
}
