package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Dev         bool   `default:"false"`
	WebAddress  string `default:":3000" envconfig:"WEB_ADDRESS"`
	SentryDsn   string `envconfig:"SENTRY_DSN"`
	MongoURI    string `envconfig:"MONGO_URI" required:"true"`
	MongoDBName string `envconfig:"MONGO_DB_NAME" default:"hearky"`
}

func Load() *Config {
	_ = godotenv.Load()
	var cfg Config
	err := envconfig.Process("SERVER", &cfg)
	if err != nil {
		log.Fatal("failed to load api config: ", err)
	}
	return &cfg
}
