package DB

import(
	MyDB "./Mysql"
	"fmt"
	"log"
	"time"
)

//UpdateUserFileAtUploadFinished: 当文件上传成功的时候，插入新的数据行, 插入成功返回true
func UpdateUserFileAtUploadFinished(userName, fileHash, fileName string, fileSize int64) bool {
	stmtIns, err := MyDB.DBConn().Prepare("insert into tbl_user_file (`user_name`,`file_sha1`,`file_name`,`file_size`" +
		",`upload_at`) values (?,?,?,?,?);")
	if err != nil{
		log.Fatal(err)
		return false
	}

	_, err = stmtIns.Exec(userName, fileHash, fileName, fileSize, time.Now()) //Time类型与数据库中的datetime类型相对应
	if err != nil{
		log.Fatal(err)
		return false
	}

	return true
}


func GetUserFilesInfo(userName string, limit int)([]UserFile, error) {
	stmtQ, err := MyDB.DBConn().Prepare(
		"select `user_name`, `file_sha1`, `file_size`, `file_name`, `upload_at`, `last_update`" +
		"from tbl_user_file where user_name = ? limit ?")
	if err != nil{
		fmt.Println(err.Error())
		return []UserFile{}, err
	}
	defer  stmtQ.Close()

	//查询多行记录使用Query，查询一行一般使用QueryRow
	rows, err := stmtQ.Query(userName, limit)
	if err != nil{
		log.Fatal(err)
		return []UserFile{}, err
	}

	var userFileInfos []UserFile
	for rows.Next() {  //从结果集从循环每一行
		userFile := UserFile{}
		err := rows.Scan(&userFile.UserName, &userFile.FileHash,
			&userFile.FileSize, &userFile.FileName, &userFile.UploadAt,&userFile.LastUpdated)
		//判空操作
		if err != nil{
			fmt.Println(err)
			break
		}

		userFileInfos = append(userFileInfos, userFile)
	}

	return userFileInfos, nil
}