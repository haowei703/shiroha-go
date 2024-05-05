package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"shiroha.com/internal/app/dao"
	"shiroha.com/internal/app/handler"
	"shiroha.com/internal/app/service"
	"shiroha.com/internal/app/utils"
	"shiroha.com/internal/pkg/auth"
	"shiroha.com/internal/pkg/database"
)

func main() {

	r := gin.Default()

	// 启用OAuth2.0认证
	auth.EnableOAuth2(r)

	// 注入数据库依赖
	db, err := database.NewPostgresDB()
	if err != nil {
		panic(err)
	}
	gameDao := dao.NewGameDao(db)
	gameService := service.NewGameService(gameDao)

	rdb, err := database.NewRedisClient()
	if err != nil {
		panic(err)
	}
	redisUtils := utils.NewRedisUtils(rdb)
	gameHandler := handler.NewGameHandler(gameService, redisUtils)
	// 游戏相关路由组
	gameHandler.Use(r)

	// 应用杂项路由组
	miscGroup := handler.MiscRouterGroup{}
	miscGroup.Use(r)

	err = r.Run(":3000")
	if err != nil {
		fmt.Println(err)
	}
}
