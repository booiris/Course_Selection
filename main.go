package main

import (
	"course_selection/types"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	types.RegisterRouter(g)

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
