package auth

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type KeyCloakConfig struct {
	BaseURL     string         `yaml:"baseUrl"`
	AdminRealm  string         `yaml:"adminRealm"`
	ClientRealm string         `yaml:"clientRealm"`
	Admin       KeycloakAdmin  `yaml:"admin"`
	Client      KeycloakClient `yaml:"client"`
}

type KeycloakAdmin struct {
	AdminClientID     string `yaml:"adminClientID"`
	AdminClientSecret string `yaml:"adminClientSecret"`
	AdminUsername     string `yaml:"adminUsername"`
	AdminPassword     string `yaml:"adminPassword"`
}

type KeycloakClient struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

// LoadConfig 加载配置文件
func LoadConfig() (*KeyCloakConfig, error) {
	rootDir := os.Getenv("ROOT")
	configPath := filepath.Join(rootDir, "configs/keycloak.yaml")

	var c KeyCloakConfig
	data, err := os.ReadFile(configPath)

	if err = yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// User DTO实体对象，仅用于注册过程
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
