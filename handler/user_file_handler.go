package handler

//秒传的应用场景：用户文件上传、离线下载、好友分享等

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	DBLayer "../DB"
)

//GetUserFileInfoHandler: 验证token后，从数据库中取出文件信息
func GetUserFileInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r.ParseForm()
	userName := r.Form.Get("username")
	limitCount,_ := strconv.Atoi(r.Form.Get("limit"))

	data, err := DBLayer.GetUserFilesInfo(userName, limitCount)
	if err != nil{
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataBytes, err := json.Marshal(data)
	if err != nil{
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(dataBytes)
}
