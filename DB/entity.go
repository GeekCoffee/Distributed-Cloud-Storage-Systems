package DB

import "database/sql"

//meta包不能使用，因为循环导入错误
//tbl_file: 文件元信息表
type TableFileMeta struct {
	FileHash sql.NullString  //在string类型上，封装了一层，可用Valid变量进行判空操作
	FileName sql.NullString
	FileSize sql.NullInt64
	FilePath sql.NullString
}

//tbl_user: 用户信息表
type User struct {
	Username 	 string
	Email 		 string
	Phone		 string
	SignupAt     string
	LastActiveAt string //最后活跃时间
	Status   	 int  //此用户状态
}

//tbl_user_token: 用户token表
type Token struct{
	Username string
	Token    string
}


//tbl_user_file: 用户文件信息表
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string //上传时间
	LastUpdated string //表最后修改时间
}

