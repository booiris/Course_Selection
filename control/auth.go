package control

import (
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var data types.LoginRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	var res types.LoginResponse
	database.Db.Table("members").Where(&data).Find(&res.Data)
	if res == (types.LoginResponse{}) {
		c.JSON(http.StatusOK, types.LoginResponse{Code: types.WrongPassword})
	} else {
		c.SetCookie("camp-session", res.Data.UserID, 3600, "/", "", false, true)
		c.JSON(http.StatusOK, types.LoginResponse{Code: types.OK, Data: res.Data})
	}
}

func Logout(c *gin.Context) {
	value, err := c.Cookie("camp-session")
	if err != nil {
		log.Println(err)
		return
	}
	c.SetCookie("camp-session", value, -1, "/", "", false, true)
	c.JSON(http.StatusOK, types.LogoutResponse{Code: types.OK})
}

func Whoami(c *gin.Context) {
	value, err := c.Cookie("camp-session")
	if err != nil {
		log.Println(err)
		return
	}
	var res types.TMember
	database.Db.Table("members").Where(&value).Find(&res)
	c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.OK, Data: res})
}
