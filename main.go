package main

import (
	"course_selection/database"
	"course_selection/router"
	"course_selection/types"
	"github.com/gin-gonic/gin"
)

//180.184.74.221
func main() {

	// 连接数据库
	database.InitDB()
	database.InitRedis()
	types.CurrentSoldOutMap = make(map[string]bool)
	// 初始化网络服务
	g := gin.Default()
	router.RegisterRouter(g)
	g.Run(":2222")

}
