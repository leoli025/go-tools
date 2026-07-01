package config

import (
	_ "embed"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ApiKey    string `yaml:"ApiKey"`
	SheetId   string `yaml:"SheetId"`
	RangeName string `yaml:"RangeName"`
	OutputDir string `yaml:"OutputDir"`
}

//go:embed dev.yaml
var configFile string

func NewConfig() *Config {
	var cfg Config
	err := yaml.Unmarshal([]byte(configFile), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &cfg
}
