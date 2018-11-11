package main

import (
	"context"
	"fmt"
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

//jobName过滤条件
type FindByJobName struct {
	JobName string `bson:"jobName"` //对过滤条件赋值
}

func main() {
	//向mongodb 插入一条记录
	var (
		client        *mongo.Client
		err           error
		database      *mongo.Database
		collection    *mongo.Collection
		clientOptions *options.ClientOptions
		cond          *FindByJobName
		findOptions   *options.FindOptions
		cursor        mongo.Cursor
		record        *LogRecord
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

	//按照jonName字段过滤，找出JobName=job10，找出5条
	cond = &FindByJobName{JobName: "job10"}

	//设置查询条件
	findOptions = options.Find()
	findOptions.SetSkip(0)
	findOptions.SetLimit(5)

	//发起查询
	if cursor, err = collection.Find(context.TODO(), cond, findOptions); err != nil {
		fmt.Println(err)
		return
	}

	//遍历结果集
	for cursor.Next(context.TODO()) {
		//定义一个日志对象
		record = &LogRecord{}
		//反序列化bson对象
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		//释放游标
		defer cursor.Close(context.TODO())
		//打印日志
		fmt.Println(*record)
	}

}
