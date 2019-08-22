package handler

import (
	"net/http"
)

//auth.go文件用于权限控制
//与func(ResponseWriter, *Request)的一样参数的func都是HandlerFunc类型

//HTTPInterceptor: HTTP协议请求的token拦截器，处理每个handler之前需要拦截一下HTTP请求，看看请求是否携带了Token令牌
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(  //这是HandlerFunc是一个类型，类型是func
		func(w http.ResponseWriter, r *http.Request) {  //先拦截request实例,并做相应的检查工作
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			if len(username) < 3 || !IsTokenValid(token) {
				w.WriteHeader(http.StatusForbidden)  //403禁止访问，权限不够
				return
			}
			h(w,r) //这个才是真正的handler处理函数
		})
}
