package controller

import (
	"context"
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func check_permission(c *gin.Context) bool {
	sessionKey, err := c.Cookie("camp-session")
	if err != nil {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.LoginRequired})
		return false
	}

	ctx := context.Background()
	userid := database.Rdb.Get(ctx, sessionKey)
	if userid.Err() == redis.Nil {
		c.JSON(200, types.CreateMemberResponse{Code: types.LoginRequired})
		return false
	}

	var res types.Member
	database.Db.Model(types.Member{}).Unscoped().Where("user_id=?", userid.Val()).Find(&res)
	if res == (types.Member{}) {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.UserNotExisted})
		return false
	} else if res.Deleted.Valid {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.UserHasDeleted})
		return false
	}
	if res.UserType != types.Admin {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.PermDenied})
		return false
	}
	return true
}

func check_param(data types.CreateMemberRequest) bool {

	// 用户昵称
	for _, v := range data.Nickname {
		if !(unicode.IsLetter(v)) {
			return false
		}
	}
	len := strings.Count(data.Nickname, "")
	if len < 4 || len > 20 {
		return false
	}

	// 用户名
	for _, v := range data.Username {
		if !(unicode.IsLetter(v)) {
			return false
		}
	}
	len = strings.Count(data.Username, "")
	if len < 8 || len > 20 {
		return false
	}

	// 密码
	upper, lower, digit := false, false, false
	for _, v := range data.Password {
		if !(unicode.IsLetter(v) || unicode.IsDigit(v)) {
			return false
		}
		if unicode.IsUpper(v) {
			upper = true
		} else if unicode.IsDigit(v) {
			digit = true
		} else if unicode.IsLower(v) {
			lower = true
		}
	}
	if !(digit && lower && upper) {
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

func Member_create(c *gin.Context) {
	if !check_permission(c) {
		return
	}
	var data types.CreateMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	if check_param(data) {
		var res types.Member
		database.Db.Model(types.Member{}).Unscoped().Where("username=?", data.Username).Find(&res)
		if res == (types.Member{}) {
			database.Db.Table("members").Create(&data)
			database.Db.Model(types.Member{}).Where("username=?", data.Username).Find(&res)
			ctx := context.Background()
			database.Rdb.Set(ctx, "usertype"+res.UserID, int(res.UserType), 0)
			c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.OK, Data: struct{ UserID string }{res.UserID}})
		} else if res.Deleted.Valid {
			c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.UserHasDeleted})
		} else {
			c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.UserHasExisted})
		}
	} else {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.ParamInvalid})
	}
}

// 获取单个成员信息
func Member_get(c *gin.Context) {
	var data types.GetMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	var res types.Member
	database.Db.Model(types.Member{}).Unscoped().Where(&data).Find(&res)
	if res == (types.Member{}) {
		c.JSON(http.StatusOK, types.GetMemberResponse{Code: types.UserNotExisted})
	} else {
		if res.Deleted.Valid {
			c.JSON(http.StatusOK, types.GetMemberResponse{Code: types.UserHasDeleted})
		} else {
			send_data := types.TMember{
				UserID:   res.UserID,
				Username: res.Username,
				UserType: res.UserType,
				Nickname: res.Nickname,
			}
			c.JSON(http.StatusOK, types.GetMemberResponse{Code: types.OK, Data: send_data})
		}
	}
}

// 获取成员列表
func Member_get_list(c *gin.Context) {
	var data types.GetMemberListRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	var users []types.TMember
	if err := database.Db.Model(types.Member{}).Limit(data.Limit).Offset(data.Offset).Find(&users).Error; err != nil {
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
	var check types.Member
	database.Db.Model(types.Member{}).Unscoped().Where("user_id=?", data.UserID).Find(&check)
	if check == (types.Member{}) {
		c.JSON(http.StatusOK, types.UpdateMemberResponse{Code: types.UserNotExisted})
		return
	} else if check.Deleted.Valid {
		c.JSON(http.StatusOK, types.UpdateMemberResponse{Code: types.UserHasDeleted})
		return
	}

	database.Db.Model(types.Member{}).Where("user_id=?", data.UserID).Update("Nickname", data.Nickname)
	c.JSON(http.StatusOK, types.UpdateMemberResponse{Code: types.OK})
}

// 删除成员
func Member_delete(c *gin.Context) {
	var data types.DeleteMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}

	var user types.Member
	database.Db.Model(types.Member{}).Unscoped().Where("user_id=?", data.UserID).Find(&user)
	if user == (types.Member{}) {
		c.JSON(http.StatusOK, types.DeleteMemberResponse{Code: types.UserNotExisted})
		return
	} else if user.Deleted.Valid {
		c.JSON(http.StatusOK, types.DeleteMemberResponse{Code: types.UserHasDeleted})
		return
	}
	if user.UserType == types.Teacher {
		database.Db.Table("courses").Where("teacher_id=?", user.UserID).Update("teacher_id", nil)
	} else if user.UserType == types.Student {
		ctx := context.Background()
		database.Rdb.Del(ctx, user.UserID)
		database.Db.Where("user_id=?", user.UserID).Delete(types.SCourse{})
	}
	ctx := context.Background()
	database.Rdb.Set(ctx, "usertype"+data.UserID, -1, 0)
	database.Db.Where("user_id=?", data.UserID).Delete(&types.Member{})
	c.JSON(http.StatusOK, types.DeleteMemberResponse{Code: types.OK})
}
