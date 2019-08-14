package handler

//multipart_upload_handler.go文件，用于分块上传接口
import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	MyRedis "../cache/redis"
	"../util"
)

const (
	UnitBlockSize = 5 * 1024 * 1025  //切分文件时的文件块大小
	ChunkFilePrefix = "MP_"
	Redis_hset= "hset"  //hash结构的hset
)

//InitMultipartUploadHandler: 初始化文件块信息，并把文件块元信息返回给客户端
func InitMultipartUploadHandler(w http.ResponseWriter, r *http.Request){
	//1.解析用户请求信息
	r.ParseForm()
	userName := r.Form.Get("username")
	fileHash := r.Form.Get("filehash")
	fileName := r.Form.Get("filename")
	fileSize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil{
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//2.初始化文件块信息结构体
	mBlock := &MultipartUploadInfo{
		fileName:fileName,
		fileHash:fileHash,
		chunkSize:UnitBlockSize,  //5MB, 最基本单位是B
		uploadId:userName + fmt.Sprintf("%x", time.Now().UnixNano()), //纳秒级别的时间戳
		chunkCount: int(math.Ceil(float64(fileSize)/(UnitBlockSize))), //文件切分为多少文件块，向上取整
	}

	//3.从redis连接池中申请一个连接
	connPool := MyRedis.RedisPool().Get() //得到redis连接池
	defer connPool.Close()

	//4.把文件块元信息存储在redis引擎中
	//使用Hash数据结构，是一个map-map形式,key -> (field, value)
	//hmset是一次性多次写入redis引擎
	connPool.Do(Redis_hset, ChunkFilePrefix + mBlock.uploadId, "fileName", mBlock.fileName)
	connPool.Do(Redis_hset, ChunkFilePrefix + mBlock.uploadId, "fileHash", mBlock.fileHash)
	connPool.Do(Redis_hset, ChunkFilePrefix + mBlock.uploadId, "chunkSize", mBlock.chunkSize)
	connPool.Do(Redis_hset, ChunkFilePrefix + mBlock.uploadId, "chunkCount", mBlock.chunkCount)

	//4.把文件块元数据返回给客户端
	w.Write(util.RespMsg{200, "OK", mBlock}.JSONBytes())
}


//UploadPartHandler: 上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1.获取用户请求参数
	r.ParseForm()
	uploadId := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("chunkindex")

	//2.向redis连接池申请一个可用连接
	connPool := MyRedis.RedisPool().Get()
	defer connPool.Close()

	//3.创建一个文件，用于存储分块文件信息，并从request.body中读取相应的数据
	err := os.MkdirAll("./data/" + uploadId,0666)
	if err != nil{
		log.Fatal(err)
		return
	}
	//fd是chunkIndex文件的可用I/O，就相当于连接磁盘扇区的一个管道
	fd, err := os.Create("./data/" + uploadId + "/" + chunkIndex)
	if err != nil{
		w.Write(util.RespMsg{-1,"upload part failed",nil}.JSONBytes())
		log.Fatal(err)
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024) //每次读取1MB
	for{
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil { //发生了错误或者读完了
			break
		}
	}

	//4.更新redis的缓存状态，用于记录文件块是否上传成功

	//把chunkindex设置为1，就表示此文件块已经上传成功了
	connPool.Do(Redis_hset, ChunkFilePrefix + uploadId, "chkidx_"+ chunkIndex, 1)

	//5.把成功的响应信息返回给客户端
	w.Write(util.RespMsg{200,"OK",nil}.JSONBytes())
}



