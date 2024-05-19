package database

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"shiroha.com/configs"
	"strconv"
)

type Redis struct {
	host     string
	password string
	port     int
}

// Init Redis配置初始化
func (r *Redis) Init() error {
	config, err := configs.LoadConfig()
	if err != nil {
		return errors.New("Error loading config: " + err.Error())
	}

	redisConfig := config.Data.Redis
	r.host = redisConfig.Host
	r.port = redisConfig.Port
	r.password = redisConfig.Password
	return nil
}

// NewRedisClient 创建并返回一个Redis客户端
// 调用者负责在适当的时候调用 Close 方法来关闭客户端连接
func NewRedisClient() (*redis.Client, error) {
	d := Redis{}
	err := d.Init()
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     d.host + ":" + strconv.Itoa(d.port),
		Password: d.password,
		DB:       0, // use default DB
	})
	ctx := context.Background()

	// 测试连接
	if _, err = rdb.Ping(ctx).Result(); err != nil {
		return nil, errors.New("Failed to connect to Redis: " + err.Error())
	}

	return rdb, nil
}
