package service

import (
	"github.com/haowei703/shiroha/internal/app/dao"
	"github.com/haowei703/shiroha/internal/app/model"
)

type GameService struct {
	gameDao *dao.GameDao
}

func NewGameService(gameDao *dao.GameDao) *GameService {
	return &GameService{
		gameDao: gameDao,
	}
}

// GetGameByID 查询指定数据
func (s *GameService) GetGameByID(id string) (*model.Game, error) {
	return s.gameDao.GetGameByID(id)
}

// CreateGame 新增数据
func (s *GameService) CreateGame(game *model.Game) error {
	return s.gameDao.CreateGame(game)
}

// UpdateGame 更新数据
func (s *GameService) UpdateGame(game *model.Game) error {
	return s.gameDao.UpdateGame(game)
}

// DeleteGame 删除数据
func (s *GameService) DeleteGame(id string) error {
	return s.gameDao.DeleteGame(id)
}

// CountGames 计数
func (s *GameService) CountGames() (int64, error) {
	return s.gameDao.CountGames()
}

// ListGames 获取全部数据
func (s *GameService) ListGames() ([]model.Game, error) {
	return s.gameDao.ListGames()
}

// ListGamesByPage 分页查询
func (s *GameService) ListGamesByPage(page int, pageSize int) ([]model.Game, error) {
	return s.gameDao.ListGamesByPage(page, pageSize)
}

// CountPages 页数计数
func (s *GameService) CountPages(pageSize int) (int64, error) {
	return s.gameDao.CountPages(pageSize)
}
