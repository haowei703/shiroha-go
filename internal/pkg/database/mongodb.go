package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shiroha.com/configs"
)

type MongoDBOptions struct {
	Host       string
	Port       int
	Username   string
	Password   string
	AuthSource string
}

func (m *MongoDBOptions) Init() error {
	config, err := configs.LoadConfig()
	if err != nil {
		return err
	}
	mongoDBConfig := config.Data.MongoDB
	m.Host = mongoDBConfig.Host
	m.Port = mongoDBConfig.Port
	m.Username = mongoDBConfig.Username
	m.Password = mongoDBConfig.Password
	m.AuthSource = mongoDBConfig.AuthSource
	return nil
}

// NewMongoDBClient 返回一个新的MongoDB客户端实例
func NewMongoDBClient() (*mongo.Client, error) {
	var m MongoDBOptions
	err := m.Init()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("mongodb://%s:%d", m.Host, m.Port)
	clientOptions := options.Client().ApplyURI(uri).SetAuth(options.Credential{
		Username:   m.Username,
		Password:   m.Password,
		AuthSource: m.AuthSource,
	})

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
