package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"../util"

	DBLayer "../DB"
)

const(
	EncSalt = "+_)(*&^%$#@!~`" //用于加强保密性的盐值，从+到1的shift键，还有最后的`
	TokenSalt = "%!*("  //token的盐值
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


//UserLoginInHandler: 用户登录处理
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

		//2. 用户每次登录时，为用户分配一个token用于后续的权限访问，token有效期为2个月，若用户2个月内没有登录，则token失效
		token := GenToken(userName)
		ok := DBLayer.UpdateUserToken(userName, token)
		if !ok {
			w.Write([]byte("failed"))
			return
		}

		//3.登录成功后重定向到首页
		resp := &util.RespMsg{
			Code:0,   //Code=0是请求成功, Code=-1为失败
			Msg:"OK",
			Data: struct {
				Location string //登录成功后重定向的URL
				Username string //用户名
				Token	 string //登录后申请的token
			}{
				Location:"http://"+r.Host+"/static/view/home.html",
				Username:userName,
				Token:token,
			},
		}
		w.Write(resp.JSONBytes())  //把自定义消息体返回给客户端
	}
}


//UserInfoHandler: 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	r.ParseForm()
	userName := r.Form.Get("username")

	//2.验证token的有效性  => 在HTTPInterceptor拦截器中已经实现了

	//3.查询用户信息
	user, err := DBLayer.GetUserInfo(userName)
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	//4.组装response并返回用户数据
	resp := util.RespMsg{
		Code:0,
		Msg:"OK",
		Data:*user,
	}
	w.Write(resp.JSONBytes())
}





//GenToken: 用户每次登录的时候生成一个token ，生成一个40位字符的string token
func GenToken(userName string) string {
	//token = MD5(username+unix_time_stamp+_tokenSalt) + unix_time_stamp[:8] => 32 + 8 = 40个字符
	ts := fmt.Sprintf("%x", time.Now().Unix()) //十六位Unix时间戳
	token := util.MD5([]byte(userName+ts+TokenSalt)) + ts[:8] //string + string
	return token
}





//IsTokenValid: 校验客户端传来的token是否有效, true = 有效
func IsTokenValid(token string) bool {
	if len(token) < 40 { //简单校验
		return false
	}

	//1.查询自从用户登录后，生成token的ts是否过期，2个月为一个token有效期，过期换token
	//使用ParseUint可以把十六进制的string转化为int64类型的十进制数
	tokenCreateTime, err := strconv.ParseUint(token[32:],16, 64)
	if err != nil{
		log.Fatal(err)
	}
	nowTime := time.Now().Unix()  //1min = 60秒 , 60min = 3600s , 2h = 120min = 7200
	sumTokenStamp := 3600 * 24 * 30 * 2  //2个月60天
	if (nowTime - int64(tokenCreateTime)) >= int64(sumTokenStamp) { //nowTime时间点 >= tokenTime时间点 + 60天
		return false
	}

	//2.若1成立，则从DB中通过token查询出token，在数据库层面进行token的比较
	if !DBLayer.GetTokenFromDB(token) {
		return false
	}

	return true
}









