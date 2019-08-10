package DB

//此文件用于操作数据库中的fileserver数据库，具体是操作tbl_file表，对应的实体是文件的元信息结构体

import (
	MyDB "./Mysql"
	"database/sql"
	"fmt"
)

//UploadFileSucFinished: 把文件元信息持久化到Mysql数据库
func UploadFileSucFinished(fSha1, fName, fPath string , fSize int64) bool {

	//预编译SQL,可以有效预防SQL注入
	stmtIns, err := MyDB.DBConn().Prepare("insert into tbl_file(`file_sha1`,`file_path`,`file_name`,`file_size`,`status`) values (?,?,?,?,1);")
	if err != nil{
		fmt.Printf("DB-Prepare() is error:%s", err.Error())
		return false
	}
	defer stmtIns.Close()

	//Result接口是对已经执行SQL的反馈总结
	ret, err := stmtIns.Exec(fSha1, fName, fPath, fSize)
	if err != nil{
		fmt.Printf("stmtIn.Exec is error: %s\n", err.Error())
		return false
	}

	//ret的类型是Result接口类型，RowsAffected()返回成功插入的行数
	if rf, err := ret.RowsAffected(); err == nil { //当insert/update/delete没有错误的时候，说SQL的执行成功的
		if rf <= 0 { //rf<=0，说明插入的某条数据行已经存在过相同的数据，插入无效，比如fileHash等
			fmt.Println("rf <= 0, fileHash has been inserted before.")
		}
		return true
	}
	return false

}


//meta包不能使用，因为循环导入错误
type TableFileMeta struct {
	FileHash sql.NullString
	FileName sql.NullString
	FileSize sql.NullInt64
	FilePath sql.NullString
}

//GetFileMetaFromDB: 通过fSha1查找文件元信息
func GetFileMetaFromDB(fSha1 string) (*TableFileMeta, error) {
	stmtQuery, err := MyDB.DBConn().Prepare("SELECT `file_sha1`, `file_name`, `file_path`, `file_size`" +
		"FROM tbl_file WHERE `file_sha1` = ? AND status = 1 LIMIT 1;")
	if err != nil{
		fmt.Println("failed to select ,error: " + err.Error())
		return &TableFileMeta{}, err
	}
	defer stmtQuery.Close()

	fMeta := TableFileMeta{}
	//QueryRow(arg)是Prepare中的问号
	err = stmtQuery.QueryRow(fSha1).Scan(&fMeta.FileHash, &fMeta.FileName, &fMeta.FilePath, &fMeta.FileSize)
	if !fMeta.FileHash.Valid { //Valid is true if string not null, else Valid is false
		fmt.Println("fMeta.FileHash is nil !!!, file path = /DB/tbl_file.go 56rows")
	}
	if err != nil{
		fmt.Println("/DB/tbl_file.go 58row is error: " + err.Error())
		return &TableFileMeta{}, err
	}

	return &fMeta, nil //数据量大的结构体都用引用传递
}
//e21414b36bc2bde69020e2aa202ca8fe9a51f7fa