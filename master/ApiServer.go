package master

import (
	"distribute-task-system/master/common"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

// 封装http.Server
type ApiServer struct {
	server *http.Server
}

// 单例设计模式
var (
	GHttpServer *ApiServer
)

// 保存前端传递过来的http任务
// job = {"name":"job1","command":"echo hello","cronExpr":"* * * * *"}
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		jobStr string
		newJob common.Job
		preJob *common.Job
		res    []byte
	)
	// 解析表单中的数据
	if err = r.ParseForm(); err != nil {
		fmt.Printf("handleJobSave r.ParseForm() err=[%v]", err)
		goto Err
	}
	// 取出表单中的字段值
	jobStr = r.PostForm.Get("job")

	// 将job中的值反序列化进common的Protocol.Job封装结构体中
	newJob = common.Job{}
	if err = json.Unmarshal([]byte(jobStr), &newJob); err != nil {
		fmt.Printf("handleJobSave json.Unmarshal err=[%v]", err)
		goto Err
	}
	// 将任务结构体信息保存进Etcd中
	if preJob, err = GJobManager.SaveJob(&newJob); err != nil {
		fmt.Printf("handleJobSave GJobManager.SaveJob err=[%v]", err)
		goto Err
	}

	if res, err = common.BuildSuccessRes("success", preJob); err == nil {
		w.Write(res)
	}
	return

Err:
	if res, err = common.BuildFailRes(err.Error(), nil); err == nil {
		w.Write(res)
	}
	return
}

// 删除指定的任务
// Post /job/delete name=job1
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		name   string
		res    []byte
		preJob *common.Job
	)
	// 解析表单获取前端需要删除的任务名
	if err = r.ParseForm(); err != nil {
		fmt.Printf("handleJobDelete r.ParseForm() err=[%v]", err)
		goto Err
	}
	name = r.PostForm.Get("name")

	if preJob, err = GJobManager.DeleteJob(name); err != nil {
		fmt.Printf("handleJobDelete GJobManager.DeleteJob(%v) err=[%v]", name, err)
		goto Err
	}

	// 返回成功
	if res, err = common.BuildSuccessRes("success", preJob); err == nil {
		w.Write(res)
		return
	}

Err:
	if res, err = common.BuildFailRes(err.Error(), nil); err == nil {
		w.Write(res)
	}
	return
}

// 查看所有任务列表
func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		jobList []common.Job
		res     []byte
	)

	if jobList, err = GJobManager.ListJob(); err != nil {
		fmt.Printf("handleJobList GJobManager.ListJob() err=[%v]", err)
		goto Err
	}
	if res, err = common.BuildSuccessRes("success", jobList); err == nil {
		w.Write(res)
		return
	}
Err:
	if res, err = common.BuildFailRes(err.Error(), nil); err == nil {
		w.Write(res)
		return
	}
}

// 干掉任务
// 原理是： 更新 /cron/killer/任务名
// worker接收到任务之后就会杀死该任务
func handleJobKill(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		res  []byte
		name string
	)
	if err = r.ParseForm(); err != nil {
		fmt.Printf("handleJobKill r.ParseForm() err=[%v]", err)
		goto Err
	}
	name = r.PostForm.Get("name")
	if name == "" {
		fmt.Printf("name can not null	")
		err = errors.New("handleJobKill name can not null")
		goto Err
	}
	// 修改任务状态
	if err = GJobManager.KillJob(name); err != nil {
		fmt.Printf("handleJobKill GJobManager.KillJob name=[%v],err=[%v]", name, err)
		goto Err
	}
	if res, err = common.BuildSuccessRes("success", nil); err == nil {
		w.Write(res)
		return
	}
Err:
	if res, err = common.BuildFailRes(err.Error(), nil); err == nil {
		w.Write(res)
		return
	}
}

// 第二种InitApiServer的写法
// 像下面这种写法
func InitApiServer2() (err error) {
	var (
		mux               *http.ServeMux
		server            *http.Server
		listener          net.Listener
		//staticWebRootPath http.Dir     // html静态文件的根目录
		staticFileHandler http.Handler // 请求静态文件对应的handle 回调
	)

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic err=[%v] \n", err)
		}
	}()

	mux = http.NewServeMux()
	// pattern遵循最大匹配原则
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)


	// 静态文件目录
	staticFileHandler = http.FileServer(http.Dir("master/webroot"))

	// 请求/时，将用staticFileHandler处理该请求
	// 它会将staticWebRootPath中的index.html返回给客户端
	mux.Handle("/", http.StripPrefix("/",staticFileHandler))

	if listener, err = net.Listen("tcp", ":9991"); err != nil {
		fmt.Printf(	"net.Listen(\"tcp\", \"0.0.0.0:9991\") err=[%v] \n", err)
		return
	}
	server = &http.Server{
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Handler:      mux,
	}

	// 单例初始化
	GHttpServer = &ApiServer{
		server: server,
	}

	// 启动服务
	go GHttpServer.server.Serve(listener)
	return nil
}

// 第一种InitApiServer的写法
// 像下面这种写法，ListenAndServe 会被阻塞住， 但是如果你把它go出去，报的错误又收不到
func InitApiServer1() (err error) {
	var (
		mux    *http.ServeMux
		server *http.Server
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	server = &http.Server{
		Addr:         "http://localhost:9999",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Handler:      mux,
	}

	// 单例初始化
	GHttpServer = &ApiServer{
		server: server,
	}

	// 启动服务
	if err = GHttpServer.server.ListenAndServe(); err != nil {
		return // 出错后，直接返回错误
	}
	return nil
}
