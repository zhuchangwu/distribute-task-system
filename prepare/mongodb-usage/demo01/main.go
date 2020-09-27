package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"github.com/mongodb/mongo-go-driver/mongo"
	//"github.com/mongodb/mongo-go-driver/mongo/options"
)

func main() {
	// 第一种：建立连接、选择数据库、选择表
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
	)
	if client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		fmt.Printf("err=[%v] \n", err)
		return
	}


	// 选择数据库
	database = client.Database("my_db")

	// 选择数据表
	collection = database.Collection("my_collection")

	fmt.Printf("database name=[%v] \n",database.Name())
	fmt.Printf("collection name=[%v] \n",collection.Name())
}
