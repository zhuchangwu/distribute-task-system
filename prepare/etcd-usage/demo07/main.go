package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)
// 学习使用Op 然后Do
func main() {
	var (
		config clientv3.Config
		cli    *clientv3.Client
		err    error
		op     clientv3.Op
		opRes  clientv3.OpResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if cli, err = clientv3.New(config); err != nil {
		fmt.Printf("err=[%v] \n", err)
		return
	}

	// 构建一个putOp动作
	op = clientv3.OpPut("testOp1", "1")
	// 执行putOp这个动作
	if opRes, err = cli.Do(context.TODO(), op); err != nil {
		fmt.Printf("err=[%v] \n", err)
		return
	}
	fmt.Printf("putOpRev=[%v] \n", opRes.Put().Header.Revision)

	// 创建一个getOp动作
	op = clientv3.OpGet("testOp1")
	// 执行getOp这个动作
	if opRes, err = cli.Do(context.TODO(), op); err != nil {
		fmt.Printf("err=[%v] \n", err)
		return
	}
	fmt.Printf("getOpRev=[%v] \n", opRes.Get().Header.Revision)
	fmt.Printf("getOpKV=[%v] \n", opRes.Get().Kvs)

	// 创建一个DeleteOp，并执行
	op = clientv3.OpDelete("testOp1",clientv3.WithPrevKV())
	if opRes, err = cli.Do(context.TODO(), op); err != nil {
		fmt.Printf("err=[%v] \n", err)
		return
	}
	fmt.Printf("delOpRev=[%v] \n",opRes.Del().Header.Revision)
	fmt.Printf("delOpKV=[%v] \n",opRes.Del().PrevKvs)// 这里想取出值来，前面需要加WithPrevKV



}
