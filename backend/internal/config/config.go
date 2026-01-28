package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

// Config 应用配置结构体
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Log      LogConfig      `mapstructure:"log"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
	Port    int    `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
	Issuer      string `mapstructure:"issuer"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Path         string   `mapstructure:"path"`
	MaxSize      int      `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// GlobalConfig 全局配置变量
var GlobalConfig *Config

// InitConfig 初始化配置
func InitConfig() error {
	// 设置默认值
	config := &Config{
		App: AppConfig{
			Name:    "online-mall",
			Version: "1.0.0",
			Debug:   false,
			Port:    8080,
		},
		Database: DatabaseConfig{
			Driver:          "mysql",
			Host:            "localhost",
			Port:            3306,
			Username:        "root",
			Password:        "123456",
			DBName:          "online_mall",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: 3600,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
			PoolSize: 100,
		},
		JWT: JWTConfig{
			Secret:      "online-mall-jwt-secret-key-2024",
			ExpireHours: 24,
			Issuer:      "online-mall",
		},
		Upload: UploadConfig{
			Path:         "./uploads",
			MaxSize:      10,
			AllowedTypes: []string{"jpg", "jpeg", "png", "gif", "mp4"},
		},
		Log: LogConfig{
			Level:      "info",
			Filename:   "./logs/app.log",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
		},
	}

	// 加载配置文件
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// 添加配置文件路径
	configPaths := []string{
		"./configs",
		"./configs/config.yaml",
		".",
	}

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Printf("Warning: Read config file failed: %v, using default config", err)
	} else {
		log.Printf("Successfully loaded config file: %s", v.ConfigFileUsed())
	}

	// 解析配置
	if err := v.Unmarshal(config); err != nil {
		return err
	}

	// 确保日志目录存在
	if err := os.MkdirAll(filepath.Dir(config.Log.Filename), 0755); err != nil {
		return err
	}

	// 确保上传目录存在
	if err := os.MkdirAll(config.Upload.Path, 0755); err != nil {
		return err
	}

	GlobalConfig = config
	return nil
}
