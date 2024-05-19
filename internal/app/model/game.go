// Package model game.go
package model

import (
	"gorm.io/datatypes"
	"time"
)

// Game 对应实体表games
type Game struct {
	GameID       string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:game_id" json:"game_id,omitempty"` // 游戏标识符，uuid
	Title        string          `gorm:"type:varchar(255);not null;column:title" json:"title,omitempty"`                          // 游戏标题
	Description  *string         `gorm:"type:text;column:description" json:"description,omitempty"`                               // 游戏描述，允许为空
	ThumbnailURL *string         `gorm:"type:varchar(255);column:thumbnail_url" json:"thumbnail_url,omitempty"`                   // 游戏缩略图地址，允许为空
	ReleaseDate  *time.Time      `gorm:"type:date;column:release_date" json:"release_date,omitempty"`                             // 发行日期，允许为空
	Developer    *string         `gorm:"type:varchar(255);column:developer" json:"developer,omitempty"`                           // 开发商，允许为空
	Rating       *float64        `gorm:"type:numeric(3,1);column:rating" json:"rating,omitempty"`                                 // 评分，允许为空
	CategoryID   *int            `gorm:"column:category_id" json:"category_id,omitempty"`                                         // 游戏分类外键，允许为空
	IsActive     bool            `gorm:"type:bool;default:true;column:is_active" json:"is_active,omitempty"`                      // 游戏是否可以正常游玩，默认为 true
	ExtraWays    *datatypes.JSON `gorm:"type:jsonb;column:extra_ways" json:"extra_ways,omitempty"`                                // 游戏额外游戏方式信息，如支持的平台和下载链接
}

// TableName overrides the table name used by Game to `games`
func (Game) TableName() string {
	return "games"
}
