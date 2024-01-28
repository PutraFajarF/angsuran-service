package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App        `yaml:"app"`
		HTTPServer `yaml:"httpserver"`
		Log        `yaml:"logger"`
		Excelize   `yaml:"excelize"`
	}

	// App -.
	App struct {
		Name        string `env-required:"true" yaml:"app_name"    env:"APP_NAME"`
		Version     string `env-required:"true" yaml:"app_version" env:"APP_VERSION"`
		Environment string `yaml:"app_environment" env:"ENVIRONMENT"`
		BaseDir     string `yaml:"app_base_dir" env:"APP_BASE_DIR"`
		TimeZone    string `yaml:"app_time_zone"`
	}

	// HTTP -.
	HTTPServer struct {
		Port    string `env-required:"true" yaml:"httpserver_port" env:"HTTP_PORT"`
		UseSSL  bool   `yaml:"httpserver_use_ssl"`
		SSLKey  string `yaml:"httpserver_ssl_key"`
		SSLCert string `yaml:"httpserver_ssl_cert"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
		Path  string `yaml:"log_path"`
	}

	// Excelize -.
	Excelize struct {
		Path string `yaml:"xlsx_path"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	baseDir := "./config/config.yml"
	err := cleanenv.ReadConfig(baseDir, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
