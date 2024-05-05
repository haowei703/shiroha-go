package configs

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

// Config 结构体用于存储配置信息
type Config struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	Data     DataConfig     `toml:"data"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	DBName   string `toml:"dbname"`
}

// DataConfig 数据配置
type DataConfig struct {
	Redis struct {
		Host     string `toml:"host"`
		Password string `toml:"password"`
		Port     int    `toml:"port"`
	} `toml:"redis"`
}

func LoadConfig() (*Config, error) {
	rootDir := os.Getenv("ROOT")
	configPath := filepath.Join(rootDir, "configs/config.toml")

	var c Config
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println("Error decoding config file:", err)
		return nil, err
	}
	return &c, nil
}

// ExportDatabaseConfig 导出数据库配置
func ExportDatabaseConfig(c *Config) *DatabaseConfig {
	return &DatabaseConfig{
		Username: c.Database.Username,
		Password: c.Database.Password,
		Host:     c.Database.Host,
		Port:     c.Database.Port,
		DBName:   c.Database.DBName,
	}
}

// ExportRedisConfig 导出redis配置
func ExportRedisConfig(c *Config) *DataConfig {
	return &DataConfig{
		Redis: c.Data.Redis,
	}
}
