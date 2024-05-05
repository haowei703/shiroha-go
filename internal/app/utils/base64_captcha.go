package utils

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"shiroha.com/internal/pkg/database"
	"time"
)

// ConfigJsonBody json request body.
type ConfigJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

type RedisStore struct {
	redisUtils *RedisUtils
}

func (store *RedisStore) Init() error {
	rdb, err := database.NewRedisClient()
	if err != nil {
		return err
	}
	store.redisUtils = &RedisUtils{rdb: rdb}
	return nil
}

func (store *RedisStore) Set(id string, value string) error {
	return store.redisUtils.SaveString(id, value, 5*time.Minute)
}

func (store *RedisStore) Get(id string, clear bool) string {
	value, err := store.redisUtils.GetString(id)
	if err != nil {
		return ""
	}
	if clear {
		if err = store.redisUtils.DeleteKey(id); err != nil {
			fmt.Printf("Redis delete key error: %v\n", err)
		}
	}
	return value
}

func (store *RedisStore) Verify(id string, answer string, clear bool) bool {
	storedValue := store.Get(id, clear)
	if storedValue == "" {
		return false
	}
	return storedValue == answer
}

func NewRedisStore() *RedisStore {
	redisStore := &RedisStore{}
	err := redisStore.Init()
	if err != nil {
		panic(err)
		return nil
	}
	return redisStore
}

var store = NewRedisStore()

// GenerateCaptcha 生成图形验证码
func GenerateCaptcha(param ConfigJsonBody) (string, string, string, error) {
	var driver base64Captcha.Driver
	//create base64 encoding captcha
	switch param.CaptchaType {
	case "audio":
		driver = param.DriverAudio
	case "string":
		driver = param.DriverString.ConvertFonts()
	case "math":
		driver = param.DriverMath.ConvertFonts()
	case "chinese":
		driver = param.DriverChinese.ConvertFonts()
	default:
		driver = param.DriverDigit
	}
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return "", "", "", err
	}
	return id, b64s, answer, nil
}

// CaptchaVerify 图形验证码校验
func CaptchaVerify(param ConfigJsonBody) bool {
	//verify the captcha
	if store.Verify(param.Id, param.VerifyValue, true) {
		return true
	}
	return false
}
