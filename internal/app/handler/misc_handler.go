package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"image/png"
	"math/rand"
	"net/http"
	"shiroha.com/internal/app/utils"
	"time"
)

type MiscRouterGroup struct {
	router *gin.RouterGroup
}

func (group *MiscRouterGroup) Use(router *gin.Engine) {
	group.router = router.Group("/misc")
	// 静态路由组，返回一些静态文件
	staticGroup := group.router.Group("/static")
	{
		staticGroup.GET("/qrcode", generateQRCodeHandleFunc)
		staticGroup.GET("/captcha", generateCaptchaHandleFunc)
	}
}

// generateQRCodeHandleFunc 生成二维码
func generateQRCodeHandleFunc(c *gin.Context) {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seed := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(seed)
	code := make([]byte, 4)
	for i := 0; i < 4; i++ {
		code[i] = charset[randGen.Intn(len(charset))]
	}
	captchaText := string(code)
	qrCode, err := utils.GenerateQRCode(captchaText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "server error"})
		return
	}

	var qrCodeBuffer bytes.Buffer
	if err = png.Encode(&qrCodeBuffer, qrCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "server error"})
	}

	c.Header("Content-Type", "image/png")
	c.Data(http.StatusOK, "image/png", qrCodeBuffer.Bytes())
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
