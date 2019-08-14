package handler


//MultipartUploadInfo: 分块上传的文件块信息
type MultipartUploadInfo struct{
	fileHash   string   //文件hash
	fileName   string   //文件名
	chunkSize  int      //块文件大小，以B作为基本单位
	uploadId    string   //文件块的id号，也是上传的ID号
	chunkCount int      //一个文件分为几块文件块
}


