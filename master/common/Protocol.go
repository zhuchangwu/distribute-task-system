package common

import (
	"encoding/json"
	"fmt"
)

var (
	//定制任务
	CronJobNamePrefix = "/cron/jobs/"
	KillCronJobNamePrefix = "/cron/killer/"
)

// 任务的封装结构体
type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

// 返回给前端的相应体封装结构体
type ResResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 构建成功的返回值
func BuildSuccessRes(msg string, data interface{}) (resByte []byte, err error) {
	var (
		res *ResResult
	)
	res = &ResResult{
		Code: 200,
		Msg:  msg,
		Data: data,
	}

	if resByte, err = json.Marshal(res); err != nil {
		fmt.Printf("BuildSuccessRes json.Marshal(res) err=[%v]", err)
		return nil, err
	}
	return
}

// 构建失败的返回值
func BuildFailRes(msg string, data interface{}) (resByte []byte, err error) {

	var (
		res *ResResult
	)
	res = &ResResult{
		Code: 500,
		Msg:  msg,
		Data: data,
	}

	if resByte, err = json.Marshal(res); err != nil {
		fmt.Printf("BuildSuccessRes json.Marshal(res) err=[%v]", err)
		return nil, err
	}
	return
}
