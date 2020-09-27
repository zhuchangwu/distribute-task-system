package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)
// demo 调度一个cron表达式
func main() {
	var(
		err error
		expr *cronexpr.Expression
		now time.Time
		nextTime time.Time
	)
	//支持到秒级、年级别（2018～2099）的调度
	if expr,err = cronexpr.Parse("*/5 * * * * * *");err!=nil{
		fmt.Printf("err=[%v]",err)
		return
	}

	// 计算得到下一次到调度时间
	// 需要注意的是，他计算出来的nextTime不一定就是从当前时间网后推5s，有可能计算出的时间为2s
	now = time.Now()
	nextTime = expr.Next(now)

	// 等待参数1位置的时间长度后，参数2位置的func
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Printf("after time=[%v] func exec",nextTime.Sub(now))
	})
	time.Sleep(999*time.Second)
}
