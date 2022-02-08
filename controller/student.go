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

	var data types.BookCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	nownum := database.Rdb.Decr(context, data.CourseID+"cnt")
	if nownum.Val() < 0 {
		database.Rdb.Incr(context, data.CourseID+"cnt")
		c.JSON(http.StatusOK, types.BookCourseResponse{Code: types.CourseNotAvailable})
	} else {
		err := database.Rdb.HSetNX(context, data.StudentID, data.CourseID, 0).Err()
		if err != nil {
			panic(err)
		}
		create_data := types.SCourse{
			UserID:   data.StudentID,
			CourseID: data.CourseID,
		}
		database.Db.Create(&create_data)
		c.JSON(http.StatusOK, types.BookCourseResponse{Code: types.OK})
	}
}

// 获取学生选课列表
func Student_course(c *gin.Context) {
	var data types.GetStudentCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	var res []types.TCourse
	database.Db.Table("courses").Where("course_id in (?)", database.Db.Table("s_courses").Select("course_id").Where("user_id=?", data.StudentID)).Find(&res)

	var temp struct {
		CourseList []types.TCourse
	}
	temp.CourseList = res
	c.JSON(http.StatusOK, types.GetStudentCourseResponse{Code: types.OK, Data: temp})
}
