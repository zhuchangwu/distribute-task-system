package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		err error
		bytes []byte
	)
	cmd = exec.Command("/bin/bash","-c","sleep 2;pwd")
	// 可以获取到执行命令得到的标准输出和标准错误
	if bytes,err = cmd.CombinedOutput();err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))
}
