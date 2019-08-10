package meta

//filemate.go文件，用于记录文件元信息，并操作提供相关操作
import (
	myDB "../DB"
	"fmt"
	"sort"
)

//FileMeta: 文件元信息结构体
type FileMeta struct{
	FileSHA1   string  //文件的hash值，用于做文件的ID
	FileName   string  //文件的名字
	Location   string  //文件路径
	UploadTime string  //文件上传时间
	FileSize   int64   //文件大小
}

//fileMetas: 定义一个map全局变量，用于存储所有的FileMeta对象
var fileMetas map[string]FileMeta

//init: 包被调用的时候，init方法第一个执行，用于生成空的map对象
func init(){
	fileMetas = make(map[string]FileMeta)
}

//UploadFileMeta: 更新或者增加文件元信息到一个全局变量中
func UpdateFileMeta(fMeta FileMeta) {
	fileMetas[fMeta.FileSHA1] = fMeta  //若FID存在，则为修改，否则为在map中新添加一个key-value
}

//UpdateFileMetaToDB: 更新或增加文件元信息到mysql数据库中
func UpdateFileMetaToDB(fMeta FileMeta) bool {
	return myDB.UploadFileSucFinished(fMeta.FileSHA1, fMeta.FileName, fMeta.Location, fMeta.FileSize)
}

//GetFileMeta: 通过文件的SHA1，获取一个文件元信息结构体的数据
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}


func GetFileMetaFromDB(fileSha1 string) (*FileMeta,error) {
	tblFileMe := &myDB.TableFileMeta{}
	tblFileMe, err := myDB.GetFileMetaFromDB(fileSha1)
	if err != nil {
		fmt.Println("/meta/filemeta.go 45row is error: " + err.Error())
		return &FileMeta{}, err
	}
	fMeta := &FileMeta{
		FileSHA1:tblFileMe.FileHash.String, //sql.NullString的结构体中有String类型
		FileName:tblFileMe.FileName.String,
		Location:tblFileMe.FilePath.String,
		FileSize:tblFileMe.FileSize.Int64,
	}

	return fMeta, nil
}

//QueryFileMetas: 按照时间排序进行查询limitCount个FileMeta的元信息组
func QueryFileMetas(limitCount int) []FileMeta{
	fMetaArr := make([]FileMeta, len(fileMetas))  //生成一个新的fileMeta数组
	for _, v := range fileMetas{ //遍历map，将fileMeta逐个存储到fileMeta数组中
		fMetaArr = append(fMetaArr, v)
	}

	//对fMetaArr文件元信息数组进行排序
	sort.Sort(ByTimeSort(fMetaArr))  //调用Sort函数，传入相应的数组就行了，底层会自动多次调用Less()等方法的
	return fMetaArr[0:limitCount]  //从0开始，返回limitCount个元素
}

//DeleteFileMeta: 通过hash值删除文件元信息，即从metas数组中剔除该元信息
func DeleteFileMeta(fSha1 string){
	delete(fileMetas, fSha1)
}
