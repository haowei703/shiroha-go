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
		Port     int    `toml:"port"`
		Password string `toml:"password"`
	} `toml:"redis"`
	MongoDB struct {
		Host       string `toml:"host"`
		Port       int    `toml:"port"`
		Username   string `toml:"username"`
		Password   string `toml:"password"`
		AuthSource string `toml:"authSource"`
	}
}

func LoadConfig() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Println(wd)
	configPath := filepath.Join(wd, "configs/toml/config.toml")

	var c Config
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println("Error decoding config file:", err)
		return nil, err
	}
	return &c, nil
}
