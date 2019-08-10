package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"../util"

	DBLayer "../DB"
)

const(
	EncSalt = "+_)(*&^%$#@!~`" //用于加强保密性的盐值，从+到1的shift键，还有最后的`
	TokenSalt = "%!*("
)

//UserSignUpHandler: 用于处理用户注册请求的handler
func UserSignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}else if r.Method == http.MethodPost { //处理用户注册
		r.ParseForm() //解析表单数据
		userName := r.Form.Get("username")
		password := r.Form.Get("password")
		email := r.Form.Get("email")

		password = util.Sha256([]byte(password+EncSalt)) //使用SHA256生成64个加密字符
		flag := DBLayer.UserSignUpToDB(userName, password, email)
		if flag {
			w.Write([]byte("success"))
			return
		}else{
			w.Write([]byte("failed"))
			return
		}
	}
}


func UserLoginInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}else if r.Method == http.MethodPost{
		//1.校验用户信息是否合法
		r.ParseForm()
		userName := r.Form.Get("username")
		password := r.Form.Get("password")
		password = util.Sha256([]byte(password + EncSalt))
		checkOk := DBLayer.UserLoginInFromDB(userName, password)
		if !checkOk {
			w.Write([]byte("failed"))
			return
		}

		//2.分配用户一个token用于后续的权限访问
		token := GenToken(userName)
		ok := DBLayer.UpdateUserToken(userName, token)
		if !ok {
			w.Write([]byte("failed"))
			return
		}

		//3.登录成功后重定向到首页
		//w.Write([]byte("http://"+r.Host+"/static/view/home.html"))
		//使用前端js做重定向转发
		w.Write([]byte("success"))
	}
}




//使用username来生成一个40位字符的string
func GenToken(userName string) string {
	//token = MD5(username+unix_time_stamp+_tokenSalt) + unix_time_stamp[:8] => 32 + 8 = 40个字符
	ts := fmt.Sprintf("%x", time.Now().Unix()) //十六位Unix时间戳
	token := util.MD5([]byte(userName+ts+TokenSalt)) + ts[:8]
	return token
}