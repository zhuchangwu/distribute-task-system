package master

import (
	"context"
	"distribute-task-system/master/common"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

// JobManager的封装结构体
type JobManager struct {
	client *clientv3.Client
	Kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	GJobManager *JobManager
)

// 保存Job至etcd中
// 保存成功，返回原来的任务
// 保存失败，返回失败的错误信息
func (JobManager *JobManager) SaveJob(job *common.Job) (preJob *common.Job, err error) {
	var (
		ctx        context.Context
		cancelFunc context.CancelFunc
		jobBytes   []byte
		putRes     *clientv3.PutResponse
		preValue   []byte
		oldJob     *common.Job
	)

	if jobBytes, err = json.Marshal(job); err != nil {
		fmt.Printf("SaveJob json.Marshal(job) err=[%v] \n", err)
		return
	}

	// 保存
	ctx, cancelFunc = context.WithTimeout(context.TODO(), 5*time.Second)
	if putRes, err = GJobManager.client.Put(ctx, common.CronJobNamePrefix+job.Name, string(jobBytes), clientv3.WithPrevKV()); err != nil {
		cancelFunc()
		fmt.Printf("GJobManager.client.Put err=[%v] \n", err)
		return
	}

	// 这里得加判断的条件，不然会报错
	if putRes.PrevKv != nil {
		preValue = putRes.PrevKv.Value
		// 方法返回值位置的变量默认是nil，所以不能直接将数据反序列化进他们中，需要像下面这样中转一下
		oldJob = &common.Job{}
		if err = json.Unmarshal(preValue, oldJob); err != nil {
			fmt.Printf("SaveJob json.Marshal(job) err=[%v] \n", err)
			return
		}
		preJob = oldJob
	}

	// 上面的逻辑都成功，执行到这里
	return preJob, nil
}

// 根据任务名称删除任务
func (JobManager *JobManager) DeleteJob(name string) (preJob *common.Job, err error) {
	var (
		delRes *clientv3.DeleteResponse
		oldJob *common.Job
	)

	if delRes, err = GJobManager.client.Delete(context.TODO(), common.CronJobNamePrefix+name, clientv3.WithPrevKV()); err != nil {
		fmt.Printf("DeleteJob GJobManager.client.Delete jobName=[%v] err=[%v]", common.CronJobNamePrefix+name, err)
		return nil, err
	}

	// 如果删除一个不存在的kv，这个if就不成立，返回的preJob就会是nil
	// 返回一个nik对象给前端，前端最终收到的就是 null
	if delRes.PrevKvs != nil && len(delRes.PrevKvs) != 0 {
		oldJob = &common.Job{}
		if err = json.Unmarshal(delRes.PrevKvs[0].Value, oldJob); err != nil {
			fmt.Printf("DeleteJob GJobManager.client.Delete jobName=[%v] err=[%v]", common.CronJobNamePrefix+name, err)
			return nil, err
		}
		preJob = oldJob
	}
	return preJob, nil
}

// 查看出所有的任务
func (JobManager *JobManager) ListJob() (jobs []common.Job, err error) {
	var (
		getRes *clientv3.GetResponse
		kvPair *mvccpb.KeyValue
	)

	if getRes, err = GJobManager.client.Get(context.TODO(), common.CronJobNamePrefix, clientv3.WithPrefix()); err != nil {
		fmt.Printf("ListJob GJobManager.client.Get err=[%v]", err)
		return nil, err
	}

	jobs = make([]common.Job, 0)

	// 遍历返回结果，反序列化进jobs数组中
	for _, kvPair = range getRes.Kvs {
		tmpJob := common.Job{}
		if err = json.Unmarshal(kvPair.Value, &tmpJob); err != nil {
			fmt.Printf("ListJob json.Unmarshal err=[%v]", err)
			// 有序列化错误的地方跳过，不影响其他任务的序列化情况
			continue
		}
		jobs = append(jobs, tmpJob)
	}

	return jobs, nil
}

// kill任务
// 原理：修改(写入) /cron/killer/jobName
// worker检测到jobName被修改后会kill任务
// 然后通过lease设置这个key到过期时间，目的是不让它浪费etcd到存储空间
func (JobManager *JobManager) KillJob(name string) (err error) {
	var (
		leaseGrantRes *clientv3.LeaseGrantResponse
		leaseId       clientv3.LeaseID
	)

	// 获取到租约
	if leaseGrantRes, err = GJobManager.lease.Grant(context.TODO(), 2); err != nil {
		fmt.Printf("KillJob GJobManager.client.Grant err=[%v]", err)
		return
	}
	// putKv with lease
	leaseId = leaseGrantRes.ID
	if _, err = GJobManager.client.Put(context.TODO(), common.KillCronJobNamePrefix+name, "", clientv3.WithLease(leaseId));err!=nil{
		fmt.Printf("KillJob GJobManager.client.Put err=[%v]", err)
		return
	}
	return
}

// 初始化JobManager，在main方法启动的时候调用了，目的是初始化单例GJobManager
func InitJobManager() (err error) {
	var (
		client *clientv3.Client
		Kv     clientv3.KV
		lease  clientv3.Lease
		config clientv3.Config
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Printf(" clientv3.New(config) err=[%v]")
		return
	}

	Kv = clientv3.NewKV(client)

	lease = clientv3.NewLease(client)

	GJobManager = &JobManager{
		client: client,
		Kv:     Kv,
		lease:  lease,
	}
	return nil
}
