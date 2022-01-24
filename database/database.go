package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//数据库配置
const (
	userName = "root"
	password = "bytedancecamp"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "booiris_user"
)

//Db数据库连接池
var database *sql.DB

func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")

	//打开数据库
	database, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	database.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	database.SetMaxIdleConns(10)
	//验证连接
	if err := database.Ping(); err != nil {
		fmt.Println("opon database fail")
		return
	}
	fmt.Println("connnect success")
}

func Find_data(values ...string) bool {
	sqlStr := "select username from login_data where username="
	sqlStr += "'" + values[0] + "' and passwd='"
	sqlStr += values[1] + "'"
	fmt.Println(sqlStr)
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	var temp string
	err := database.QueryRow(sqlStr).Scan(&temp)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return false
	}
	return true
}
