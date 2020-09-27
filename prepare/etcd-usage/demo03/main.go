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
		delRps   *clientv3.DeleteResponse
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
	// 根据key，删除value
	if delRps, err = cli.Delete(ctx, "sample_key1",clientv3.WithPrevKV()); err != nil {
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
	fmt.Printf("Revision=[%v] \n", delRps.Header.Revision)
}
