/*
 * Hearky Server
 * Copyright (C) 2021 Hearky
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
