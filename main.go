package main

import (
	"course_selection/database"
	"course_selection/globals"
	"course_selection/router"
	"course_selection/types"

	"github.com/gin-gonic/gin"
)

func main() {

	// 连接数据库
	database.InitDB()
	database.InitRedis()

	database.Db.AutoMigrate(&types.Member{}, &types.Course{}, &types.SCourse{})

	// 初始化网络服务
	// gin.SetMode(gin.ReleaseMode)
	globals.G = gin.Default()
	router.RegisterRouter(globals.G)

	err := globals.G.Run(":1319")
	if err != nil {
		return
	}

}
