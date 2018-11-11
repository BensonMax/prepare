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
		logArr        []interface{} // C语言里面的addr,type,JAVA Object
		result        *mongo.InsertManyResult
		insertId      interface{}
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
		Command:   "echo hello1",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	//批量插入
	logArr = []interface{}{record, record, record} //初始化列表 3个record
	if result, err = collection.InsertMany(context.TODO(), logArr); err != nil {
		fmt.Println(err)
	}

	//微博、推特很早时候开源的，tweetID 算法
	//snowflake: 毫秒/微妙的当前时间 + 机器ID + 当前毫秒/微妙内的自增ID（每当毫秒变化了，重置成0，继续自增）
	for _, insertId = range result.InsertedIDs {
		//	拿着interface{},反射成objectID
		docId = insertId.(objectid.ObjectID)
		fmt.Println("自增ID:", docId.Hex())
	}
}
