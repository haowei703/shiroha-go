package utils

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisUtils struct {
	rdb *redis.Client
}

func NewRedisUtils(rdb *redis.Client) *RedisUtils {
	return &RedisUtils{rdb: rdb}
}

// SaveString 保存一个字符串到 Redis
func (utils *RedisUtils) SaveString(key string, value string, expiration time.Duration) error {
	return utils.rdb.Set(context.Background(), key, value, expiration).Err()
}

// GetString 从 Redis 获取一个字符串
func (utils *RedisUtils) GetString(key string) (string, error) {
	return utils.rdb.Get(context.Background(), key).Result()
}

// SaveObject 保存一个对象到 Redis，对象通过 JSON 进行序列化
func (utils *RedisUtils) SaveObject(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return utils.rdb.Set(context.Background(), key, jsonData, expiration).Err()
}

// GetObject 从 Redis 获取一个对象，对象通过 JSON 进行反序列化
func (utils *RedisUtils) GetObject(key string, target interface{}) error {
	result, err := utils.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), target)
}

// SaveList 保存一个列表到 Redis
func (utils *RedisUtils) SaveList(key string, values []string) error {
	return utils.rdb.RPush(context.Background(), key, values).Err()
}

// GetList 从 Redis 获取一个列表
func (utils *RedisUtils) GetList(key string) ([]string, error) {
	return utils.rdb.LRange(context.Background(), key, 0, -1).Result()
}

// AddToList 向 Redis 中的已存在列表添加新元素
func (utils *RedisUtils) AddToList(key string, value string) error {
	return utils.rdb.RPush(context.Background(), key, value).Err()
}

// RemoveFromList 从 Redis 列表中删除元素
func (utils *RedisUtils) RemoveFromList(key string, value string, count int) error {
	return utils.rdb.LRem(context.Background(), key, int64(count), value).Err()
}

// SaveSet 保存一个集合到 Redis
func (utils *RedisUtils) SaveSet(key string, members []string) error {
	return utils.rdb.SAdd(context.Background(), key, members).Err()
}

// GetSet 从 Redis 获取一个集合
func (utils *RedisUtils) GetSet(key string) ([]string, error) {
	return utils.rdb.SMembers(context.Background(), key).Result()
}

// AddToSet 向 Redis 中的已存在集合添加新成员
func (utils *RedisUtils) AddToSet(key string, member string) error {
	// 使用 SAdd 向集合添加一个新成员
	return utils.rdb.SAdd(context.Background(), key, member).Err()
}

// RemoveFromSet 从 Redis 集合中删除一个或多个元素
func (utils *RedisUtils) RemoveFromSet(key string, members ...string) error {
	return utils.rdb.SRem(context.Background(), key, members).Err()
}

// DeleteKey 删除 Redis 中的一个键
func (utils *RedisUtils) DeleteKey(key string) error {
	return utils.rdb.Del(context.Background(), key).Err()
}

// SaveHashField 保存哈希表中的一个字段值
func (utils *RedisUtils) SaveHashField(key, field, value string) error {
	return utils.rdb.HSet(context.Background(), key, field, value).Err()
}

// GetHashField 获取哈希表中的一个字段值
func (utils *RedisUtils) GetHashField(key, field string) (string, error) {
	return utils.rdb.HGet(context.Background(), key, field).Result()
}

// DeleteHashField 删除哈希表中的一个字段
func (utils *RedisUtils) DeleteHashField(key, field string) error {
	return utils.rdb.HDel(context.Background(), key, field).Err()
}

// GetAllHashFields 获取哈希表中的所有字段和对应的值
func (utils *RedisUtils) GetAllHashFields(key string) (map[string]string, error) {
	return utils.rdb.HGetAll(context.Background(), key).Result()
}

// GetHashKeys 获取哈希表中的所有字段
func (utils *RedisUtils) GetHashKeys(key string) ([]string, error) {
	return utils.rdb.HKeys(context.Background(), key).Result()
}

// GetHashValues 获取哈希表中的所有值
func (utils *RedisUtils) GetHashValues(key string) ([]string, error) {
	return utils.rdb.HVals(context.Background(), key).Result()
}

// GetHashLength 获取哈希表的长度
func (utils *RedisUtils) GetHashLength(key string) (int64, error) {
	return utils.rdb.HLen(context.Background(), key).Result()
}

// CheckHashExists 检查哈希表中是否存在指定字段
func (utils *RedisUtils) CheckHashExists(key, field string) (bool, error) {
	return utils.rdb.HExists(context.Background(), key, field).Result()
}

// SaveHashFields 保存多个字段到哈希表中
func (utils *RedisUtils) SaveHashFields(key string, fields map[string]interface{}) error {
	return utils.rdb.HMSet(context.Background(), key, fields).Err()
}

// SaveHashIfNotExists 保存哈希表中的一个字段值，如果字段不存在则保存
func (utils *RedisUtils) SaveHashIfNotExists(key, field, value string) (bool, error) {
	return utils.rdb.HSetNX(context.Background(), key, field, value).Result()
}

// DeleteHash 删除哈希表
func (utils *RedisUtils) DeleteHash(key string) error {
	return utils.rdb.Del(context.Background(), key).Err()
}

// ExpireKey 设置键的过期时间
func (utils *RedisUtils) ExpireKey(key string, expiration time.Duration) error {
	return utils.rdb.Expire(context.Background(), key, expiration).Err()
}
func (utils *RedisUtils) ExpireKeyIfNotExists(key string, expiration time.Duration) (bool, error) {
	return utils.rdb.ExpireNX(context.Background(), key, expiration).Result()
}
