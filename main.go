package main

import (
	"log"

	"net/http"

	"./handler"
)

func main(){
	//使用http.HandleFunc建立路由规则，本质上就是建立RESTful API与处理程序的对应关系
	//tbl_file文件相关的Handler
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadFileHandler))  //传入参数的是处理函数
	http.HandleFunc("/file/upload/", handler.UploadFileHandler )  //传入参数的是处理函数
	http.HandleFunc("/file/upload/suc", handler.UploadFileSucHandler)  //上传文件成功后会调用此函数
	http.HandleFunc("/file/meta", handler.GetFileMetaInfoHandler) //通过hash值，得到文件元信息的json数据
	http.HandleFunc("/file/meta/query_all", handler.QueryFileMetasHandler) //通过limit参数，得到批量的文件元信息数组
	http.HandleFunc("/file/download", handler.DownloadFileHandler) //通过fileHash下载对应的文件
	http.HandleFunc("/file/meta/update", handler.UpdateFileMetaInfoHandler) //更新文件信息
	http.HandleFunc("/file/delete", handler.DeleteFileHandler)  //通过hash删除相应的文件

	//tbl_user, tbl_user_token用户相关的Handler
	http.HandleFunc("/user/signup", handler.UserSignUpHandler) //用户注册，可不用token访问
	http.HandleFunc("/user/signin", handler.UserLoginInHandler) //用户登录，可不用token访问
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))  //通过token，查询用户信息

	//tbl_user_file 用户文件信息相关的Handler
	http.HandleFunc("/user/file/meta/query", handler.HTTPInterceptor(handler.GetUserFileInfoHandler))

	//设置静态资源访问，把./static映射为/static/
	http.Handle("/static/",http.StripPrefix("/static/",http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(":8080", nil)) //让webserver监听8000端口，还没有中间件程序
}
