package DB

import (
	MyDB "./Mysql"
	"database/sql"
	"fmt"
	"log"
)

//UserSignUp: 用户注册，包括检查重复注册
func UserSignUpToDB(username , password, email, phone string) bool {
	//先检查username、email、phone是否被人注册过，可使用Result.RowsAffected().int64来检查成功插入了几行
	//逻辑：使用or查询出所有的username或者email或者phone，若三者为空，则可以注册，否则返回username或者email、或者phone重复
	//-- 情况1：三个参数都没有在表中存在相应的记录，此时可以进一步注册
	//-- 情况2：三个参数中存在一个以上与表中的记录相同，此时注册失败
	//stmtQuery, err := MyDB.DBConn().Prepare("SELECT user_name, email, phone FROM tbl_user WHERE user_name = ? OR email = ? OR phone = ? limit 1;")
	//if err != nil{
	//	fmt.Println("select ... , error: " + err.Error())
	//	return false
	//}
	//var _username sql.NullString
	//var _email sql.NullString
	//var _phone sql.NullString
	//err = stmtQuery.QueryRow(username, email, phone).Scan(&_username, &_email, &_phone)
	//if err != nil{
	//	fmt.Println(" queryRow ... , error: " + err.Error())
	//	return false
	//}
	//if _username.Valid || _email.Valid || _phone.Valid { //这三者只要有一者存在于数据表中，传入参数重叠，注册失败
	//	return false
	//}

	stmtIns, err := MyDB.DBConn().Prepare("INSERT INTO tbl_user (`user_name`,`user_pwd`, `email`, `phone`) VALUES (?,?,?,?);")
	if err != nil{
		fmt.Println("Failed to insert tbl_user, error: " + err.Error())
		return false
	}
	defer stmtIns.Close()

	ret, err := stmtIns.Exec(username,password, email, phone)
	if err != nil{
		fmt.Println("Failed to stmtIns.Exec ,error: " + err.Error())
		return false
	}

	//插入失败，肯定是三个参数至少有一个已经存在于数据表中了
	if afRows,err := ret.RowsAffected(); err == nil && afRows <= 0 {
		fmt.Println("Failed to insert, because insert data has ready")
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
	if userPwd.Valid && userPwd.String == password && userName.String == username{
		return true
	}else { //userPwd.Valid == false, 即返回为nil
		fmt.Println("username not found !")
		return false
	}

}


//UpdateUserToken: 刷新用户登录的token，同一个用户新旧token可覆盖replace
func UpdateUserToken(username , token string) bool {
	stmtReplace, err := MyDB.DBConn().Prepare("replace into tbl_user_token (`user_name`, `user_token`) values (?,?);")
	if err != nil{
		fmt.Println("prepare is error: " + err.Error())
		return false
	}
	defer stmtReplace.Close()

	_, err = stmtReplace.Exec(username, token)
	if err != nil {
		fmt.Println("stmt.Exec is error: " + err.Error())
		return false
	}

	return true
}


//GetUserInfo: 用于获取用户信息
func GetUserInfo(username string)(*User, error){
	user := &User{}

	stmtQuery, err := MyDB.DBConn().Prepare("SELECT user_name, signup_at FROM tbl_user WHERE user_name = ? LIMIT 1;")
	if err != nil{
		fmt.Println(err.Error())
		return user, err
	}
	defer stmtQuery.Close()

	//从结果集中把数据存入变量的内存地址中
	err = stmtQuery.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil{
		fmt.Println(err.Error())
		return user, err
	}

	return user, nil
}


//UpdateAvatarURL: 更新用户头像URL，默认URL是开发者们自己放置的一张默认图片
func UpdateAvatarURL(userName, avatarURL string) bool {

	stmtUpd, err := MyDB.DBConn().Prepare("UPDATE tbl_user AS u SET u.`avatar_img` = ? WHERE u.`user_name` = ?;")
	if err != nil {
		log.Fatal(err)
		return false
	}

	ret, err := stmtUpd.Exec(avatarURL, userName)
	if err != nil {
		log.Fatal(err)
		return false
	}

	if resRow, err := ret.RowsAffected(); err == nil && resRow <= 0 { //说明插入不成功
		log.Println("avatarImgURL，插入不成功，因为已经存在相同的记录数据...")
	}

	return true
}


//GetAvatarImgURL: 从数据库中获取用户的avatarURL
func GetAvatarImgURL(userName string) (string, bool) {
	stmtSelect, err := MyDB.DBConn().Prepare("SELECT u.`avatar_img` FROM tbl_user AS u WHERE u.`user_name` = ?;")
	if err != nil{
		log.Fatal(err)
		return "", false
	}

	var avatarImgURL sql.NullString

	err = stmtSelect.QueryRow(userName).Scan(&avatarImgURL)
	if err != nil{
		log.Fatal(err)
		return "", false
	}

	if !avatarImgURL.Valid {
		log.Fatal("avatarImgURL is null")
		return "", false
	}

	return avatarImgURL.String, true

}















