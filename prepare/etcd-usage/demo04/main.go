package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		cli                *clientv3.Client
		err                error
		config             clientv3.Config
		leaseGrantResponse *clientv3.LeaseGrantResponse
		leaseId            clientv3.LeaseID
		putRes             *clientv3.PutResponse
		getRes             *clientv3.GetResponse
	)

	config = clientv3.Config{
		// 如果有多个节点的话，可以配置多个
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}

	if cli, err = clientv3.New(config); err != nil {
		fmt.Printf("err=[%v]", err)
		return
	}
	defer cli.Close()

	// 租赁
	if leaseGrantResponse, err = cli.Lease.Grant(context.TODO(), 10); err != nil {
		fmt.Printf("err=[%v]", err)
		return
	}

	// 获取租约id
	leaseId = leaseGrantResponse.ID

	// put + lease
	// 当租约过期时，会自动将这对KV删除
	if putRes, err = cli.Put(context.TODO(), "/addr/tj", "tianjin", clientv3.WithLease(leaseId)); err != nil {
		fmt.Printf("err=[%v]", err)
		return
	}

	fmt.Printf("putRes.Header.Revision=[%v] \n",	putRes.Header.Revision)

	// 开启协程不断检查key是否还存在
	go func(cli *clientv3.Client) {
		for {
			if getRes, err = cli.Get(context.TODO(), "/addr/tj"); err != nil {
				fmt.Printf("err=[%v]", err)
				return
			}
			if getRes.Count == 0 {
				fmt.Println("kv has been deleted \n")
				goto EOF
			}
			fmt.Printf("getRes=[%v] \n", getRes.Kvs)
			time.Sleep(1*time.Second)
		}
		EOF:
	}(cli)

	time.Sleep(999 * time.Second)
}
