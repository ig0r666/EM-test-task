package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type APIConfig struct {
	AgeAPIUrl    string `yaml:"age_url" env:"AGE_URL" env-default:"agify.io"`
	GenderAPIUrl string `yaml:"gender_url" env:"GENDER_URL" env-default:"genderize.io"`
	NationAPIUrl string `yaml:"nation_url" env:"NATION_URL" env-default:"nationalize.io"`
}

type HTTPConfig struct {
	Address string        `yaml:"address" env:"ADDRESS" env-default:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
}

type Config struct {
	LogLevel   string     `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	DBAddress  string     `yaml:"db_address" env:"DB_ADDRESS"`
	HTTPConfig HTTPConfig `yaml:"http"`
	APIConfig  APIConfig  `yaml:"api"`
}

func MustLoad(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return cfg
}
