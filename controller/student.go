package controller

import (
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 选课
func Student_book_course(c *gin.Context) {

	var data types.BookCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}

	var usertype types.UserType
	database.Db.Table("members").Where("user_id = ?", data.StudentID).Find("user_type = ?", usertype)

	if usertype != types.Student {
		// StudentNotExisted  ErrNo = 11 // 学生不存在
		c.JSON(http.StatusOK, types.BookCourseResponse{Code: types.StudentNotExisted})
	}

	key, get := types.Course_Cap[data.CourseID]

	if !get {
		c.JSON(http.StatusOK, types.BookCourseResponse{Code: types.CourseNotExisted})
	} else if key == 0 {
		// CourseNotAvailable ErrNo = 7  // 课程已满
		c.JSON(http.StatusOK, types.BookCourseResponse{Code: types.CourseNotAvailable})
	} else {
		copyContext := c.Copy()
		// 异步处理
		go func() {
			var input types.BookCourseRequest
			copyContext.ShouldBind(&input)
			database.Db.Table("s_courses").Create(input)
		}()
	}
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
