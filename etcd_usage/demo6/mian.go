package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

/*
https://coding.imooc.com/lesson/281.html#mid=18378
4-8 lease租约实现kv过期
*/

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		putResp        *clientv3.PutResponse
		kv             clientv3.KV
		getResp        *clientv3.GetResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"http://47.75.179.127:2379/"}, //服务集群
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//申请租约
	lease = clientv3.NewLease(client)

	//获得一个10秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	//拿到租约ID
	leaseId = leaseGrantResp.ID

	//自动续租
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
		return
	}
	//处理续租应答的协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已经失效")
					goto END
				} else {
					fmt.Println("收到续租应答：", keepResp.ID)
				}
			}
		}
	END:
	}()

	//h获得 KV API子集
	kv = clientv3.NewKV(client)

	//申请成功，put 一个KV,让它与租约关联起来，从而实现10后自动过期
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "666", clientv3.WithLease(leaseId)); err != nil {
		println(err)
		return
	}
	fmt.Println("写入成功", putResp.Header.Revision)
	//检查KEY过期了没有
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期:", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}
