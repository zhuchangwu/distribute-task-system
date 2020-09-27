package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)


// 通过下面这两个结构体主要是用来构建这样的filter： {"time_detail.start_time":{"$lt":timestamp}}
type TimeBefore struct {
	Before int64 `bson:"$lt"`
}

type DelCondition struct {
	Before TimeBefore `bson:"time_detail.start_time"`
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
	// todo 按照时间删除, 开始时间早于当前时间删除掉
	delCondition := &DelCondition{
		Before: TimeBefore{
			Before: time.Now().Unix(),
		},
	}
	if delResult, err = collection.DeleteMany(context.TODO(),delCondition); err != nil {
		fmt.Printf("collection.DeleteOne err=[%v]", err)
		return
	}
	fmt.Printf("delResult.DeletedCount=[%v]", delResult.DeletedCount)
}
