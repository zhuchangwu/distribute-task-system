package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
// 不添加任何限制条件，查询全部
// 从那个mongo中读取一个记录
func main() {
	var (
		err        error
		client     *mongo.Client
		collection *mongo.Collection
		cursor     *mongo.Cursor
	)

	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		fmt.Printf("err=[%v] \n", err)
		return
	}
	collection = client.Database("my_db1").Collection("my_collection1")
	//bson.A{}
	//bson.D{}
	//bson.M{}
	//bson.E{}

	// 入参2位置上应该填写一个filter，它也是一个bson，我们可以像下面这样自定义查询条件
	// 也可以使用依赖中自带的bson
	// 查询条件为自定义的结构体
	if cursor, err = collection.Find(context.TODO(), FindByJobName{JobName: "create db"}); err != nil {
		fmt.Printf("collection.Find err=[%v] \n", err)
		return
	}
	for cursor.Next(context.TODO()) {
		var detail = &LogRecord{}
		if err = cursor.Decode(&detail); err != nil {
			fmt.Printf("cursor.Decode err=[%v] \n")
			return
		}
		fmt.Printf("logDetail=[%#v] \n", detail)
	}

}
