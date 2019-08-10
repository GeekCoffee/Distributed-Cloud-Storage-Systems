package Mysql
//用于初始化连接mysql数据库对象

import(
	"database/sql" //导入sql包，至少需要再导入一个数据库的驱动
	"fmt"
	_ "github.com/go-sql-driver/mysql" //mysql数据库驱动程序，加入下划线，可以让驱动程序配置到sql包中
	"os"
)

var db *sql.DB

func init(){
	var err error
	//Open()用于检验参数是否有效，不做DB的连接操作，可以安全地被多个go协程调用，其自身会维护数据库连接池
	db, err = sql.Open("mysql", "root:abc5518988@tcp(localhost:3306)/fileserver?charset=utf8")
	if err != nil{
		fmt.Printf("open the mysql_db fail, error:%s\n", err.Error())
		os.Exit(-1)
	}

	//设置最大同时活跃连接数
	db.SetMaxOpenConns(1000)

	//测试连接数据库是否成功，这个是真的是连接了数据库
	err = db.Ping()
	if err != nil{
		fmt.Printf("db.Ping() is fail, error:%s\n", err.Error())
		os.Exit(-1)
	}
}

//DBConn : 返回已经连接数据库的对象引用
func DBConn() *sql.DB{
	return db
}














