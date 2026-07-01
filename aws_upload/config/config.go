package config

import (
	_ "embed"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AccessId     string `yaml:"AccessId"`     // 访问密钥ID
	AccessSecret string `yaml:"AccessSecret"` // 访问密钥
	Bucket       string `yaml:"Bucket"`       // 存储桶
	Region       string `yaml:"Region"`       // 区域
	S3Prefix     string `yaml:"S3Prefix"`     // S3 路径前缀
	LocalDir     string `yaml:"LocalDir"`     // 本地文件夹路径
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
