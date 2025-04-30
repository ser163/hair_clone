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

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	PoolKey  string `yaml:"poolKey"`
	SetKey   string `yaml:"setKey"`
}

type CheckConfig struct {
	PortTimeout time.Duration `yaml:"portTimeout"`
	HTTPTimeout time.Duration `yaml:"httpTimeout"`
	TargetURL   string        `yaml:"targetURL"`
}

type PoolConfig struct {
	IpNumber int64 `yaml:"ipNumber"`
}

type Config struct {
	Database DBConfig    `yaml:"database"`
	Redis    RedisConfig `yaml:"redis"`
	Check    CheckConfig `yaml:"check"`
	Pool     PoolConfig  `yaml:"pool"`
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
