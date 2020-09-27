package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FindByJobName struct {
	JobName string `bson:"job_name"`
}

// 封装开始时间和结束时间
type TimePoint struct {
	StartTime int64 `bson:"start_time"`
	EndTime   int64 `bson:"end_time"`
}

type LogRecord struct {
	JobName    string     `bson:"job_name"`
	Command    string     `bson:"command"`
	JobErr     string     `bson:"job_err"`
	Content    string     `bson:"content"`
	TimeDetail *TimePoint `bson:"time_detail"`
}

// 删除记录
func main() {
	var (
		err        error
		op         *options.ClientOptions
		client     *mongo.Client
		collection *mongo.Collection
		delResult  *mongo.DeleteResult
	)

	// 构建获取连接相关的op
	var timeOut = 5 * time.Second
	op = &options.ClientOptions{
		ConnectTimeout: &timeOut,
	}

	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"), op); err != nil {
		fmt.Printf("mongo.Connect err=[%v]", err)
		return
	}
	// 获取连接
	collection = client.Database("my_db1").Collection("my_collection1")

	// 它就删除掉第一个匹配到记录
	if delResult,err = collection.DeleteOne(context.TODO(),bson.M{"time_detail":"empty"});err!=nil{
		fmt.Printf("collection.DeleteOne err=[%v]", err)
		return
	}
	fmt.Printf("delResult.DeletedCount=[%v]",delResult.DeletedCount)
}
