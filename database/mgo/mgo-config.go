package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var (
	MgoDB = MongoDB{}
)

type MongoDB struct {
	DBContext context.Context
	MongoDB   *mongo.Client
}

// 收藏一下詳細用法: https://www.jianshu.com/p/0eca791fd6b7
type SearchCMD struct {
	DBName    string        // 数据库名称
	CName     string        // 数据表名称
	SortField string        // 排序条件
	LenLimit  int           // 数量限制
	ItemID    bson.ObjectId // 数据ID
	Query     interface{}   // 查询条件
	Insert    interface{}   // 新增內容
	Update    interface{}   // 更新内容
	Delete    interface{}   // 刪除內容
	Skip      int           // 数量起始偏移值
}

type TestData struct {
	Name  string  `bson:"name"`
	Value float64 `bson:"value"`
}
