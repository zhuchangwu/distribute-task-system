package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thirdparty-modules-bak/github.com/siddontang/go/bson"
	"time"
)

// 封装开始时间和结束时间
type TimePoint struct {
	StartTime int64 `bson:"start_time"`
	EndTime   int64 `bson:"end_time"`
}

// 日志记录
type LogRecord struct {
	JobName    string     `bson:"job_name"`
	Command    string     `bson:"command"`
	JobErr     string     `bson:"job_err"`
	Content    string     `bson:"content"`
	TimeDetail *TimePoint `bson:"time_detail"`
}

func main() {
	//第二种： 建立连接、选择数据库、选择表
	var (
		client          *mongo.Client
		err             error
		database        *mongo.Database
		collection      *mongo.Collection
		ctx             context.Context
		cancelFunc      context.CancelFunc
		insertOneResult *mongo.InsertOneResult
	)

	ctx, cancelFunc = context.WithCancel(context.TODO())
	defer cancelFunc()
	// 获取连接
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		fmt.Printf("err=[%v] \n", err)
		return
	}
	// 方法退出时，断开连接
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// 选择数据库
	database = client.Database("my_db1")
	// 选择数据表
	collection = database.Collection("my_collection1")

	// 构建需要插入的bson
	doc := &LogRecord{
		JobName: "create db",
		Command: "shell xxx",
		JobErr:  "empty",
		Content: "success",
		TimeDetail: &TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix(),
		},
	}
	// 插入
	if insertOneResult, err = collection.InsertOne(context.TODO(), doc); err != nil {
		fmt.Printf("err=[%v] \n", err)
		return
	}

	// 返回给我们一个全局唯一的ObjectID，12字节的二进制
	//fmt.Printf("InsertId=[%v] \n", insertOneResult.InsertedID)

}
