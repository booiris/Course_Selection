package control

import (
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func check_permission(c *gin.Context) bool {
	value, err := c.Cookie("camp-session")
	if err != nil {
		log.Println(err)
		return false
	}
	var res types.TMember
	database.Db.Table("TMember").Where(&value).Find(&res.UserType)
	return res.UserType == types.Admin
}

func check_param(c *gin.Context) bool {
	return true
}

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
	check_param(c)
}

func Member_get(c *gin.Context) {
	if !check_permission(c) {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.PermDenied})
		return
	}
	var data types.GetMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	//TODO 完成
}

func Member_get_list(c *gin.Context) {
	if !check_permission(c) {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.PermDenied})
		return
	}
	//TODO 完成
}

func Member_update(c *gin.Context) {
	if !check_permission(c) {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.PermDenied})
		return
	}
	var data types.UpdateMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	database.Db.Table("TMember").Where(&data.UserID).Update("Nickname", data.Nickname)
	c.JSON(http.StatusOK, types.UpdateMemberResponse{Code: types.OK})
}

func Member_delete(c *gin.Context) {
	if !check_permission(c) {
		c.JSON(http.StatusOK, types.CreateMemberResponse{Code: types.PermDenied})
		return
	}
	var data types.DeleteMemberRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	// TODO 软删除
	database.Db.Table("TMember").Where(&data)
	c.JSON(http.StatusOK, types.DeleteMemberResponse{Code: types.OK})
}
