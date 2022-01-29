package controller

import (
	"course_selection/database"
	"course_selection/types"
	"log"
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
}

// 获取单个成员信息
func Member_get(c *gin.Context) {
}

// 获取成员列表
func Member_get_list(c *gin.Context) {
}

// 更新成员
func Member_update(c *gin.Context) {
}

// 删除成员
func Member_delete(c *gin.Context) {
}
