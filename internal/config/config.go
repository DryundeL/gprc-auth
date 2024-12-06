package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	DB       DBConfig      `yaml:"db"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"1h"`
	GRPC     GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DBConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode" env-default:"disable"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)

	if err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
