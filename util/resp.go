package util

import (
	"encoding/json"
	"fmt"
	"log"
)

//resp.go文件，用于自定义响应体信息

//自定义HTTP响应体定义
type RespMsg struct{
	Code int `json:"code"`  //自定义状态码
	Msg  string `json:"msg"` //自定义message
	Data interface{} `json:"data"` //自定义数据结构
}

//NewRespMsg: 生成Response对象
func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code:code,
		Msg:msg,
		Data:data,
	}
}

//JSONBytes: 把response对象转换为json二进制数组
func (resp *RespMsg) JSONBytes() []byte{
	data, err := json.Marshal(resp)
	if err != nil{
		log.Println(err)
	}
	return data
}

//JSONBytes: 把response对象转换为json-string
func (resp *RespMsg) JSONString() string{
	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return string(data)
}

// GenSimpleRespBytes: 只包含code和msg的响应体，以字节流形式返回
func GenSimpleRespBytes(code int, msg string) []byte {
	return []byte(fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg))  //使用``,其中的双引号"",就不需要
}

// GenSimpleRespString: 只包含code和msg的响应体，以字符串形式返回
func GenSimpleRespString(code int, msg string) string {
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg)
}


