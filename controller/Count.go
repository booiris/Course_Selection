package controller

import (
	"course_selection/types"

	"github.com/gin-gonic/gin"
)

func Count(c *gin.Context) {
	println(types.Count)
	types.Count -= 1
}
