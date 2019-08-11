package handler

import (
	"net/http"
)

//auth.go文件用于权限控制
//与func(ResponseWriter, *Request)的一样参数的func都是HandlerFunc类型

//HTTPInterceptor: HTTP协议请求的token拦截器，处理每个handler之前需要拦截一下HTTP请求，看看请求是否携带了Token令牌
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(  //这是HandlerFunc是一个类型，类型是func
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			if len(username) < 3 || !IsTokenValid(token) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h(w,r) //这个才是handler.UserInfoHandler
		})
}
