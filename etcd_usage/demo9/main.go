package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
	)
	//客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"http://47.75.179.127:2379/"},
		DialTimeout: 5 * time.Second,
	}
	//建立链接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	/*
		lease 实现锁自动过期\ 续租
		OP操作
		txn事务:if else then
	*/

	//	上锁 （创建租约，自动续租，拿着租约去抢占一个key）
	//  处理业务
	//  释放锁（取消自动续租，释放租约）

}
