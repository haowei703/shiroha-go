package tests

import (
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"reflect"
	"shiroha.com/internal/app/utils"
	"shiroha.com/internal/pkg/database"
	"sort"
	"testing"
)

func TestDatabaseConnect(t *testing.T) {
	db, err := database.NewPostgresDB()
	if err != nil {
		err := fmt.Errorf("database error: %s", err)
		t.Error(err)
	}
	if db == nil {
		err := fmt.Errorf("database nil")
		t.Error(err)
	}
}

func TestRedisConnect(t *testing.T) {
	rdb, err := database.NewRedisClient()
	if err != nil {
		err := fmt.Errorf("database error: %s", err)
		t.Error(err)
	}
	if rdb == nil {
		err := fmt.Errorf("database nil")
		t.Error(err)
	}
}

func TestRedisWriteAndRead(t *testing.T) {
	rdb, err := database.NewRedisClient()
	if err != nil {
		t.Fatalf("Failed to create Redis client: %s", err)
	}
	if rdb == nil {
		t.Fatal("Redis client is nil")
	}

	redisUtils := utils.NewRedisUtils(rdb)

	// 测试保存字符串
	keyString := "test"
	valueString := "1"
	err = redisUtils.SaveString(keyString, valueString, 200)
	if err != nil {
		t.Fatalf("Failed to save string: %s", err)
	}
	value, err := redisUtils.GetString(keyString)
	if err != nil {
		t.Fatalf("Failed to get string: %s", err)
	}
	if value != valueString {
		t.Errorf("Expected string value %s, got %s", valueString, value)
	}
	if err = redisUtils.DeleteKey(keyString); err != nil {
		t.Fatalf("Failed to delete key: %s", err)
	}
	t.Logf("Deleted key: %s", keyString)

	// 测试保存和读取结构体
	student := struct {
		Name string
		Age  int
		Sex  string
	}{
		Name: "123",
		Age:  20,
		Sex:  "male",
	}
	keyStudent := "student"
	err = redisUtils.SaveObject(keyStudent, &student, 200)
	if err != nil {
		t.Fatalf("Failed to save object: %s", err)
	}
	var testStudent struct {
		Name string
		Age  int
		Sex  string
	}
	err = redisUtils.GetObject(keyStudent, &testStudent)
	if err != nil {
		t.Fatalf("Failed to get object: %s", err)
	}
	if testStudent != student {
		t.Errorf("Expected student %#v, got %#v", student, testStudent)
	}

	err = redisUtils.DeleteKey(keyStudent)
	if err != nil {
		t.Fatalf("Failed to delete key: %s", err)
	}
	t.Logf("Deleted key: %s", keyStudent)

	// 测试保存和读取列表
	keyList := "list"
	list := []string{"1", "2", "3", "4", "5"}
	err = redisUtils.SaveList(keyList, list)
	if err != nil {
		t.Fatalf("Failed to save list: %s", err)
	}
	testList, err := redisUtils.GetList(keyList)
	if err != nil {
		t.Fatalf("Failed to get list: %s", err)
	}
	if !reflect.DeepEqual(list, testList) {
		t.Errorf("Expected list %#v, got %#v", list, testList)
	}

	err = redisUtils.DeleteKey(keyList)
	if err != nil {
		t.Fatalf("Failed to delete key: %s", err)
	}
	t.Logf("Deleted key: %s", keyList)

	// 测试保存和读取集合
	keySet := "set"
	err = redisUtils.SaveSet(keySet, list)
	if err != nil {
		t.Fatalf("Failed to save set: %s", err)
	}
	testSet, err := redisUtils.GetSet(keySet)
	if err != nil {
		t.Fatalf("Failed to get set: %s", err)
	}
	sort.Strings(testSet) // Redis 集合是无序的，需要排序后再比较
	if !reflect.DeepEqual(list, testSet) {
		t.Errorf("Expected set %#v, got %#v", list, testSet)
	}
	err = redisUtils.DeleteKey(keySet)
	if err != nil {
		t.Fatalf("Failed to delete key: %s", err)
	}
	t.Logf("Deleted key: %s", keySet)
}

func TestRedisHashTable(t *testing.T) {
	rdb, err := database.NewRedisClient()
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rdb)
	if err != nil {
		err := fmt.Errorf("database error: %s", err)
		t.Error(err)
	}
	if rdb == nil {
		err := fmt.Errorf("database nil")
		t.Error(err)
	}
	redisUtils := utils.NewRedisUtils(rdb)

	// 保存和获取哈希表字段值
	err = redisUtils.SaveHashField("user", "name", "新垣结衣")
	if err != nil {
		t.Errorf("保存哈希表字段失败: %v", err)
	}
	name, err := redisUtils.GetHashField("user", "name")
	if err != nil {
		t.Errorf("获取哈希表字段失败: %v", err)
	}
	if name != "新垣结衣" {
		t.Errorf("获取到的哈希表字段值不符合预期: %s", name)
	}

	// 删除哈希表字段
	err = redisUtils.DeleteHashField("user", "name")
	if err != nil {
		t.Errorf("删除哈希表字段失败: %v", err)
	}
	_, err = redisUtils.GetHashField("user", "name")
	if !errors.Is(err, redis.Nil) {
		t.Errorf("预期的哈希表字段已被删除，但依然存在: %v", err)
	}

	// 保存多个哈希表字段
	fields := map[string]interface{}{
		"age":     18,
		"address": "Japan",
	}
	err = redisUtils.SaveHashFields("user", fields)
	if err != nil {
		t.Errorf("保存多个哈希表字段失败: %v", err)
	}

	// 获取所有哈希表字段和值
	allFields, err := redisUtils.GetAllHashFields("user")
	if err != nil {
		t.Errorf("获取所有哈希表字段失败: %v", err)
	}
	fmt.Println("所有哈希表字段和值:")
	for k, v := range allFields {
		fmt.Printf("%s: %s\n", k, v)
	}

	// 获取哈希表字段列表
	keys, err := redisUtils.GetHashKeys("user")
	if err != nil {
		t.Errorf("获取哈希表字段列表失败: %v", err)
	}
	fmt.Println("哈希表字段列表:", keys)

	// 获取哈希表值列表
	values, err := redisUtils.GetHashValues("user")
	if err != nil {
		t.Errorf("获取哈希表值列表失败: %v", err)
	}
	fmt.Println("哈希表值列表:", values)

	// 获取哈希表长度
	length, err := redisUtils.GetHashLength("user")
	if err != nil {
		t.Errorf("获取哈希表长度失败: %v", err)
	}
	fmt.Println("哈希表长度:", length)

	// 检查哈希表字段是否存在
	exists, err := redisUtils.CheckHashExists("user", "age")
	if err != nil {
		t.Errorf("检查哈希表字段是否存在失败: %v", err)
	}
	fmt.Println("哈希表字段是否存在:", exists)

	// 保存哈希表字段，如果字段不存在则保存
	success, err := redisUtils.SaveHashIfNotExists("user", "name", "石原里美")
	if err != nil {
		t.Errorf("保存哈希表字段失败: %v", err)
	}
	fmt.Println("保存哈希表字段是否成功:", success)

	//// 删除哈希表
	err = redisUtils.DeleteHash("user")
	if err != nil {
		t.Errorf("删除哈希表失败: %v", err)
	}
}
