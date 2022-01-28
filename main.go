package main

import (
	"Course_Selection/auth"
	"Course_Selection/dataBase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	dataBase.CreateConnection()
	g := gin.Default()
	g.Handle("GET", "/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong, method is GET")
	})
	g.Handle("POST", "/ping/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong, method is POST")
	})
	g.Handle("GET", "/say_hello", func(c *gin.Context) {
		name := c.Query("name")
		c.String(http.StatusOK, fmt.Sprintf("hello %s", name))
	})

	g.POST("/auth/say", auth.Login)
	RegisterRouter(g)
	err := g.Run(":1319")
	if err != nil {
		return
	}
}
