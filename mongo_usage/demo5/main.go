package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
	"time"
)

//startTime 小于某时间
//{"$lt":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

//{"timePoint.startTime":{"$lt":timestamp}}
type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	//向mongodb 插入一条记录
	var (
		client        *mongo.Client
		err           error
		database      *mongo.Database
		collection    *mongo.Collection
		clientOptions *options.ClientOptions
		delCond       *DeleteCond
		result        *mongo.DeleteResult
	)
	clientOptions = options.Client()
	clientOptions.SetConnectTimeout(time.Duration(5 * time.Second))
	//建立连接
	if client, err = mongo.Connect(context.TODO(), "mongodb://47.75.179.127:27017", clientOptions); err != nil {
		fmt.Println(err)
		return
	}
	//选择数据库
	database = client.Database("cron")
	//选择表log
	collection = database.Collection("log")

	//删除所有在当前时间前建立的日志
	//delete({"timePint.startTime":{"$lt":当前时间}})
	delCond = &DeleteCond{BeforeCond: TimeBeforeCond{Before: time.Now().Unix()}}
	//执行删除
	if result, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除行数", result.DeletedCount)
}
