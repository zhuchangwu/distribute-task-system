package main

import (
	"distribute-task-system/master"
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

var (
	confFilePath string
)

func initArgs() {
	flag.StringVar(&confFilePath, "config", "./master.conf", "指定master配置文件的位置")
	flag.Parse()
}


/*func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", doLogin)
	server := &http.Server{
		Addr:         ":8081",
		WriteTimeout: time.Second * 2,
		Handler:      mux,
	}
	log.Fatal(server.ListenAndServe())
}
*/
func doLogin(writer http.ResponseWriter,req *http.Request){
	_, err := writer.Write([]byte("do login"))
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
}

func main() {
	var (
		err error
	)
	// 获取命令行中指定的配置文件地址
	initArgs()
	fmt.Printf("加载配置文件：filePath=[%v] \n", confFilePath)

	// 设置master占用系统全部核心数
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化httpServer
	if err = master.InitApiServer2(); err != nil {
		goto Err
	}

	// 初始化JobManager
	if err = master.InitJobManager(); err != nil {
		goto Err
	}


Err:
	fmt.Printf("Err=[%v] \n", err)

	time.Sleep(999 * time.Second)
}
