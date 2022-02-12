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
	// TODO！！！！！！为了方便初始化数据库添加，记得删除
	database.Db.Exec("DROP TABLE courses")
	database.Db.Exec("DROP TABLE s_courses")

	database.Db.AutoMigrate(&types.Member{}, &types.Course{}, &types.SCourse{})

	database.InitRedis()

	// 初始化网络服务
	// gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	router.RegisterRouter(g)
	g.Run(":2000")

}
