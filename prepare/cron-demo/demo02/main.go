package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

// 任务封装结构体
type CronJob struct {
	nextTime time.Time
	expr *cronexpr.Expression
}

func main() {
	// 同时调度多个cron表达式
	var (
		cronJobTable map[string]*CronJob
		err          error
		exp          *cronexpr.Expression
	)
	// map是引用类型，使用前，先初始化
	cronJobTable = make(map[string]*CronJob)

	// job1
	if exp, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Printf("err=[%v]", err)
	}
	cronJobTable["job1"] = &CronJob{
		expr: exp,
		nextTime: exp.Next(time.Now()),
	}

	// job2
	//if exp, err = cronexpr.Parse("*/2 * * * * * *"); err != nil {
	//	fmt.Printf("err=[%v]", err)
	//}
	//cronJobTable["job2"] = &CronJob{
	//	expr: exp,
	//	nextTime: exp.Next(time.Now()),
	//}

	// 拉起一个协程调度任务
	go func() {
		var (
			name     string
			cronJob  *CronJob
			now      time.Time
			nextTime time.Time
		)
		for {
			now = time.Now()
			for name, cronJob = range cronJobTable {
				// 计算当前任务的下一次执行时间
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					go func(jobjName string) {
						fmt.Printf("exec job=[%v] \n", jobjName)
					}(name)
					// 如果是周期任务，任务被执行之后，重新计算任务的执行时间
					nextTime = cronJob.expr.Next(now)
					fmt.Printf("next exec time=[%v] \n",nextTime)
					// 更新下一次job执行的时间
					cronJob.nextTime = nextTime
				}
			}
			// 每轮循环间隔1秒
			select {
			case <-time.NewTimer(time.Millisecond * 100).C:
			}
		}
	}()
	time.Sleep(999*time.Second)
}
