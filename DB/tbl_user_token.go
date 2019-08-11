package DB

import (
	MyDB "./Mysql"
	"database/sql"
	"log"
)


//GetTokenFromDB: 把token传入数据库进行比较，在数据库层面进行比较，若存在则返回token值就行了，且为true
func GetTokenFromDB(token string) bool {
	stmtQ,err := MyDB.DBConn().Prepare("select user_token from tbl_user_token where user_token = ? limit 1;")
	if err != nil{
		log.Fatal(err)
		return false
	}
	defer stmtQ.Close()

	var tmpToken sql.NullString
	err = stmtQ.QueryRow(token).Scan(&tmpToken)
	if err != nil {
		log.Fatal(err)
		return false
	}

	if !tmpToken.Valid {  // Valid is true if String is not NULL
		return false  //从数据库中查询的token为null，说明token不一样或不存在，token已失效
	}

	return true
}
