package controller

import (
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

func check_permission(c *gin.Context) bool {
	value, err := c.Cookie("camp-session")
	if err != nil {
		log.Println(err)
		return false
	}
	var res struct{ UserType types.UserType }
	database.Db.Table("members").Where(&value).Find(&res)
	return res.UserType == types.Admin
}

func check_param(data types.CreateMemberRequest) bool {

	// 用户昵称
	len := strings.Count(data.Nickname, "")
	if len < 4 || len > 20 {
		return false
	}

	// 用户名
	for _, v := range data.Username {
		if !(unicode.IsLetter(v) || unicode.IsDigit(v)) {
			return false
		}
	}
	len = strings.Count(data.Username, "")
	if len < 8 || len > 20 {
		return false
	}

	// 密码
	letter, digit := false, false
	for _, v := range data.Password {
		if !(unicode.IsLetter(v) || unicode.IsDigit(v)) {
			return false
		}
		if unicode.IsLetter(v) {
			letter = true
		} else if unicode.IsDigit(v) {
			digit = true
		}
	}
	if !(letter && digit) {
		return false
	}
	len = strings.Count(data.Password, "")
	if len < 8 || len > 20 {
		return false
	}

	// 用户类型
	if data.UserType < 1 || data.UserType > 3 {
		return false
	}

	return true
}

// 创建成员
func Member_create(c *gin.Context) {
	if !check_permission(c) {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.PermDenied})
		return
	}
	var data types.CreateMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	if check_param(data) {
		var res struct{ UserId string }
		database.Db.Model(types.Member{}).Where("username=?", data.Username).Find(&res)
		if res == (struct{ UserId string }{}) {
			database.Db.Model(types.Member{}).Create(&data)
			database.Db.Model(types.Member{}).Where("username=?", data.Username).Find(&res)
			c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.OK, Data: struct{ UserID string }{res.UserId}})
		} else {
			c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.UserHasExisted})
		}
	} else {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.ParamInvalid})
	}
}

// 获取单个成员信息
func Member_get(c *gin.Context) {
}

// 获取成员列表
func Member_get_list(c *gin.Context) {
	if check_permission(c) {
		c.JSON(http.StatusOK, types.GetMemberListResponse{Code: types.PermDenied})
		return
	}
	var data types.GetMemberListRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	var users []types.TMember
	if err := database.Db.Table("members").Limit(data.Limit).Offset(data.Offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, types.GetMemberListResponse{Code: types.UnknownError})
	}
	c.JSON(http.StatusOK, types.GetMemberListResponse{Code: types.OK, Data: struct{ MemberList []types.TMember }{MemberList: users}})
}

// 更新成员
func Member_update(c *gin.Context) {
	var data types.UpdateMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	res := database.Db.Model(types.Member{}).Where("user_id=?", data.UserID).Update("Nickname", data.Nickname)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusOK, types.UpdateMemberResponse{Code: types.UserNotExisted})
	} else {
		c.JSON(http.StatusOK, types.UpdateMemberResponse{Code: types.OK})
	}
}

// 删除成员
func Member_delete(c *gin.Context) {
	var data types.DeleteMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	var user types.Member
	database.Db.Model(&types.Member{}).Where(&data).Find(&user)
	if user == (types.Member{}) {
		c.JSON(http.StatusOK, types.DeleteMemberResponse{Code: types.UserNotExisted})
		return
	}
	if user.UserType == types.Teacher {
		database.Db.Table("courses").Where("teacher_id=?", user.UserID).Update("teacher_id", nil)
	} else if user.UserType == types.Student {
		database.Db.Where("user_id=?", user.UserID).Delete(types.SCourse{})
	}
	database.Db.Where("user_id=?", data.UserID).Delete(&types.Member{})
	c.SetCookie("camp-session", data.UserID, -1, "/", "", false, true)
	c.JSON(http.StatusOK, types.DeleteMemberResponse{Code: types.OK})
}
