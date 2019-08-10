package DB

import (
	MyDB "./Mysql"
	"database/sql"
	"fmt"
)

//UserSignUp: 用户注册，包括检查重复注册
func UserSignUpToDB(username , password, email string) bool {
	stmtIns, err := MyDB.DBConn().Prepare("INSERT INTO tbl_user (`user_name`,`user_pwd`, `email`) VALUES (?,?,?);")
	if err != nil{
		fmt.Println("Failed to insert tbl_user, error: " + err.Error())
		return false
	}
	defer stmtIns.Close()

	ret, err := stmtIns.Exec(username,password, email)
	if err != nil{
		fmt.Println("Failed to stmtIns.Exec ,error: " + err.Error())
		return false
	}

	if afRows,err := ret.RowsAffected(); err == nil && afRows <= 0 {
		fmt.Println("Failed to insert, because insert_data has ready")
		return false
	}

	return true
}

//UserLoginInFromDB: 从数据库中查询用户数据
func UserLoginInFromDB(username, password string) bool {
	stmtQuery, err := MyDB.DBConn().Prepare("select user_name, user_pwd from tbl_user where user_name = ? LIMIT 1;")
	if err != nil{
		fmt.Println("Failed to select, error: " + err.Error())
		return false
	}
	defer stmtQuery.Close()

	var userPwd sql.NullString
	var userName sql.NullString
	err = stmtQuery.QueryRow(username).Scan(&userName, &userPwd)
	if err != nil {
		fmt.Println("Failed to query, error: "+err.Error())
		return false
	}
	if userPwd.Valid && userPwd.String == password && userName.String == username{ //字段为null
		return true
	}else { //userPwd.Valid == false, 即返回为nil
		fmt.Println("username not found !")
		return false
	}

}