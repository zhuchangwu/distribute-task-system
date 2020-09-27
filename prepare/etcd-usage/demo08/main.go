package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

// 使用Op实现乐观锁
func main() {
	var (
		err                   error
		cli                   *clientv3.Client
		cfg                   clientv3.Config
		lease                 clientv3.Lease
		leaseGrantRes         *clientv3.LeaseGrantResponse
		leaseId               clientv3.LeaseID
		leaseKeepAliveResChan <-chan *clientv3.LeaseKeepAliveResponse
		res                   *clientv3.LeaseKeepAliveResponse
		txn                   clientv3.Txn
		txnRes                *clientv3.TxnResponse
	)

	cfg = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 上锁：创建租约、自动续租、拿着租约去抢key
	// 这个租约保证了key的存活时间，防止抢锁成功后，斜程挂了而出现的死锁
	if cli, err = clientv3.New(cfg); err != nil {
		fmt.Printf("err=[%v]", err)
		return
	}
	lease = clientv3.NewLease(cli)

	// 创建租约
	if leaseGrantRes, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Printf("err=[%v]", err)
		return
	}

	// 获取租约id
	leaseId = leaseGrantRes.ID
	fmt.Printf("leaseId=[%v] \n", leaseId)

	// 通过这个cancelFunc实现取消自动续租
	cancelContext, cancelFunc := context.WithCancel(context.TODO())
	// 确保方法退出前，取消续租
	defer cancelFunc()
	// 保证释放租约
	defer lease.Revoke(context.TODO(), leaseId)

	// 为了避免在抢锁阶段租约过期，所以自动续租
	if leaseKeepAliveResChan, err = lease.KeepAlive(cancelContext, leaseId); err != nil {
		fmt.Printf("err=[%v]", err)
		return
	}

	// 启动协程，处理续租的结果
	go func() {
		for {
			select {
			case res = <-leaseKeepAliveResChan:
				if res != nil {
					fmt.Printf("keepAlive success id=[%v] \n", res.ID)
					return
				}
				fmt.Printf("keepAlive fail \n")
				goto END
			}
		}
	END:
	}()

	// 抢锁：抢锁的逻辑是：如果key不存在就设置上kv，表示抢到了锁，如果key存在表示抢锁失败
	// 使用etcd的事物
	txn = cli.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision("/tryLock"), "=", 0)).
		Then(clientv3.OpPut("/tryLock", "locked")).
		Else(clientv3.OpGet("/tryLock")) // 进入到else表示抢锁失败

	if txnRes, err = txn.Commit(); err != nil {
		fmt.Printf("err=[%v]", err)
		return
	}

	// 判断是否抢到了锁
	if !txnRes.Succeeded {
		fmt.Printf("未抢到锁：[%v]", string(txnRes.Responses[0].GetResponseRange().Kvs[0].Value))
		// 没抢到锁直接退出
		return
	}

	// 业务
	// 当代码执行到这里时，已经在锁内了
	fmt.Printf("抢到锁了，处理任务～")
	time.Sleep(5 * time.Second)

	// 释放锁：通过上面的两个defer 实现取消自动续租、释放租约

}
