package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"github.com/mongodb/mongo-go-driver/mongo"
	//"github.com/mongodb/mongo-go-driver/mongo/options"
)

func main() {
	//第二种： 建立连接、选择数据库、选择表
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		ctx        context.Context
		cancelFunc context.CancelFunc
	)

	ctx,cancelFunc = context.WithCancel(context.TODO())
	defer cancelFunc()
	// 获取连接
	if client,err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"));err!=nil{
		fmt.Printf("err=[%v] \n",err)
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

	fmt.Printf("database name=[%v] \n", database.Name())
	fmt.Printf("collection name=[%v] \n", collection.Name())
}
