package controller

import (
	"course_selection/database"
	"course_selection/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 创建成员
func Member_create(c *gin.Context) {
	var createCourseRequest types.CreateMemberRequest
	//表单数据绑定
	if err := c.ShouldBind(&createCourseRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	value, err := c.Cookie("camp-session")
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	var user types.Member
	database.Db.Where("username=?", value).First(&user)
	if user.UserType == 1 {
		database.Db.Select("Nickname", "Username", "Password", "UserType").Create(&types.Member{
			Nickname: createCourseRequest.Nickname,
			Username: createCourseRequest.Username,
			Password: createCourseRequest.Password,
			UserType: createCourseRequest.UserType})
		var u types.Member
		database.Db.Where("username=?", createCourseRequest.Username).First(&u)
		c.JSON(http.StatusOK, types.CreateMemberResponse{
			types.OK,
			struct {
				UserID string
			}{u.UserID},
		})
	} else {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.PermDenied})
	}
}

// 获取单个成员信息
func Member_get(c *gin.Context) {
	var getMemberRequest types.GetMemberRequest
	//表单数据绑定
	if err := c.ShouldBind(&getMemberRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	var user types.Member
	database.Db.First(&user, getMemberRequest.UserID)
	if user != (types.Member{}) {
		c.JSON(http.StatusOK, types.GetMemberResponse{
			types.OK,
			types.TMember{
				UserID:   user.UserID,
				Username: user.Username,
				UserType: user.UserType,
				Nickname: user.Nickname},
		})
	} else {
		database.Db.Unscoped().First(&user, getMemberRequest.UserID)
		if user == (types.Member{}) {
			c.JSON(http.StatusOK, types.GetMemberResponse{Code: types.UserNotExisted})
		} else {
			c.JSON(http.StatusOK, types.GetMemberResponse{Code: types.UserHasDeleted})
		}
	}
}

// 获取成员列表
func Member_get_list(c *gin.Context) {
	var getMemberListRequest types.GetMemberListRequest
	//表单数据绑定
	if err := c.ShouldBind(&getMemberListRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	var users []types.TMember
	database.Db.Table("Members").Limit(getMemberListRequest.Limit).Offset(getMemberListRequest.Offset).Find(&users)

	c.JSON(http.StatusOK, types.GetMemberListResponse{
		types.OK,
		struct{ MemberList []types.TMember }{MemberList: users},
	})
}

// 更新成员
func Member_update(c *gin.Context) {
	var updateMemberRequest types.UpdateMemberRequest
	//表单数据绑定
	if err := c.ShouldBind(&updateMemberRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	c.JSON(http.StatusOK, types.UpdateMemberResponse{Code: types.OK})
	database.Db.Model(&types.Member{}).Where("user_id = ?", updateMemberRequest.UserID).Update("nickname", updateMemberRequest.Nickname)
}

// 删除成员
func Member_delete(c *gin.Context) {
	var deleteMemberRequest types.DeleteMemberRequest
	//表单数据绑定
	if err := c.ShouldBind(&deleteMemberRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	c.JSON(http.StatusOK, types.DeleteMemberResponse{Code: types.OK})
	database.Db.Delete(&types.Member{}, deleteMemberRequest.UserID)
}
