package controller

import (
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

// type GetMemberRequest struct {
// 	UserID string `form:"UserID" json:"UserID" xml:"UserID"  binding:"required"`
// }
// type GetMemberResponse struct {
// 	Code ErrNo
// 	Data TMember
// }

func check_permission(c *gin.Context) bool {
	value, err := c.Cookie("camp-session")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, types.WhoAmIResponse{Code: types.LoginRequired})
		return false
	}
	var res struct{ UserType types.UserType }
	database.Db.Table("members").Where(&value).Find(&res)
	return res.UserType == types.Admin
}
func check_param(data types.CreateMemberRequest) bool {

	// 用户昵称
	lens := strings.Count(data.Nickname, "")
	if lens < 4 || lens > 20 {
		return false
	}

	// 用户名
	for _, v := range data.Username {
		if !(unicode.IsLetter(v) || unicode.IsDigit(v)) {
			return false
		}
	}
	lens = strings.Count(data.Username, "")
	if lens < 8 || lens > 20 {
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
	lens = strings.Count(data.Password, "")
	if lens < 8 || lens > 20 {
		return false
	}

	// 用户类型
	if data.UserType < 1 || data.UserType > 3 {
		return false
	}

	return true
}
func check_state_byName(username string) types.ErrNo {
	var res types.Member
	database.Db.Table("members").Unscoped().Where("username = ?", username).Find(&res)
	if res == (types.Member{}) {
		// 无记录
		return types.UserNotExisted
	} else {
		if res.Deleted.Time != (time.Time{}) {
			// 存在
			return types.UserHasExisted
		} else {
			// 软删
			return types.UserHasDeleted
		}
	}
}
func check_state_byID(userId string) types.ErrNo {
	var res types.Member
	database.Db.Table("members").Unscoped().Where("user_id = ?", userId).Find(&res)
	if res == (types.Member{}) {
		// 无记录
		return types.UserNotExisted
	} else {
		if res.Deleted.Time != (time.Time{}) {
			// 存在
			return types.UserHasExisted
		} else {
			// 软删
			return types.UserHasDeleted
		}
	}
}
func Member_get(c *gin.Context) {
	if !check_permission(c) {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.PermDenied})
		return
	}

	/* 解析json */
	var data types.GetCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}

	var res types.GetMemberResponse
	database.Db.Table("members").Where(&data).Find(&res.Data)
	if res.Data != (types.TMember{}) {
		res.Code = types.OK
	} else {
		// UserNotExisted     ErrNo = 4  // 用户不存在
		res.Code = types.UserNotExisted
	}
	c.JSON(http.StatusOK, res)
}

// ???
func Member_get_list(c *gin.Context) {
}

// type CreateMemberRequest struct {
// 	Nickname string   `form:"Nickname" json:"Nickname" xml:"Nickname"  binding:"required"` // required，不小于 4 位 不超过 20 位
// 	Username string   `form:"Username" json:"Username" xml:"Username"  binding:"required"` // required，只支持大小写，长度不小于 8 位 不超过 20 位
// 	Password string   `form:"Password" json:"Password" xml:"Password"  binding:"required"` // required，同时包括大小写、数字，长度不少于 8 位 不超过 20 位
// 	UserType UserType `form:"UserType" json:"UserType" xml:"UserType"  binding:"required"` // required, 枚举值
// }

// type CreateMemberResponse struct {
// 	Code ErrNo
// 	Data struct {
// 		UserID string // int64 范围
// 	}
// }

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

	var res types.CreateMemberResponse

	if check_param(data) {

		state := check_state_byName(data.Username)
		if state == types.UserNotExisted {
			// 无记录
			res.Code = types.OK
			database.Db.Table("members").Create(&data)
			database.Db.Table("members").Where(&data).Find(&res.Data)
		} else {
			res.Code = state
		}
	} else {
		res.Code = types.ParamInvalid
	}
	c.JSON(http.StatusOK, res)
}

// type UpdateMemberRequest struct {
// 	UserID   string `form:"UserID" json:"UserID" xml:"UserID"  binding:"required"`
// 	Nickname string `form:"Nickname" json:"Nickname" xml:"Nickname"  binding:"required"`
// }

// type UpdateMemberResponse struct {
// 	Code ErrNo
// }

func Member_update(c *gin.Context) {
	if !check_permission(c) {
		c.JSON(http.StatusOK, types.UpdateMemberResponse{Code: types.PermDenied})
		return
	}

	var data types.UpdateMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}

	var res types.UpdateMemberResponse
	state := check_state_byID(data.UserID)

	if state == types.UserHasExisted {
		res.Code = types.OK
		database.Db.Table("members").Where("user_id=?", data.UserID).Update("Nickname", data.Nickname)
	} else {
		res.Code = state
	}

	c.JSON(http.StatusOK, res)
}

func Member_delete(c *gin.Context) {
	if !check_permission(c) {
		c.JSON(http.StatusOK, types.DeleteMemberResponse{Code: types.PermDenied})
		return
	}
	var data types.DeleteMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}

	var res types.DeleteMemberResponse
	state := check_state_byID(data.UserID)

	if state == types.UserHasExisted {
		res.Code = types.OK
		database.Db.Table("members").Where("user_id=?", data.UserID).Delete(&types.Member{})
	} else {
		res.Code = state
	}
	c.JSON(http.StatusOK, res)
}
