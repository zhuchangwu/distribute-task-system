package main

import (
	"fmt"
	"os/exec"
)

// go 执行shell
func main() {
	var (
		command *exec.Cmd
		err error
	)
	command = exec.Command("/bin/bash", "-c","sleep 5")
	// 调用Run后会执行我们传递给他的命令，这个过程是阻塞的
	// 标准输入、标准输出、标准错误都没问题的话，返回状态为0
	err = command.Run()
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
}
