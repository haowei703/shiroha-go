package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/haowei703/shiroha/internal/app/model"
	response2 "github.com/haowei703/shiroha/internal/app/response"
	"github.com/haowei703/shiroha/internal/app/service"
	"github.com/haowei703/shiroha/internal/app/utils"
	"net/http"
	"strconv"
	"time"
)

type GameHandler struct {
	gameService service.GameService
	redisCache  utils.RedisUtils
}

func NewGameHandler(gameService *service.GameService, redisUtils *utils.RedisUtils) *GameHandler {
	return &GameHandler{gameService: *gameService, redisCache: *redisUtils}
}

func (handler *GameHandler) Use(group *gin.RouterGroup) {
	group.GET("/games", handler.listGamesByPage)
}

func (handler *GameHandler) listGamesByPage(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "page must be int"})
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "pageSize must be int"})
		return
	}

	tableName := "games"
	// 响应体
	var response = &response2.PaginatedQueryResponse{}
	response.Pagination.CurrentPage = page
	response.Pagination.PageSize = pageSize

	// 获取游戏总数
	totalCount, _ := handler.gameService.CountGames()
	response.Pagination.TotalCount = totalCount

	// 先去缓存中查找数据
	cacheKey := handler.getCacheKey(tableName, page, pageSize)
	var result map[string]string
	result, err = handler.redisCache.GetAllHashFields(cacheKey)
	for _, value := range result {
		var game model.Game
		_ = json.Unmarshal([]byte(value), &game)
		response.Games = append(response.Games, game)
	}
	// 查询成功则直接返回查询结果
	if err == nil && len(result) > 0 {
		c.AbortWithStatusJSON(http.StatusOK, response)
		return
	}

	var games []model.Game
	// 缓存中不存在则在数据库中查询
	games, err = handler.gameService.ListGamesByPage(page, pageSize)
	// 查询成功则返回查询结果
	if err == nil {
		response.Games = games
		c.JSON(http.StatusOK, response)

		// 存储在redis的多字段
		fields := make(map[string]interface{}, totalCount)
		for _, game := range games {
			jsonData, _ := json.Marshal(game)
			fields[game.GameID] = jsonData
		}
		_ = handler.redisCache.SaveHashFields(cacheKey, fields)
		// 缓存时间为1小时
		_, _ = handler.redisCache.ExpireKeyIfNotExists(cacheKey, 1*time.Hour)
	}
}

func (handler *GameHandler) getCacheKey(tableName string, page int, pageSize int) string {
	return fmt.Sprintf("%s:page:%d:size:%d", tableName, page, pageSize)
}
