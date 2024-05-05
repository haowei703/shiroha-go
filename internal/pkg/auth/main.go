package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shiroha.com/internal/app/utils"
	"strconv"
	"strings"
)

type KeyCloakRouter struct {
	authGroup *gin.RouterGroup
}

func EnableOAuth2(r *gin.Engine) {

	// 认证路由组
	authGroup := r.Group("/auth")
	keyCloakApi, err := Init()
	if err != nil {
		panic(err)
	}
	authGroup.Use(ApiMiddleware(keyCloakApi))

	authGroup.POST("/login", login)
	authGroup.POST("/register", register)

}

// login 登录
func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	}

	api, exists := c.Get("keyCloakApi")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	keyCloakApi := api.(*KeyCloakApi)
	token, userInfo, err := keyCloakApi.loginByPassword(user)
	if err != nil {
		if strings.Contains(err.Error(), "401") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}
	c.SetCookie("csrf_token", token.AccessToken, token.ExpiresIn, "/", "localhost", false, true) // Secure 和 HttpOnly 标志根据需要调整
	c.SetCookie("session", token.SessionState, token.ExpiresIn, "/", "localhost", false, true)
	c.SetCookie("uid", *userInfo.Sub, token.ExpiresIn, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// register 注册
func register(c *gin.Context) {
	var request struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		CaptchaId    string `json:"captchaId"`
		CaptchaValue int    `json:"captchaValue"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	}

	param := utils.ConfigJsonBody{
		Id:          request.CaptchaId,
		VerifyValue: strconv.Itoa(request.CaptchaValue),
	}
	if !utils.CaptchaVerify(param) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Captcha verify failed"})
		return
	}

	user := User{
		Email:    request.Email,
		Password: request.Password,
	}

	api, exists := c.Get("keyCloakApi")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	keyCloakApi := api.(*KeyCloakApi)
	err := keyCloakApi.createUser(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "Use created successfully"})
	}
}

// getUserID 获取用户ID
//func getUserID(c *gin.Context) {
//	var user Use
//	if err := c.ShouldBindJSON(&user); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
//	}
//
//	api, exists := c.Get("keyCloakApi")
//	if !exists {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
//		return
//	}
//
//	keyCloakApi := api.(*KeyCloakApi)
//	uid, err := keyCloakApi.
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//	} else {
//		c.JSON(http.StatusOK, gin.H{"data": uid})
//	}
//}
