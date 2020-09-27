package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"time"
)

func main() {
	var (
		cli    *clientv3.Client
		err    error
		config clientv3.Config
		resp   *clientv3.GetResponse
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

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// 第一个参数也可使用超时机制，超时方法返回，出现超时err
	// 第三个参数用户仅统计个数，而不返回详细值
	if resp, err = cli.Get(ctx, "sample_key1",clientv3.WithCountOnly()); err != nil {
		switch err {
		case context.Canceled:
			fmt.Printf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			fmt.Printf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			fmt.Printf("client-side error: %v", err)
		default:
			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	fmt.Printf("Revision=[%v] \n", resp.Header.Revision)
	fmt.Printf("Kvs=[%v] \n", resp.Kvs)
}
