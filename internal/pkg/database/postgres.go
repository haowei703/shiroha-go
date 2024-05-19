package database

import (
	"errors"
	"fmt"
	"github.com/haowei703/shiroha/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	username string
	password string
	host     string
	port     int
	dbname   string
}

func (d *PostgresDB) Init() error {
	config, err := configs.LoadConfig()
	if err != nil {
		return errors.New("Error loading config: " + err.Error())
	}

	databaseConfig := config.Database
	d.username = databaseConfig.Username
	d.password = databaseConfig.Password
	d.host = databaseConfig.Host
	d.port = databaseConfig.Port
	d.dbname = databaseConfig.DBName
	return nil
}

// NewPostgresDB 创建并返回一个postgresDB连接
func NewPostgresDB() (*gorm.DB, error) {
	d := PostgresDB{}
	err := d.Init()
	if err != nil {
		return nil, err
	}

	// 连接数据库
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d.host, d.port, d.username, d.password, d.dbname)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}
	return db, nil
}
