package auth

import (
	"encoding/json"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/haowei703/shiroha/configs"
	"github.com/haowei703/shiroha/internal/app/utils"
	"net/http"

	"strconv"
	"strings"
)

var (
	serverConfig *configs.ServerConfig
)

type KeyCloakRouter struct {
	engine    *gin.Engine
	authGroup *gin.RouterGroup
	api       *KeyCloakApi
}

func NewKeyCloakRouter(engine *gin.Engine) *KeyCloakRouter {
	return &KeyCloakRouter{engine: engine}
}

// Init 初始化操作
func (key *KeyCloakRouter) Init() error {
	config, err := configs.LoadConfig()
	if err != nil {
		return err
	}
	serverConfig = &config.Server

	key.authGroup = key.engine.Group("/auth")
	keyCloakApi, err := Init()
	if err != nil {
		return err
	}

	key.api = keyCloakApi
	key.authGroup.Use(ApiMiddleware(keyCloakApi))
	key.authGroup.POST("/login", Login)
	key.authGroup.POST("/register", Register)
	key.authGroup.GET("/verify", Verify)

	return nil
}

// SetProtectRule 设置保护规则，保护路由
func (key *KeyCloakRouter) SetProtectRule(routes gin.IRoutes) {
	routes.Use(ApiMiddleware(key.api))
	routes.Use(AuthenticateMiddleware())
}

// Login 登录
func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "invalid json")
	}

	api, exists := c.Get("keyCloakApi")
	if !exists {
		utils.Error(c, http.StatusInternalServerError, "internal server error")
	}
	keyCloakApi := api.(*KeyCloakApi)
	token, userInfo, err := keyCloakApi.LoginByPassword(user)
	if err != nil {
		if strings.Contains(err.Error(), "401") {
			utils.Error(c, http.StatusUnauthorized, "unauthorized")
		} else if strings.Contains(err.Error(), "400") {
			utils.Error(c, http.StatusBadRequest, "The mailbox is not verified")
		} else {
			utils.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.SetCookie("csrf_token", token.AccessToken, token.ExpiresIn, "/", serverConfig.Host, false, true)
	c.SetCookie("refresh_token", token.RefreshToken, token.RefreshExpiresIn, "/", serverConfig.Host, false, true)
	c.SetCookie("session", token.SessionState, token.ExpiresIn, "/", serverConfig.Host, false, true)
	c.SetCookie("uid", *userInfo.Sub, token.ExpiresIn, "/", serverConfig.Host, false, false)
	userDetail, err := keyCloakApi.GetUserByID(userInfo.Sub)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}
	var userData UserInfo
	if attributes := userDetail.Attributes; attributes != nil {
		if avatar, ok := (*attributes)["avatar"]; ok {
			userData.Avatar = avatar[0]
		}
	}

	utils.Success(c, userData)
}

// Register 注册
func Register(c *gin.Context) {
	var request struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		CaptchaId    string `json:"captchaId"`
		CaptchaValue int    `json:"captchaValue"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "invalid json")
	}

	param := utils.ConfigJsonBody{
		Id:          request.CaptchaId,
		VerifyValue: strconv.Itoa(request.CaptchaValue),
	}
	if !utils.CaptchaVerify(param) {
		utils.Error(c, http.StatusUnprocessableEntity, "invalid captcha")
	}

	user := User{
		Email:    request.Email,
		Password: request.Password,
	}

	api, exists := c.Get("keyCloakApi")
	if !exists {
		utils.Error(c, http.StatusInternalServerError, "internal server error")
	}

	keyCloakApi := api.(*KeyCloakApi)
	uid, err := keyCloakApi.CreateUser(user)
	if err != nil {
		utils.Error(c, http.StatusConflict, "user already exists")
	} else {
		data := struct {
			Uid   string `json:"uid"`
			Email string `json:"email"`
		}{
			Uid:   gocloak.PString(uid),
			Email: user.Email,
		}
		var jsonData []byte
		jsonData, err = json.Marshal(data)
		err = keyCloakApi.SendMail("E_100489674255", user.Email, string(jsonData))
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "unable to send email")
		}
		utils.Success(c, nil)
	}
}

// Verify 验证用户邮箱
func Verify(c *gin.Context) {
	userId := c.Query("uid")
	if userId == "" {
		utils.Error(c, http.StatusUnprocessableEntity, "invalid user id")
		return
	}

	api, exists := c.Get("keyCloakApi")
	keyCloakApi := api.(*KeyCloakApi)
	if !exists {
		utils.Error(c, http.StatusInternalServerError, "internal server error")
	}

	user := gocloak.User{
		ID:            gocloak.StringP(userId),
		EmailVerified: gocloak.BoolP(true),
	}

	err := keyCloakApi.UpdateUser(user)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	successURL := serverConfig.Host
	http.Redirect(c.Writer, c.Request, successURL, http.StatusSeeOther)
}
