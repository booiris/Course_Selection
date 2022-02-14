package main

import (
	"course_selection/database"
	"course_selection/router"
	"course_selection/types"

	"github.com/gin-gonic/gin"
)

func main() {

	// 连接数据库
	database.InitDB()

	database.Db.AutoMigrate(&types.Member{}, &types.Course{}, &types.SCourse{})

	database.InitRedis()

	// 初始化网络服务
	//gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	router.RegisterRouter(g)
	g.Run(":80")
}
