package main

import (
	"course_selection/database"
	"course_selection/router"

	"github.com/gin-gonic/gin"
)

func main() {

	// 连接数据库
	database.InitDB()

	// database.Db.AutoMigrate(&types.Member{})

	// 初始化网络服务
	g := gin.Default()
	router.RegisterRouter(g)
	g.Run(":2000")

}
