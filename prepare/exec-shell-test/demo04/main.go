package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func main() {
	var (
		command   *exec.Cmd
		err       error
		bytes     []byte
		pid       int
		exitCode1 int
		exitCode2 int
		ctx context.Context
	)

	// 设置两秒超时，但是shell会睡上3s, 肯定会超时，那我们看下下面的execCode的打印情况
	ctx,_ = context.WithTimeout(context.TODO(),2*time.Second)
	command = exec.CommandContext(ctx,"/bin/bash", "-c", "sleep 1s;pwd")

	// 方法将返回执行cmd得到的标准输出，任何错误通常都是ExitErr类型的
	bytes, err = command.Output()

	pid = command.ProcessState.Pid()
	exitCode1 = command.ProcessState.ExitCode()
	statusCode, ok := command.ProcessState.Sys().(syscall.WaitStatus)
	if ok {
		exitCode2 = statusCode.ExitStatus()
	} else if err != nil {
		// 无法获取获取到执行的结果，标记错误返回值为-1
		exitCode2 = -1
	}

	fmt.Printf("os pid=[%v] \n", os.Getpid()) // 当前进程pid  os pid=[71762]
	fmt.Printf("main pid=[%v] \n", goID())    // 当前协程pid  main pid=[1]
	fmt.Printf("pid=[%v] \n", pid)            // 操作系统fork出来的那个执行shell的协程的pid  pid=[71763]
	fmt.Printf("exitCode1=[%v] \n", exitCode1)// 正常退出为0
	fmt.Printf("exitCode2=[%v] \n", exitCode2)

	if err != nil {
		// shell执行超时状态码124、正常结束0、ctrl+c退出状态码：130
		if exitCode2 == 124 {
			fmt.Printf("err timeout=[%v] \n", err)
			return
		}

		// 触发的时机：比如我们通过contex希望command执行时间超过2s后为超时，一旦超时后，执行shell的进程就会被kill
		// 那cmd.Output()就会接收到返回值，err为ExitError，程序执行到这里
		if ee, ok := err.(*exec.ExitError); ok {
			fmt.Printf("ee=[%v] \n", ee)
			fmt.Printf("ok=[%v] \n", ok)
			return
		}
	}

	fmt.Printf("result=[%v] \n", string(bytes))
}

// 为了防止协程id的滥用，影响正常的gc。go 从1.4后去掉了获取协程id的方法
// 但是可以从程序调用堆栈中获取协程的id
func goID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
