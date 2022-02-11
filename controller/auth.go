package controller

import (
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type LoginRequest struct {
// 	Username string `form:"Username" json:"Username" xml:"Username"  binding:"required"`
// 	Password string `form:"Password" json:"Password" xml:"Password"  binding:"required"`
// }

// // 登录成功后需要 Set-Cookie("camp-session", ${value})
// // 密码错误范围密码错误状态码

// type LoginResponse struct {
// 	Code ErrNo
// 	Data struct {
// 		UserID string
// 	}
// }

func Login(c *gin.Context) {
	/* 解析json */
	var request types.LoginRequest
	if err := c.ShouldBind(&request); err != nil {
		log.Println(err)
		return
	}

	/* 提取数据 */
	member := types.Member{}
	database.Db.Debug().Model(&types.Member{}).Where("username = ?", request.Username).Find(&member)
	// fmt.Println("userId : ", member)
	// fmt.Println("request : ", request)

	/* 构建Response */
	var response types.LoginResponse
	/* 判断返回ErrNo */
	// UserHasDeleted     ErrNo = 3  // 用户已删除
	// UserNotExisted     ErrNo = 4  // 用户不存在
	// WrongPassword      ErrNo = 5  // 密码错误
	if member == (types.Member{}) {
		// 空， 不存在该用户
		response = types.LoginResponse{
			Code: types.UserNotExisted,
			Data: struct{ UserID string }{UserID: member.UserID},
		}
	} else if member.Password != request.Password {
		response = types.LoginResponse{
			Code: types.WrongPassword,
			Data: struct{ UserID string }{UserID: member.UserID},
		}
	} else {
		response = types.LoginResponse{
			Code: types.OK,
			Data: struct{ UserID string }{UserID: member.UserID},
		}
		c.SetCookie("camp-session", response.Data.UserID, 3600, "/", "", false, true)
	}
	c.JSON(http.StatusOK, response)
}

func Logout(c *gin.Context) {
	value, err := c.Cookie("camp-session")
	if err != nil {
		log.Println(err)
		// LoginRequired      ErrNo = 6  // 用户未登录
		c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.LoginRequired})
		return
	}
	c.SetCookie("camp-session", value, -1, "/", "", false, true)
	c.JSON(http.StatusOK, types.LogoutResponse{Code: types.OK})
}

func Whoami(c *gin.Context) {
	value, err := c.Cookie("camp-session")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.LoginRequired})
		return
	}
	var res types.TMember
	database.Db.Table("members").Where(&value).Find(&res)
	c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.OK, Data: res})
}
