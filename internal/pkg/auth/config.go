package auth

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

type KeyCloakConfig struct {
	KeyCloak struct {
		BaseURL     string         `toml:"baseUrl"`
		AdminRealm  string         `toml:"adminRealm"`
		ClientRealm string         `toml:"clientRealm"`
		Admin       KeycloakAdmin  `toml:"admin"`
		Client      KeycloakClient `toml:"client"`
	} `toml:"keycloak"`
	EmailApiConfig struct {
		RequestUrl string `toml:"requestUrl"`
		AppKey     string `toml:"appKey"`
	} `toml:"email_api"`
}

type KeycloakAdmin struct {
	AdminClientID     string `toml:"adminClientID"`
	AdminClientSecret string `toml:"adminClientSecret"`
	AdminUsername     string `toml:"adminUsername"`
	AdminPassword     string `toml:"adminPassword"`
}

type KeycloakClient struct {
	ClientID     string `toml:"clientID"`
	ClientSecret string `toml:"clientSecret"`
}

// LoadConfig 加载配置文件
func LoadConfig() (*KeyCloakConfig, error) {
	rootDir := os.Getenv("ROOT")
	configPath := filepath.Join(rootDir, "configs/toml/keycloak.toml")

	var c KeyCloakConfig
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println("Error:", err)
	}
	return &c, nil
}

// User DTO实体对象，仅用于注册过程
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserInfo 用户信息
type UserInfo struct {
	Avatar string `json:"avatar"`
}
