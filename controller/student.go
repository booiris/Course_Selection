package controller

import (
	"course_selection/database"
	"course_selection/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 选课
func Student_book_course(c *gin.Context) {
	var bookCourseRequest types.BookCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&bookCourseRequest); err != nil {
		c.JSON(http.StatusOK, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	var course types.Course
	//查询课程
	database.Db.First(&course, bookCourseRequest.CourseID)
	if course == (types.Course{}) {
		c.JSON(http.StatusOK, types.ResponseMeta{Code: types.CourseNotExisted})
		return
	}
	if course.Cap > 0 {
		var student types.SCourse
		database.Db.First(&student, bookCourseRequest.StudentID)
		if student == (types.SCourse{}) {
			c.JSON(http.StatusOK, types.ResponseMeta{Code: types.StudentNotExisted})
		} else {
			course.Cap -= 1
			student.CourseID = bookCourseRequest.CourseID
			//修改后的值还需更新到数据库
			c.JSON(http.StatusOK, types.ResponseMeta{Code: types.OK})
		}
	} else {
		c.JSON(http.StatusOK, types.ResponseMeta{Code: types.CourseNotAvailable})
	}
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
