package controller

import (
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 登录
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

// 登出
func Logout(c *gin.Context) {

}

// 获取成员信息
func Whoami(c *gin.Context) {
}
