// Package model game.go
package model

import (
	"gorm.io/gorm"
	"time"
)

// Game 对应实体表games
type Game struct {
	GameID       string                  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:gameid"` // 游戏标识符，uuid
	Title        string                  `gorm:"type:varchar(255);not null;column:title"`                       // 游戏标题
	Description  *string                 `gorm:"type:text;column:description"`                                  // 游戏描述，允许为空
	ThumbnailURL *string                 `gorm:"type:varchar(255);column:thumbnailurl" `                        // 游戏缩略图地址，允许为空
	ReleaseDate  *time.Time              `gorm:"type:date;column:releasedate"`                                  // 发行日期，允许为空
	Developer    *string                 `gorm:"type:varchar(255);column:developer"`                            // 开发商，允许为空
	Rating       *float64                `gorm:"type:numeric(3,1);column:rating"`                               // 评分，允许为空
	CategoryID   *int                    `gorm:"column:categoryid"`                                             // 游戏分类外键，允许为空
	IsActive     bool                    `gorm:"type:bool;default:true;column:isactive"`                        // 游戏是否可以正常游玩，默认为 true
	ExtraWays    *map[string]interface{} `gorm:"type:jsonb;column:extraways"`                                   // 游戏额外游戏方式信息，如支持的平台和下载链接

	//// GORM-related internal fields
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the table name used by Game to `games`
func (Game) TableName() string {
	return "games"
}
