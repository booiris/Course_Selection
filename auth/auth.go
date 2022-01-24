package auth

import (
	"course_selection/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	username := c.PostForm("Username")
	passwd := c.PostForm("Passwd")
	fmt.Println(username, passwd)
	result := database.Find_data(username, passwd)
	if result {
		c.String(http.StatusOK, "OK")
	} else {
		c.String(http.StatusOK, "NO")
	}
}

func Logout(c *gin.Context) {

}
func Whoami(c *gin.Context) {

}
