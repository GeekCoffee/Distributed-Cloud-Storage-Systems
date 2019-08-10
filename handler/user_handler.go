package handler

import (
	"io/ioutil"

	"net/http"

	"../util"

	DBLayer "../DB"
)

const(
	EncSalt = "+_)(*&^%$#@!~`" //用于加强保密性的盐值，从+到1的shift键，还有最后的`
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
		}else{
			w.Write([]byte("Failed Sign up"))
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
		if checkOk {
			w.Write([]byte("success"))
		}else{
			w.Write([]byte("failed"))
		}
	}



	//2.分配用户一个token用于后续的权限访问

	//3.
}