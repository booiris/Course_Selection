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

	// TODO :测试用
	// user_type := database.FindUserType(data.StudentID)
	// if user_type == 0 {
	// 	c.JSON(http.StatusOK, types.BookCourseResponse{Code: types.UserNotExisted})
	// 	return
	// } else if user_type < 0 {
	// 	c.JSON(http.StatusOK, types.BookCourseResponse{Code: types.UserHasDeleted})
	// 	return
	// } else if user_type != types.Student {
	// 	c.JSON(http.StatusOK, types.BookCourseResponse{Code: types.StudentNotExisted})
	// 	return
	// }

	check_course := database.CheckCourse(data.CourseID)
	if !check_course {
		c.JSON(http.StatusOK, types.BookCourseResponse{Code: types.CourseNotExisted})
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
		database.Rdb.RPush(context, "Sync_mysql", data.StudentID+","+data.CourseID)
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

	// TODO :测试用
	// user_type := database.FindUserType(data.StudentID)
	// if user_type == 0 {
	// 	c.JSON(http.StatusOK, types.GetStudentCourseResponse{Code: types.UserNotExisted})
	// 	return
	// } else if user_type < 0 {
	// 	c.JSON(http.StatusOK, types.GetStudentCourseResponse{Code: types.UserHasDeleted})
	// 	return
	// } else if user_type != types.Student {
	// 	c.JSON(http.StatusOK, types.GetStudentCourseResponse{Code: types.StudentNotExisted})
	// 	return
	// }

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
