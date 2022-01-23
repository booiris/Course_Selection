package main

import (
	"course_selection/types"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	types.RegisterRouter(g)
}
