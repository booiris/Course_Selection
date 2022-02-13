package controller

import (
	"context"
	"course_selection/database"
	"course_selection/types"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
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
	sessionKey := uuid.NewV4().String()
	ctx := context.Background()
	database.Rdb.Set(ctx, sessionKey, member.UserID, time.Second*3600)
	c.SetCookie("camp-session", sessionKey, 3600, "/", "", false, true)
	c.JSON(200, types.LoginResponse{Code: types.OK, Data: struct{ UserID string }{UserID: member.UserID}})
}

// 登出
func Logout(c *gin.Context) {
	sessionKey, err := c.Cookie("camp-session")
	if err != nil {
		c.JSON(200, types.LogoutResponse{Code: types.LoginRequired})
		return
	}
	ctx := context.Background()
	database.Rdb.Del(ctx, sessionKey)
	c.SetCookie("camp-session", sessionKey, -1, "/", "", false, true)
	c.JSON(200, types.LogoutResponse{Code: types.OK})
}

// 获取成员信息
func Whoami(c *gin.Context) {
	sessionKey, err := c.Cookie("camp-session")
	if err != nil {
		c.JSON(200, types.WhoAmIResponse{Code: types.LoginRequired})
		return
	}
	ctx := context.Background()
	userid := database.Rdb.Get(ctx, sessionKey)
	if userid.Err() == redis.Nil {
		c.JSON(200, types.WhoAmIResponse{Code: types.LoginRequired})
		return
	}
	var member types.Member
	database.Db.Where("user_id = ?", userid.Val()).Find(&member)
	var data types.TMember
	data.UserID = member.UserID
	data.Nickname = member.Nickname
	data.Username = member.Username
	data.UserType = member.UserType
	c.JSON(200, types.WhoAmIResponse{Code: types.OK, Data: data})
}
