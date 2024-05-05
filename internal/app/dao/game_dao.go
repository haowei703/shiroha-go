// Package dao game_dao.go
// 对games表的CRUD操作
package dao

import (
	"gorm.io/gorm"
	"math"
	"shiroha.com/internal/app/model"
)

type GameDao struct {
	db *gorm.DB
}

func NewGameDao(db *gorm.DB) *GameDao {
	return &GameDao{db: db}
}

// CreateGame 增
func (dao *GameDao) CreateGame(game *model.Game) error {
	return dao.db.Create(game).Error
}

// DeleteGame 删
func (dao *GameDao) DeleteGame(id string) error {
	return dao.db.Delete(&model.Game{}, "gameid = ?", id).Error
}

// UpdateGame 改
func (dao *GameDao) UpdateGame(game *model.Game) error {
	return dao.db.Save(game).Error
}

// GetGameByID 查询指定id的数据
func (dao *GameDao) GetGameByID(id string) (*model.Game, error) {
	var game model.Game
	if err := dao.db.Where("gameid = ?", id).First(&game).Error; err != nil {
		return nil, err
	}
	return &game, nil
}

// ListGames 查（全部数据）
func (dao *GameDao) ListGames() ([]model.Game, error) {
	var games []model.Game
	if err := dao.db.Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}

// CountGames 计算游戏总数
func (dao *GameDao) CountGames() (int64, error) {
	var count int64
	if err := dao.db.Model(&model.Game{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ListGamesByPage 分页查询
func (dao *GameDao) ListGamesByPage(page int, pageSize int) ([]model.Game, error) {
	var games []model.Game
	offset := (page - 1) * pageSize
	if err := dao.db.Offset(offset).Limit(pageSize).Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}

// CountPages 计算分页查询的总页数
func (dao *GameDao) CountPages(pageSize int) (int64, error) {
	count, err := dao.CountGames()
	if err != nil {
		return 0, err
	}
	totalPages := int64(math.Ceil(float64(count) / float64(pageSize)))
	return totalPages, nil
}
