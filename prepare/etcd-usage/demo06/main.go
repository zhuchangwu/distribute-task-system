package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		cli          *clientv3.Client
		err          error
		config       clientv3.Config
		getRes       *clientv3.GetResponse
		nextReVision int64
		watcher      clientv3.Watcher
		watchChan    clientv3.WatchChan
		watchRes     clientv3.WatchResponse
		event        *clientv3.Event
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

	// 开启一条协程，不断的put、delete
	go func() {
		for {
			cli.Put(context.TODO(), "testPut", "...")
			cli.Delete(context.TODO(), "testPut")
			time.Sleep(1 * time.Second)
		}
	}()

	// 执行一次get，目的是获取出一个 revision，因为监听时需要执行从哪个版本开始监听
	if getRes, err = cli.Get(context.TODO(), "testPut"); err != nil {
		fmt.Printf("err=[%v]", err)
		return
	}
	nextReVision = getRes.Header.Revision
	fmt.Printf("watch form Revision=[%v] \n", nextReVision)

	// 开启监听器
	watcher = clientv3.NewWatcher(cli)

	watchChan = watcher.Watch(context.TODO(), "testPut", clientv3.WithRev(nextReVision))


	for watchRes = range watchChan {
		for _, event = range watchRes.Events {
			switch event.Type{
			case mvccpb.PUT:
				fmt.Printf("Put: K=[%v] V=[%v] Revision=[%v] \n",string(event.Kv.Key),string(event.Kv.Value),event.Kv.CreateRevision)
			case mvccpb.DELETE:
				fmt.Printf("Delete: modRevision=[%v] \n",event.Kv.ModRevision)
			}
		}
	}
}
