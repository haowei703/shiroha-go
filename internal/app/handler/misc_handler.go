package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/haowei703/shiroha/internal/app/utils"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

type MiscHandler struct {
}

func NewMiscHandler() *MiscHandler {
	return &MiscHandler{}
}

func (m *MiscHandler) Use(group *gin.RouterGroup) {
	group.GET("/captcha", generateCaptchaHandleFunc)
}

// generateCaptchaHandleFunc 生成图形验证码
func generateCaptchaHandleFunc(c *gin.Context) {
	//parse request parameters
	driverDigit := base64Captcha.DriverDigit{
		Height:   80,  // 验证码图片的高度
		Width:    240, // 验证码图片的宽度
		Length:   6,   // 验证码的长度（数字个数）
		MaxSkew:  0.7, // 最大扭曲幅度
		DotCount: 80,  // 干扰点数量
	}
	param := utils.ConfigJsonBody{
		CaptchaType: "digit",      // 指定验证码类型为 'digit'
		DriverDigit: &driverDigit, // 将数字图形驱动器配置加入
	}

	id, b64s, _, err := utils.GenerateCaptcha(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, gin.H{
		"code": 200, "data": b64s, "captchaId": id, "msg": "success",
	})
}
