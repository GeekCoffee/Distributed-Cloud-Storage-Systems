package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mgutz/str"

	"../util"

	DBLayer "../DB"
)

const(
	EncSalt = "+_)(*&^%$#@!~`" //用于加强保密性的盐值，从+到1的shift键，还有最后的`
	TokenSalt = "%!*("  //token的盐值
	AvatarImgPathPrefix = "./static/img/avatar/"
	HttpDomain = "http://localhost:8080/" //HTTP域名前缀
	VideoDir = ""
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
		phone := r.Form.Get("phone")

		password = util.Sha256([]byte(password+EncSalt)) //使用SHA256生成64个加密字符
		flag := DBLayer.UserSignUpToDB(userName, password, email, phone)
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
			Code:200,   //参照HTTP Status Code
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
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//1.解析请求参数
	r.ParseForm()
	userName := r.Form.Get("username")

	//2.验证token的有效性  => 在HTTPInterceptor拦截器中已经实现了

	//3.查询用户信息
	user, err := DBLayer.GetUserInfo(userName)
	respErr := util.RespMsg{
		Code:500,
		Msg:"failed",
		Data:nil,
	}
	if err != nil{
		w.Write(respErr.JSONBytes())
		return
	}

	//4.组装response并返回用户数据
	resp := util.RespMsg{
		Code:200,
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

//UpdateAvatarImg: 上传用户头像
func UpdateAvatarImg(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {   //从HTTP请求消息体中可以获取得到是否是GET方法
		// "./"指的是根目录
		// io调用了操作系统的系统调用接口，委托文件系统进行文件的读写操作，然后把数据放入byteStream指定的变量域中
		byteStream, err := ioutil.ReadFile("./static/view/testImg.html")
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)  //读取文件的时候发生错误
			return
		}
		//把byteStream变量的值写入到响应体的content部分中，然后由web服务程序返回到客户端
		io.WriteString(w, string(byteStream))
	}

	if r.Method != http.MethodPut {  //按表单方式上传
		w.Write([]byte("please use GET！"))
		return
	}

	r.ParseForm()  //解析HTTP消息请求体中的表单数据，结构放入Form结构体中
	userName := r.Form.Get("username")

	userName = "root"
	userNamePre := userName + "/"
	//token := r.Form.Get("token")

	// 前端表单中的name域
	imgFile, fileHeader, err := r.FormFile("imgfile")
	if err != nil{
		fmt.Printf("Failured to read img_file, err:%s\n", err.Error())
		return
	}
	defer imgFile.Close()  //文件系统打开了这个文件进程，不使用就关闭这个文件描述符对应的进程

	//创建一个头像目录，目录名是username，头像文件名是img_hash后的名字
	err = os.MkdirAll(AvatarImgPathPrefix + userNamePre, 0666)
	if err != nil{
		log.Fatal(err)
		return
	}

	fileType := getFileType(fileHeader.Filename)  //得到文件类型，带"."的后缀名

	oldPath := AvatarImgPathPrefix + userNamePre + fileHeader.Filename
	newFile, err := os.Create(oldPath)  //创建一个空文件
	if err != nil{
		fmt.Printf("Faliured to create file, err:%s\n", err.Error())
		return
	}

	_, err = io.Copy(newFile, imgFile)  //调用文件系统接口，委托文件系统做两个文件间的数据覆盖工作
	if err != nil{
		fmt.Printf("Failured to copy file, error:%s\n", err.Error())
		return
	}

	//使用SHA1算法，生成hash值，在此之前把文件指针归为文件开头位置，因为后面要使用到copy函数
	newFile.Seek(0,0) //whence是相对位置，0表示从文件头开始，偏移offset个位置，然后开始读写文件
	imgFileHash := util.FileMD5(newFile)
	newPath := AvatarImgPathPrefix + userNamePre + imgFileHash + fileType
	newFile.Close()  //关闭掉此文件的文件句柄，不然这个文件句柄进程不能被其他进程所读写
	err = os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println("os.Rename is error: ", err)
		return
	}

	//http://localhost:8080/static/img/avatar/williamchen/3f41fa9ae00a16b32a0a1ad68847e053.
	avatarImgURL := HttpDomain + "static/img/avatar/" + userNamePre + imgFileHash + fileType

	//把用户头像的URL存入数据库
	ok := DBLayer.UpdateAvatarURL(userName, avatarImgURL)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

//getFileType: 得到文件后缀名
func getFileType(fileName string) string {
	pos1 := str.IndexOf(fileName, ".", 0)
	if pos1 == -1 { //匹配失败
		return ""
	}

	//说明匹配到了，提取出所有页面的pageId，包括主页的index.jsp的ID，即index
	pos1 = pos1 + len(".")
	fileType := str.Substr(fileName, pos1, 3) //从pos1位置开始，截取pos2-pos1长度的字符
	if fileType == "jpe" || fileType == "JPE"{
		fileType = "jpeg"
	}
	fileType = "." + strings.ToLower(fileType)  //把所有后缀名都转换为小写
	return fileType
}

//GetUserAvatarImg: 单独获取用户头像
func GetUserAvatarImg(w http.ResponseWriter, r *http.Request){
	//传入username和token，token在HTTP拦截器已经检查过了，所以只需要获取username或者其他参数就行了
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest) //请求方法不对
		return
	}

	r.ParseForm()
	userName := r.Form.Get("username")
	userName = "root"

	//从数据库获取avatarImgURL
	avatarImgURL, ok := DBLayer.GetAvatarImgURL(userName)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError) //服务器内部错误
		return
	}

	w.Write([]byte(avatarImgURL))
}



func GetVideo(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	vid := r.Form.Get("vid")

	videoPath := VideoDir + vid    //得到video文件的完整path

	fmt.Println(videoPath)
	videoFile, err := os.Open(videoPath)
	if err != nil{
		//服务端内部错误，HTTP：500
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	//videoFile是*File类型，File对象表示一个已经打开的文件，即其文件数据已经被加载到用户态内存区域中了
	defer videoFile.Close()

	//得到videoFile对象引用后，进行转化为mp4编码的二进制流，并输出给client
	//设置响应体Header，添加头字段，设置MIME类型
	w.Header().Set("Content-Type","video/mp4")  //client请求video视频文件，以mp4编码格式返回给client
	http.ServeContent(w, r, "", time.Now(), videoFile)


}



