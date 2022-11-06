package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type TotalVolume struct {
	Name  string               `json:"name"`
	Total primitive.Decimal128 `json:"total"`
}

// get mongoDB Collection
func (DB *MongoDB) GetMongoDBCollection(dbName string, cName string, timeout int) (*mongo.Collection, context.Context, context.CancelFunc) {
	collection := DB.MongoDB.Database(dbName).Collection(cName)
	ctx, cancel := context.WithTimeout(DB.DBContext, time.Duration(timeout)*time.Second)
	return collection, ctx, cancel
}

// ping mongoDB
func (DB *MongoDB) Ping() error {
	time.Sleep(time.Duration(2) * time.Second)
	fmt.Println("Ping Mongo")
	ctx, cancel := context.WithTimeout(DB.DBContext, 2*time.Second)
	defer cancel()
	err := DB.MongoDB.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("Ping MongoDB  err.", err)
		return err
	}
	return nil
}

// insert mongoDB
func (DB *MongoDB) Insert(cmd SearchCMD) (interface{}, error) {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, cmd.Insert)
	if err != nil {
		fmt.Println(fmt.Sprintf("MongoDB Insert err. data:%v", cmd.Insert))
		return nil, err
	}

	return res.InsertedID, err
}

// search MongoDB
func (DB *MongoDB) FindAll(cmd SearchCMD, data interface{}) error {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, cmd.Query)
	if err != nil {
		fmt.Println("Find", err)
		return err
	}
	defer cur.Close(ctx)

	if err := cur.Err(); err != nil {
		fmt.Println("cur.Err:", err)
		return err
	}
	if err := cur.All(ctx, data); err != nil {
		fmt.Println("", err)
		return err
	}

	return nil
}

// 資料總數
func (DB *MongoDB) Count(cmd SearchCMD) (int64, error) {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, cmd.Query)
	if err != nil {
		fmt.Println("", err)
		return 0, err
	}
	return count, nil
}

// search MongoDB use find
func (DB *MongoDB) FindOptions(cmd SearchCMD, data interface{}, opts options.FindOptions) error {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, cmd.Query, &opts)
	if err != nil {
		fmt.Println("Find", err)
		return err
	}
	defer cur.Close(ctx)

	if err := cur.Err(); err != nil {
		fmt.Println("cur.Err:", err)
		return err
	}
	if err := cur.All(ctx, data); err != nil {
		fmt.Println("", err)
		return err
	}
	return nil
}

// 資料總數
func (DB *MongoDB) CountOptions(cmd SearchCMD, opts options.CountOptions) (int64, error) {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, cmd.Query, &opts)
	if err != nil {
		fmt.Println("", err)
		return 0, err
	}
	return count, nil
}

func (DB *MongoDB) FindOne(cmd SearchCMD, result interface{}) error {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, cmd.Query).Decode(result)
	if err != nil {
		fmt.Println("Mongo FindOne err", err)
	}
	return err
}

// 修改資料
func (DB *MongoDB) Update(cmd SearchCMD) error {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 10*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true) //若沒資料新增一筆
	filter := cmd.Query
	update := cmd.Update

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		fmt.Println("Update", err)
		return err
	}

	if result.MatchedCount != 0 {
		return nil
	}
	if result.UpsertedCount != 0 {
		fmt.Println(fmt.Sprintf("找不到匹配 已新增資料 ID: %v\n", result.UpsertedID))
		return nil
	}

	return nil
}

// 修改資料(只修改 沒資料不新增
func (DB *MongoDB) Update2(cmd SearchCMD) error {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 10*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(false) //若沒資料不新增
	filter := cmd.Query
	update := cmd.Update

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		fmt.Println("Update", err)
		return err
	}

	if result.MatchedCount != 0 {
		return nil
	}
	if result.UpsertedCount != 0 {
		fmt.Println(fmt.Sprintf("找不到匹配 已新增資料 ID: %v\n", result.UpsertedID))
		return nil
	}

	return nil
}

// 修改資料 Many
func (DB *MongoDB) UpdateMany(cmd SearchCMD) error {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 10*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true) //若沒資料新增一筆
	filter := cmd.Query
	update := cmd.Update

	result, err := collection.UpdateMany(ctx, filter, update, opts)
	if err != nil {
		fmt.Println("Update", err)
		return err
	}

	if result.MatchedCount != 0 {
		return nil
	}
	if result.UpsertedCount != 0 {
		fmt.Println(fmt.Sprintf("找不到匹配 已新增資料 ID: %v\n", result.UpsertedID))
		return nil
	}

	return nil
}

// 刪除資料
func (DB *MongoDB) Delete(cmd SearchCMD) error {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, cmd.Delete)
	if err != nil {
		fmt.Println("Delete: ", err)
		return err
	}

	return nil
}

// 刪除資料 Many
func (DB *MongoDB) DeleteMany(cmd SearchCMD) error {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 10*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, cmd.Delete)
	if err != nil {
		fmt.Println("Delete: ", err)
		return err
	}

	return nil
}

// 檢查 Mongo 的 err 是否為找不到資料
func (DB *MongoDB) IsNotFindData(err error) (bool, error) {
	if err != nil && err.Error() == "mongo: no documents in result" {
		return true, err
	}
	return false, err
}

func (DB *MongoDB) FindOneAndUpdateMongo(cmd SearchCMD, result interface{}) (bool, error) {
	var err error
	isFound := false
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	ctx, cancel := context.WithTimeout(DB.DBContext, 5*time.Second)
	defer cancel()
	filter := cmd.Query
	update := cmd.Update

	err = collection.FindOneAndUpdate(ctx, filter, update).Decode(result)
	if err != nil {
		fmt.Println("db FindOneMongo err3:", err)
		return isFound, err
	}
	if result != nil {
		isFound = true
	}
	return isFound, err
}

// UniqueIndex MongoDB
func (DB *MongoDB) CreateUniqueIndex(cmd SearchCMD, keys []string, isUnique bool) {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := collection.Indexes()
	keysDoc := bsonx.Doc{}

	// 复合索引
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}

	// 创建索引
	result, err := indexView.CreateOne(
		DB.DBContext,
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(isUnique),
		},
		opts,
	)
	if result == "" || err != nil {
		panic(err.Error())
	}
}

// Remove UniqueIndex MongoDB
func (DB *MongoDB) RemoveUniqueIndex(cmd SearchCMD, key string) {
	collection := DB.MongoDB.Database(cmd.DBName).Collection(cmd.CName)
	indexView := collection.Indexes()

	// 移除索引
	result, err := indexView.DropOne(
		DB.DBContext,
		key,
	)
	if err != nil {
		fmt.Println("Drop Index error", err)
		return
	}
	fmt.Println("Drop Index Result", result)
}
