package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const CONFIG_PATH = "../../local.yaml"

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	GrpcClient	`yaml:"grpc_client"`
	JwtAppSecret string	`yaml:"jwt_app_secret"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD" env-default:"inf"`
}

type GrpcClient struct {
	Address		string			`yaml:"address"`
	Timeout		time.Duration	`yaml:"timeout"`
	Retries		int				`yaml:"retries"`
}

func MustLoad() *Config {
    _, b, _, _ := runtime.Caller(0)

    configPath := filepath.Join(filepath.Dir(b), CONFIG_PATH)

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
