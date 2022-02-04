package controller

import (
	"context"
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 选课
func Student_book_course(c *gin.Context) {

	context := context.Background()

	database.Rdb.Incr(context, "123")
}

// 获取学生选课列表
func Student_course(c *gin.Context) {
	var data types.GetStudentCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	var courseids []struct{ CourseID string }
	database.Db.Table("s_courses").Where(&data).Find(&courseids)
	var res []types.TCourse
	database.Db.Table("courses").Where(&courseids).Find(&res)
	var temp struct {
		CourseList []types.TCourse
	}
	temp.CourseList = res
	c.JSON(http.StatusOK, types.GetStudentCourseResponse{Code: types.OK, Data: temp})
}
