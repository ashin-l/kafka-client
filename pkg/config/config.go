package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// AppConfig 应用总配置
type AppConfig struct {
	AppName      string       `mapstructure:"app_name"`
	ServerConfig ServerConfig `mapstructure:"server"`
	LogConfig    LogConfig    `mapstructure:"log"`
	KafkaConfig  KafkaConfig  `mapstructure:"kafka"`
}

type ServerConfig struct {
	HTTPAddr string `mapstructure:"http_addr"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level        string `mapstructure:"level"`
	Format       string `mapstructure:"format"`
	EnableCaller bool   `mapstructure:"enable_caller"`
	// OutputPaths      []string `mapstructure:"output_paths"`
	// ErrorOutputPaths []string `mapstructure:"error_output_paths"`
	LumberjackConfig `mapstructure:"lumberjack"`
}

// LumberjackConfig lumberjack 轮转配置
type LumberjackConfig struct {
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	GroupID string   `mapstructure:"group_id"`
	Topics  []string `mapstructure:"topics"`
}

// Init 函数用于初始化配置加载
func Init(configPath string) (*AppConfig, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	var cfg AppConfig
	// 将配置反序列化到 Conf 变量
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	return &cfg, nil
}
