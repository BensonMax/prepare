package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
	"time"
)

//任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"starttime"`
	EndTime   int64 `bson:"endtime"`
}

//一条日志
type LogRecord struct {
	JobName   string    `bson:"jobName"`   //任务名
	Command   string    `bson:"command"`   //shell命令
	Err       string    `bson:"err"`       //脚本错误
	Content   string    `bson:"content"`   //脚本输出
	TimePoint TimePoint `bson:"timePoint"` //执行时间
}

func main() {
	//向mongodb 插入一条记录
	var (
		client        *mongo.Client
		err           error
		database      *mongo.Database
		collection    *mongo.Collection
		clientOptions *options.ClientOptions
		record        *LogRecord
		result        *mongo.InsertOneResult
		docId         objectid.ObjectID
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
	//插入记录(bson）
	record = &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	if result, err = collection.InsertOne(context.TODO(), record); err != nil {
		fmt.Println(err)
	}
	// id：默认生成一个全局唯一ID,ObjectID: 12字节的二进制
	docId = result.InsertedID.(objectid.ObjectID)
	fmt.Println("自增ID:", docId.Hex())

}
