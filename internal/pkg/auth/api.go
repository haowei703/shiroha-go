package auth

import (
	"bytes"
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/goccy/go-json"
	"io"
	"net/http"
)

// KeyCloakApi 与keycloak交互方法
type KeyCloakApi struct {
	Config KeyCloakConfig
	Client *gocloak.GoCloak
}

type RequestBody struct {
	AppKey     string `json:"app_key"`
	TemplateId string `json:"template_id"`
	To         string `json:"to"`
	Data       string `json:"data"`
}

func NewKeyCloakApi(config KeyCloakConfig, client *gocloak.GoCloak) *KeyCloakApi {
	return &KeyCloakApi{Config: config, Client: client}
}

// Init 初始化操作
func Init() (*KeyCloakApi, error) {
	keyCloakConfig, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	client := gocloak.NewClient(keyCloakConfig.KeyCloak.BaseURL)
	api := NewKeyCloakApi(*keyCloakConfig, client)
	return api, nil
}

// LoginAdmin 管理员登录获取token
func (api *KeyCloakApi) LoginAdmin() (*gocloak.JWT, error) {
	client := api.Client
	config := api.Config.KeyCloak
	ctx := context.Background()

	token, err := client.LoginClient(ctx, config.Admin.AdminClientID, config.Admin.AdminClientSecret, config.AdminRealm)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// CreateUser 新建用户
func (api *KeyCloakApi) CreateUser(user User) (*string, error) {
	client := api.Client
	config := api.Config.KeyCloak
	ctx := context.Background()
	token, err := api.LoginAdmin()
	if err != nil {
		return nil, err
	}

	credential := []gocloak.CredentialRepresentation{
		{
			Type:      gocloak.StringP("password"),
			Temporary: gocloak.BoolP(false),
			Value:     gocloak.StringP(user.Password),
		},
	}

	newUser := gocloak.User{
		Email:         gocloak.StringP(user.Email),
		Enabled:       gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(false),
		Username:      gocloak.StringP("user@" + user.Email),
		Credentials:   &credential,
	}

	var uid string
	uid, err = client.CreateUser(ctx, token.AccessToken, config.ClientRealm, newUser)
	if err != nil {
		return nil, err
	}
	return &uid, nil
}

// GetUserByID 获取用户详细信息
func (api *KeyCloakApi) GetUserByID(id *string) (*gocloak.User, error) {
	client := api.Client
	config := api.Config.KeyCloak
	ctx := context.Background()
	token, err := api.LoginAdmin()
	if err != nil {
		return nil, err
	}
	return client.GetUserByID(ctx, token.AccessToken, config.ClientRealm, gocloak.PString(id))
}

// UpdateUser 更新用户
func (api *KeyCloakApi) UpdateUser(user gocloak.User) error {
	client := api.Client
	config := api.Config.KeyCloak
	ctx := context.Background()

	token, err := api.LoginAdmin()
	if err != nil {
		return err
	}
	err = client.UpdateUser(ctx, token.AccessToken, config.ClientRealm, user)
	if err != nil {
		return err
	}
	return nil
}

// SendMail 发送邮件接口
func (api *KeyCloakApi) SendMail(templateId string, to string, data string) error {
	emailConfig := api.Config.EmailApiConfig

	body := RequestBody{
		AppKey:     emailConfig.AppKey,
		TemplateId: templateId,
		To:         to,
		Data:       data,
	}

	jsonData, err := json.Marshal(body)
	req, err := http.NewRequest("POST", emailConfig.RequestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		responseBody, _ := io.ReadAll(resp.Body)
		return errors.New(string(responseBody))
	}
	return nil
}

// LoginByPassword 用户登录并获取token
func (api *KeyCloakApi) LoginByPassword(user User) (*gocloak.JWT, *gocloak.UserInfo, error) {
	client := api.Client
	config := api.Config.KeyCloak
	ctx := context.Background()

	// 获取登录token，使用邮箱+密码登录方式
	token, err := client.Login(ctx, config.Client.ClientID, config.Client.ClientSecret, config.ClientRealm, user.Email, user.Password)
	if err != nil {
		return nil, nil, err
	}
	userInfo, err := client.GetUserInfo(ctx, token.AccessToken, api.Config.KeyCloak.ClientRealm)
	if err != nil {
		return nil, nil, err
	}

	return token, userInfo, nil
}

// Logout 用户登出
func (api *KeyCloakApi) Logout(refreshToken string) error {
	client := api.Client
	config := api.Config.KeyCloak
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
	config := api.Config.KeyCloak
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
	config := api.Config.KeyCloak
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
	config := api.Config.KeyCloak
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

// ValidateToken 验证用户的 JWT token 是否有效
func (api *KeyCloakApi) ValidateToken(token string) (*bool, error) {
	client := api.Client
	config := api.Config.KeyCloak
	ctx := context.Background()

	// 使用 RPT Token, 也可以选择使用其他方式来验证
	rptResult, err := client.RetrospectToken(ctx, token, config.Client.ClientID, config.Client.ClientSecret, config.ClientRealm)
	if err != nil {
		return gocloak.BoolP(false), err
	}

	return rptResult.Active, nil
}
