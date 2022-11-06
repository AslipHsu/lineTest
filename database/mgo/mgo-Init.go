package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoUser = ""
	MongoPwd  = ""
	MongoAddr = "0.0.0.0:27017"
)

// conn mongoDB
func (DB *MongoDB) Init() {
	// 先設一個context
	DB.DBContext = context.Background()

	dbURL := "mongodb://"
	if MongoUser != "" {
		dbURL = dbURL + MongoUser + ":" + MongoPwd + "@"
	}
	dbURL = dbURL + MongoAddr
	fmt.Println("Mongo Init MongoAddr: ", dbURL)

	// 連線方式, MgoDB.DBContext = ctx
	ctx, cancel := context.WithTimeout(DB.DBContext, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
	if err != nil {
		fmt.Println("Mongo connect error", dbURL, err)
		return
	}
	MgoDB.MongoDB = client

	// 判斷是否與DB連線成功 (範例)
	err = MgoDB.Ping()
	if err != nil {
		fmt.Println("Mongo Ping error", dbURL, err)
		return
	}
}

// close mongoDB
func (DB *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(DB.DBContext, 10*time.Second)
	defer func() {
		cancel()
		DB.DBContext.Done()
	}()

	fmt.Println("closeMongo")
	if err := DB.MongoDB.Disconnect(ctx); err != nil {
		panic(err)
	}
}
