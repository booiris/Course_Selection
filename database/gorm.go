package database

import (
	"course_selection/types"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

//数据库配置
const (
	userName = "root"
	password = "1234"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "main"
)

var Db *gorm.DB

func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8mb4"
	dsn := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8mb4"}, "")
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	Db.AutoMigrate(&types.Member{})
	Db.AutoMigrate(&types.Course{})
	Db.AutoMigrate(&types.SCourse{})
	if err != nil {
		fmt.Println("open database fail")
		return
	}
	fmt.Println("connnect success")
}
