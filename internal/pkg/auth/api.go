package auth

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
)

// KeyCloakApi 与keycloak交互方法
type KeyCloakApi struct {
	Config KeyCloakConfig
	Client *gocloak.GoCloak
}

func NewKeyCloakApi(config KeyCloakConfig, client *gocloak.GoCloak) *KeyCloakApi {
	return &KeyCloakApi{Config: config, Client: client}
}

// Init 初始化操作
func Init() (*KeyCloakApi, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	client := gocloak.NewClient(config.BaseURL)
	api := NewKeyCloakApi(*config, client)
	return api, nil
}

// createUser 新建用户
func (api *KeyCloakApi) createUser(user User) error {
	client := api.Client
	config := api.Config
	ctx := context.Background()
	token, err := client.LoginClient(ctx, config.Admin.AdminClientID, config.Admin.AdminClientSecret, config.AdminRealm)
	if err != nil {
		return err
	}

	credential := []gocloak.CredentialRepresentation{
		{
			Type:      gocloak.StringP("password"),
			Temporary: gocloak.BoolP(false),
			Value:     gocloak.StringP(user.Password),
		},
	}

	newUser := gocloak.User{
		Email:       gocloak.StringP(user.Email),
		Enabled:     gocloak.BoolP(true),
		Username:    gocloak.StringP("user@" + user.Email),
		Credentials: &credential,
	}

	_, err = client.CreateUser(ctx, token.AccessToken, config.ClientRealm, newUser)

	if err != nil {
		return err
	}
	return nil
}

// loginByPassword 用户登录并获取token
func (api *KeyCloakApi) loginByPassword(user User) (*gocloak.JWT, *gocloak.UserInfo, error) {
	client := api.Client
	config := api.Config
	ctx := context.Background()

	// 获取登录token，使用邮箱+密码登录方式
	token, err := client.Login(ctx, config.Client.ClientID, config.Client.ClientSecret, config.ClientRealm, user.Email, user.Password)
	if err != nil {
		return nil, nil, err
	}
	userInfo, err := client.GetUserInfo(ctx, token.AccessToken, api.Config.ClientRealm)
	if err != nil {
		return nil, nil, err
	}

	return token, userInfo, nil
}

// Logout 用户登出
func (api *KeyCloakApi) Logout(refreshToken string) error {
	client := api.Client
	config := api.Config
	ctx := context.Background()

	// 调用 gocloak 的 Logout 方法
	err := client.Logout(ctx, config.Client.ClientID, config.Client.ClientSecret, config.ClientRealm, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

// RefreshToken 刷新访问令牌
func (api *KeyCloakApi) RefreshToken(refreshToken string) (*gocloak.JWT, error) {
	client := api.Client
	config := api.Config
	ctx := context.Background()

	// 调用 gocloak 的 RefreshToken 方法
	newToken, err := client.RefreshToken(ctx, refreshToken, config.Client.ClientID, config.Client.ClientSecret, config.ClientRealm)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

// CheckUserRole 检查用户是否拥有特定角色
func (api *KeyCloakApi) CheckUserRole(userID, roleName string) (bool, error) {
	client := api.Client
	config := api.Config
	ctx := context.Background()
	token, err := client.LoginClient(ctx, config.Admin.AdminClientID, config.Admin.AdminClientSecret, config.AdminRealm)
	if err != nil {
		return false, err
	}

	// 获取用户的角色
	roles, err := client.GetRealmRolesByUserID(ctx, token.AccessToken, config.ClientRealm, userID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if *role.Name == roleName {
			return true, nil
		}
	}

	return false, nil
}

// AddUserToGroup 将用户加入群组
func (api *KeyCloakApi) AddUserToGroup(userID, groupID string) error {
	client := api.Client
	config := api.Config
	ctx := context.Background()
	token, err := client.LoginClient(ctx, config.Admin.AdminClientID, config.Admin.AdminClientSecret, config.AdminRealm)
	if err != nil {
		return err
	}

	// 调用 gocloak 的 AddUserToGroup 方法
	err = client.AddUserToGroup(ctx, token.AccessToken, config.ClientRealm, userID, groupID)
	if err != nil {
		return err
	}

	return nil
}
