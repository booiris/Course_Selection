package controller

import (
	"context"
	"course_selection/database"
	"course_selection/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 选课
func Student_book_course(c *gin.Context) {
	var ctx = context.Background()
	var bookCourseRequest types.BookCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&bookCourseRequest); err != nil {
		c.JSON(http.StatusOK, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}

	//课程已满，直接返回
	if types.CurrentSoldOutMap[bookCourseRequest.CourseID] {
		c.JSON(http.StatusOK, types.ResponseMeta{Code: types.CourseNotAvailable})
		return
	}

	//获取缓存中数据
	nownum := database.Rdb.Decr(ctx, bookCourseRequest.CourseID+"cnt")
	if nownum.Val() < 0 {
		//标记课程已满
		types.CurrentSoldOutMap[bookCourseRequest.CourseID] = true
		//恢复，防止少卖
		database.Rdb.Incr(ctx, bookCourseRequest.CourseID+"cnt")
		c.JSON(http.StatusOK, types.ResponseMeta{Code: types.CourseNotAvailable})
		return
	}

	//抢课成功，创建记录
	scourse := types.SCourse{
		bookCourseRequest.CourseID,
		bookCourseRequest.StudentID,
	}

	if err := database.Db.Create(&scourse); err.Error != nil {
		//课程创建失败，恢复缓存
		nownum := database.Rdb.Incr(ctx, bookCourseRequest.CourseID+"cnt")
		//恢复课程到未满状态
		if nownum.Val() > 0 && types.CurrentSoldOutMap[bookCourseRequest.CourseID] {
			types.CurrentSoldOutMap[bookCourseRequest.CourseID] = false
		}
		c.JSON(http.StatusOK, types.ResponseMeta{Code: types.CourseNotAvailable})
		return
	}
	c.JSON(http.StatusOK, types.ResponseMeta{Code: types.OK})
}

// 获取学生选课列表
func Student_course(c *gin.Context) {
	var getStudentCourseRequest types.GetStudentCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&getStudentCourseRequest); err != nil {
		c.JSON(http.StatusOK, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}

}
