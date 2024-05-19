package utils

import (
	"bytes"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBUtils struct {
	client *mongo.Client
}

func NewMongoDBUtils(client *mongo.Client) *MongoDBUtils {
	return &MongoDBUtils{client}
}

// SaveJson 将 JSON 文档保存到指定集合中
func (m *MongoDBUtils) SaveJson(database string, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := m.client.Database(database).Collection(collection)
	insertResult, err := coll.InsertOne(context.TODO(), document)
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// FindJson 根据过滤器查询文档
func (m *MongoDBUtils) FindJson(database string, collection string, filter interface{}) ([]bson.M, error) {
	coll := m.client.Database(database).Collection(collection)
	var results []bson.M
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.TODO())
	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateJson 根据过滤器更新文档
func (m *MongoDBUtils) UpdateJson(database string, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	coll := m.client.Database(database).Collection(collection)
	updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

// DeleteJson 根据过滤器删除文档
func (m *MongoDBUtils) DeleteJson(database string, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	coll := m.client.Database(database).Collection(collection)
	deleteResult, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return deleteResult, nil
}

// UploadFile 上传文件到 GridFS
func (m *MongoDBUtils) UploadFile(database string, bucketName string, fileName string, fileData []byte) (interface{}, error) {
	bucket, err := gridfs.NewBucket(
		m.client.Database(database),
		options.GridFSBucket().SetName(bucketName),
	)
	if err != nil {
		return nil, err
	}

	uploadStream, err := bucket.OpenUploadStream(fileName)
	if err != nil {
		return nil, err
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(fileData)
	if err != nil {
		return nil, err
	}

	return uploadStream.FileID, nil
}

// DownloadFile 从 GridFS 下载文件
func (m *MongoDBUtils) DownloadFile(database string, bucketName string, fileID interface{}) ([]byte, error) {
	bucket, err := gridfs.NewBucket(
		m.client.Database(database),
		options.GridFSBucket().SetName(bucketName),
	)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = bucket.DownloadToStream(fileID, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DeleteFile 从 GridFS 删除文件
func (m *MongoDBUtils) DeleteFile(database string, bucketName string, fileID interface{}) error {
	bucket, err := gridfs.NewBucket(
		m.client.Database(database),
		options.GridFSBucket().SetName(bucketName),
	)
	if err != nil {
		return err
	}

	err = bucket.Delete(fileID)
	if err != nil {
		return err
	}

	return nil
}
