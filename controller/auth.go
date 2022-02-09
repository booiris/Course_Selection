package controller

import (
	"course_selection/database"
	"course_selection/types"

	"github.com/gin-gonic/gin"
)

// 登录
func Login(c *gin.Context) {
	var loginRequest types.LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		c.JSON(200, types.LoginResponse{Code: types.ParamInvalid})
		return
	}
	var member types.Member
	database.Db.Where("username = ?", loginRequest.Username).Find(&member)
	if member.Password != loginRequest.Password {
		c.JSON(200, types.LoginResponse{Code: types.WrongPassword})
		return
	}
	c.SetCookie("camp-session", member.UserID, 3600, "/", "", false, true)
	c.JSON(200, types.LoginResponse{Code: types.OK, Data: struct{ UserID string }{UserID: member.UserID}})
}

// 登出
func Logout(c *gin.Context) {
	userid, err := c.Cookie("camp-session")
	if err != nil {
		c.JSON(200, types.LogoutResponse{Code: types.LoginRequired})
		return
	}
	c.SetCookie("camp-session", userid, -1, "/", "", false, false)
	c.JSON(200, types.LogoutResponse{Code: types.OK})
}

// 获取成员信息
func Whoami(c *gin.Context) {
	userid, err := c.Cookie("camp-session")
	if err != nil {
		c.JSON(200, types.WhoAmIResponse{Code: types.LoginRequired})
		return
	}
	var member types.Member
	database.Db.Where("user_id = ?", userid).Find(&member)
	var data types.TMember
	data.UserID = member.UserID
	data.Nickname = member.Nickname
	data.Username = member.Username
	data.UserType = member.UserType
	c.JSON(200, types.WhoAmIResponse{Code: types.OK, Data: data})
}
