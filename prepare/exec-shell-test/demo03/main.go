package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// 实现：新开启一个go程去执行shell命令，然后将执行的结果返回给main线程，
// 子go程会睡眠2s，main线程在1s时将其kill掉

// 返回值封装结构体
type res struct {
	bytes []byte
	err error
}

// 主程和子程交互

func main() {
	var (
		ctx context.Context
		cancelFunc context.CancelFunc
		resultChan chan *res
		// 主程接受从chan中获取到的返回值
		target *res
	)
	// chan是引用类型，需要通过make初始化一下再使用
	resultChan = make(chan *res,100)

	ctx,cancelFunc = context.WithCancel(context.TODO())

	go func() {
		var (
			cmd *exec.Cmd
			bytes []byte
			err error
		)

		cmd = exec.CommandContext(ctx,"/bin/bash","-c","sleep 2s;pwd")
		bytes,err = cmd.CombinedOutput()
		// 写会结果
		resultChan<-&res{
			err: err,
			bytes: bytes,
		}
	}()

	time.Sleep(3*time.Second)
	cancelFunc()

	target = <- resultChan
	fmt.Printf("err=[%v] \n",target.err)
	fmt.Printf("result=[%v] \n",string(target.bytes))
}
