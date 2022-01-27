package main

import (
	"course_selection/database"
	"course_selection/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	database.InitDB()

	// 初始化网络服务
	g := gin.Default()
	router.RegisterRouter(g)
	g.Run(":2000")

	// g.Handle("GET", "/ping", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "pong, method is GET")
	// })
	// g.Handle("POST", "/ping", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "pong, method is POST")
	// })
	// g.Handle("GET", "/say_hello", func(c *gin.Context) {
	// 	name := c.Query("name")
	// 	c.String(http.StatusOK, fmt.Sprintf("hello %s", name))
	// })
	// g.Run(":2000")
}
