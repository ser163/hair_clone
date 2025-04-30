package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	AgentName     string        `mapstructure:"agent_name"`
	ProxyAPI      string        `mapstructure:"proxy_api"`
	TargetURL     string        `mapstructure:"target_url"`
	PingTimeout   time.Duration `mapstructure:"ping_timeout"`
	PortTimeout   time.Duration `mapstructure:"port_timeout"`
	HTTPTimeout   time.Duration `mapstructure:"http_timeout"`
	CheckInterval time.Duration `mapstructure:"check_interval"`
	Database      struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"database"`
}

func LoadConfig(config *Config) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}
