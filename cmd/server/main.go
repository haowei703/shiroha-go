package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haowei703/shiroha/internal/app/dao"
	"github.com/haowei703/shiroha/internal/app/handler"
	"github.com/haowei703/shiroha/internal/app/service"
	"github.com/haowei703/shiroha/internal/app/utils"
	"github.com/haowei703/shiroha/internal/pkg/auth"
	"github.com/haowei703/shiroha/internal/pkg/database"
)

func main() {

	r := gin.Default()

	// 启用OAuth2.0认证
	authRouter := auth.NewKeyCloakRouter(r)
	err := authRouter.Init()
	if err != nil {
		panic(err)
	}
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
	gameRouterGroup := r.Group("/")
	authRouter.SetProtectRule(gameRouterGroup)
	gameHandler := handler.NewGameHandler(gameService, redisUtils)
	// 游戏相关路由组
	gameHandler.Use(gameRouterGroup)

	// 杂项路由组
	miscGroup := r.Group("misc")
	miscHandler := handler.NewMiscHandler()
	miscHandler.Use(miscGroup)

	err = r.Run(":3000")
	if err != nil {
		fmt.Println(err)
	}
}
