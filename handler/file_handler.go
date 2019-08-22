package handler

// file_handler.go文件用于具体实现相关文件的RESTful API的处理函数

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"../meta"
	"../util"

	DBLayer "../DB"
)

const (
	Upload_File_Path_Prefix= "./tmp/"   //上传到的文件路径的前缀
	BaseTimeFormat = "2006-01-02 15:04:05"  //标准时间模板
)


//上传文件的REST ful接口，用于给client访问的URI资源程序
func UploadFileHandler(w http.ResponseWriter, r *http.Request){  //响应体w，和请求消息体r
	if r.Method == http.MethodGet { //返回上传文件的HTML页面
	//使用ioutil工具读取文件数据
	//把文件转化为字节流
	byteStream, err := ioutil.ReadFile("./static/view/upload.html")  // "./"指的是根目录
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)  //读取文件的时候发生错误
		return
	}
	io.WriteString(w, string(byteStream)) //把字节流转化为string类型，写入响应体中的content部分，并返回

	}else if r.Method == http.MethodPost {  //REST ful中规定，POST是上传文件,接收文件并存储到本地

		//从request请求体中，读取form表单中<input>中的文件file数据
		//表单中的File是Multipart.File类型的File，并不是我们文件系统中的File类型
		file, fileHeader, err := r.FormFile("file") //key是<input />标签中的name的值
		if err != nil{
			fmt.Printf("Failured to read file, err:%s\n", err.Error())
			return
		}
		defer file.Close()

		//记录文件元数据信息
		fileMeta := meta.FileMeta{
			FileName:fileHeader.Filename,  //包括文件后缀名
			Location:Upload_File_Path_Prefix + fileHeader.Filename,
			UploadTime:time.Now().Format(BaseTimeFormat),
		}


		//然后在本地磁盘创建一个对应的文件，用于存储用户上传的文件数据
		//os.Create方法只能创建文件，创建不了目录，创建目录使用os.Mkdir或者os.MkdirAll
		newFile, err := os.Create(fileMeta.Location)
		if err != nil{
			fmt.Printf("Faliured to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		//把上传上来的文件数据通过I/O拷贝到新文件中 , 核心重点!
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil{
			fmt.Printf("Failured to copy file, error:%s\n", err.Error())
			return
		}

		//使用SHA1算法，生成hash值，在此之前把文件指针归为文件开头位置，因为后面要使用到copy函数
		newFile.Seek(0,0) //whence是相对位置，0表示从文件头开始，偏移offset个位置，然后开始读写文件
		fileMeta.FileSHA1 = util.FileSha1(newFile)

		//把文件元信息添加到meta.fileMetas中
		//meta.UpdateFileMeta(fileMeta)

		//把文件元信息存储到mysql数据库中，即tbl_file文件信息表中
		_ = meta.UpdateFileMetaToDB(fileMeta)

		//更新完DB中的tbl_file表，就要更新tbl_user_file用户文件信息表了，要知道是哪个用户上传了文件
		r.ParseForm()
		userName := r.Form.Get("username")
		ok := DBLayer.UpdateUserFileAtUploadFinished(userName, fileMeta.FileSHA1, fileMeta.FileName, fileMeta.FileSize)
		if ok {
			//先alert一下，再重定向到home页面
			//http.Redirect(w, r, "/user/info", http.StatusFound)
			w.Write([]byte("success"))
		}else{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}


//UploadFileSuc: 上传文件已完成的handler
func UploadFileSucHandler(w http.ResponseWriter, r *http.Request){
	io.WriteString(w, "Upload file success!")
}


//GetFileMetaInfo: 通过filehash得到文件的元信息
func GetFileMetaInfoHandler(w http.ResponseWriter, r *http.Request) {
	//从URL或者表单中解析出参数数据
	r.ParseForm() //解析后的结果存储在r.Form中

	//r.Form的类型是url.Values结构体，Values类型的最终基础类型是map[string][]string
	filehash := r.Form.Get("filehash")  //也可以这样写r.Form.Get("filehash")，默认取数组中的第一个元素
	//fileMeta := meta.GetFileMeta(filehash)
	fileMeta, err := meta.GetFileMetaFromDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("/handler/file_handler.go 101row is error: ", err.Error())
		return
	}

	//得到文件元信息结构体后，变为json格式返回到client端
	//即将struct序列化的过程
	jsonBytes,err :=json.Marshal(*fileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)  //向client返回HTTP状态码
		return
	}

	w.Write(jsonBytes)  //返回给client的Json字符串数据
}


//QueryFileMetas: 通过Get方法获取批量文件元数据信息
func QueryFileMetasHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm() //解析URL或者form表单的数据，结构存储在r.Form结构体中
	limitCount,_ := strconv.Atoi(r.Form.Get("limit")) //string转换为int
	fMetas := meta.QueryFileMetas(limitCount)

	//把fMetas数组，序列化为json对象数组
	jsonBytesArr,err := json.Marshal(fMetas)
	if err != nil{
		fmt.Printf("整体序列化 fMetas时发生错误, error:%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonBytesArr) //把序列化后的json对象数组以string的形式输出回client
}

//DownloadFile: 用于下载文件，所选的文件就是通过前端传来的file_hash来识别的
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()  //从request对象中解析form表单中的数据
	fSha1 := r.Form.Get("filehash") //结果存储在Form对象，Form类型是url.Values -> map[string][]string
	fMeta,err := meta.GetFileMetaFromDB(fSha1)
	if err != nil{
		log.Fatal("GetFileMetaFromDB is error: " + err.Error())
		return
	}

	fd, err := os.Open(fMeta.Location)
	fmt.Println("file_path:", fMeta.Location)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//设置响应头信息的头字段，用于告诉浏览器相关的信息
	w.Header().Set("Content-type","application/octet-stream") //设置MIME类型为下载文件类型
	//指明内容以什么方式返回给client，=inline内联，即在线展现；=attachment;filename...是以附件下载保存到本地的方式返回给客户端
	w.Header().Set("Content-Disposition","attachment; filename=" + fMeta.FileName)
	w.Write(data) //返回字节流信息给client
}


//UpdateFileMetaInfoHandler: 更新文件名，更新元信息和磁盘文件名
func UpdateFileMetaInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fSha1 := r.Form.Get("filehash") //文件SHA1的哈希值
	opType := r.Form.Get("optype") //操作类型，目前只有0=修改文件名
	newFileName := r.Form.Get("newfname") //新文件名

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden) //403禁止
		return
	}
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed) //不为PUT方法，则返回405，即请求方法不允许，不符合REST风格接口
		return
	}

	//从metas数组中查询文件元信息，并修改filename和location
	curFileMeta := meta.GetFileMeta(fSha1)
	oldLocation := curFileMeta.Location
	curFileMeta.FileName = newFileName
	curFileMeta.Location = Upload_File_Path_Prefix + newFileName
	meta.UpdateFileMeta(curFileMeta)

	//在磁盘上更改文件名
	err := os.Rename(oldLocation,curFileMeta.Location)
	if err != nil{
		fmt.Printf("os.Rename, error is:%s", err.Error())
	}

	//把新的文件元信息结构的json对象序列化为字节流
	dataBytes, err := json.Marshal(curFileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(dataBytes)

}

//DeleteFile: 通过文件hash删除文件元信息，以及从磁盘上删除，即硬件上删除
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fSha1)

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed) //405，请求方法不允许
		return
	}

	//先从磁盘上删除
	os.Remove(fMeta.Location)

	//再删除内存数组中的元素
	meta.DeleteFileMeta(fSha1)
}




















