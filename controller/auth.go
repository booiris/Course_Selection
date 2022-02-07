package controller

import (
	"course_selection/database"
	"course_selection/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 登录
func Login(c *gin.Context) {
	var loginRequest types.LoginRequest
	//表单数据绑定
	if err := c.ShouldBind(&loginRequest); err != nil {
		c.JSON(http.StatusOK, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	//查询
	var user types.Member
	database.Db.Where("username=?", loginRequest.Username).First(&user)
	if user == (types.Member{}) {
		c.JSON(http.StatusOK, types.LoginResponse{Code: types.WrongPassword})
	} else {
		if user.Password != loginRequest.Password {
			c.JSON(http.StatusOK, types.LoginResponse{Code: types.WrongPassword})
		} else {
			// 登录成功设置cookie
			value := loginRequest.Username
			c.SetCookie("camp-session", value, 60*60*60*24, "/", "localhost", false, true)
			c.JSON(http.StatusOK, types.LoginResponse{
				Code: types.OK,
				Data: struct {
					UserID string
				}{user.UserID},
			})
		}
	}
}

// 登出
func Logout(c *gin.Context) {
	//登出删除cookie
	if value, err := c.Cookie("camp-session"); err != nil {
		c.SetCookie("camp-session", value, -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, types.LogoutResponse{Code: types.LoginRequired})
		return
	}
	c.JSON(http.StatusOK, types.LogoutResponse{Code: types.OK})
}

// 获取成员信息
func Whoami(c *gin.Context) {
	value, err := c.Cookie("camp-session")
	if err != nil {
		c.JSON(http.StatusOK, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	var user types.Member
	database.Db.Where("username=?", value).First(&user)

	c.JSON(http.StatusOK, types.WhoAmIResponse{
		types.OK,
		types.TMember{
			UserID:   user.UserID,
			Username: user.Username,
			UserType: user.UserType,
			Nickname: user.Nickname},
	})
}
