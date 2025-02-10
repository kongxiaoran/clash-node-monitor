package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Email EmailConfig `yaml:"email"`
	Clash ClashConfig `yaml:"clash"`
}

type EmailConfig struct {
	SMTPHost     string   `yaml:"smtp_host"`
	SMTPPort     int      `yaml:"smtp_port"`
	Username     string   `yaml:"username"`
	Password     string   `yaml:"password"`
	From         string   `yaml:"from"`
	To           []string `yaml:"to"`
	Subject      string   `yaml:"subject"`
}

type ClashConfig struct {
	ConfigPath string `yaml:"config_path"`
	Timeout    int    `yaml:"timeout"`
	Interval   int    `yaml:"interval"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}