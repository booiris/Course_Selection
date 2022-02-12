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
		res := database.Rdb.HIncrBy(context, data.StudentID, data.CourseID, 1)
		if res.Err() != nil {
			panic(res.Err())
		}
		if res.Val() > 1 {
			database.Rdb.HIncrBy(context, data.StudentID, data.CourseID, -1)
			database.Rdb.Incr(context, data.CourseID+"cnt")
			c.JSON(http.StatusOK, types.BookCourseResponse{Code: types.StudentHasCourse})
			return
		}
		create_data := types.SCourse{
			UserID:   data.StudentID,
			CourseID: data.CourseID,
		}
		create_err := database.Db.Table("s_courses").Select("user_id", "course_id").Create(&create_data)
		if create_err.Error != nil {
			panic(create_err.Error)
		}
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
	// var res []types.TCourse
	// database.Db.Table("courses").Where("course_id in (?)", database.Db.Table("s_courses").Select("course_id").Where("user_id=?", data.StudentID)).Find(&res)

	ctx := context.Background()
	result := database.Rdb.HGetAll(ctx, data.StudentID)
	ids := make([]string, len(result.Val()))
	index := 0
	for k := range result.Val() {
		ids[index] = k
		index++
	}
	var res []types.TCourse
	database.Db.Table("courses").Where("course_id in (?)", ids).Find(&res)

	var temp struct {
		CourseList []types.TCourse
	}
	temp.CourseList = res
	c.JSON(http.StatusOK, types.GetStudentCourseResponse{Code: types.OK, Data: temp})
}
