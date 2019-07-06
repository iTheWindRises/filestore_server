package utils

import (
	"encoding/json"
	"log"
)

//http响应通用结构
type RespMsg struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code:code,
		Msg:msg,
		Data:data,
	}
}

//转换成json byte
func (resp *RespMsg)JSONBytes() []byte {
	r,err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return r
}

//转换成json byte
func (resp *RespMsg)JSONString() string {
	r,err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return string(r)
}