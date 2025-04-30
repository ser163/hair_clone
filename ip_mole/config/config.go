package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

type CheckConfig struct {
	PortTimeout time.Duration `yaml:"portTimeout"`
	HTTPTimeout time.Duration `yaml:"httpTimeout"`
	TargetURL   string        `yaml:"targetURL"`
}

type Config struct {
	Database DBConfig    `yaml:"database"`
	Check    CheckConfig `yaml:"check"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
