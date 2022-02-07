package controller

import (
	"course_selection/database"
	"course_selection/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 创建课程
func Course_create(c *gin.Context) {
	var createCourseRequest types.CreateCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&createCourseRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	database.Db.Select("Name", "Cap").Create(&types.Course{
		Name: createCourseRequest.Name,
		Cap:  createCourseRequest.Cap})
	var u types.Course
	database.Db.Where("Name=?", createCourseRequest.Name).First(&u)

	c.JSON(http.StatusOK, types.CreateCourseResponse{
		types.OK,
		struct{ CourseID string }{CourseID: u.CourseID},
	})
}

// 获取课程对应老师
func Course_get(c *gin.Context) {
	var getCourseRequest types.GetCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&getCourseRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	var course types.TCourse
	database.Db.Table("Courses").First(&course, getCourseRequest.CourseID)

	c.JSON(http.StatusOK, types.GetCourseResponse{
		Code: types.OK,
		Data: course,
	})
}

// 老师绑定课程
func Teacher_bind_course(c *gin.Context) {
	var bindCourseRequest types.BindCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&bindCourseRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	var course types.Course
	//查询课程是否被绑定
	database.Db.First(&course, bindCourseRequest.CourseID)
	if course.TeacherID == "" {
		course.TeacherID = bindCourseRequest.TeacherID
		c.JSON(http.StatusOK, types.BindCourseResponse{Code: types.OK})
	} else {
		c.JSON(http.StatusOK, types.BindCourseResponse{Code: types.CourseHasBound})
	}
}

// 解绑
func Teacher_unbind_course(c *gin.Context) {
	var unbindCourseRequest types.UnbindCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&unbindCourseRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	var course types.Course
	//查询课程是否被绑定
	database.Db.First(&course, unbindCourseRequest.CourseID)
	if course.TeacherID == "" {
		c.JSON(http.StatusOK, types.UnbindCourseResponse{Code: types.CourseNotBind})
	} else {
		course.TeacherID = ""
		c.JSON(http.StatusOK, types.UnbindCourseResponse{Code: types.OK})
	}
}

// 获取老师所有能上的课
func Teacher_get_course(c *gin.Context) {
	var getTeacherCourseRequest types.GetTeacherCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&getTeacherCourseRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	var courses []*types.TCourse
	database.Db.Table("Courses").Where("TeacherID", getTeacherCourseRequest.TeacherID).Find(&courses)
	c.JSON(http.StatusOK, types.GetTeacherCourseResponse{
		types.OK,
		struct{ CourseList []*types.TCourse }{CourseList: courses},
	})
}

// 二分图匹配对应课程
func Course_schedule(c *gin.Context) {
	var scheduleCourseRequest types.ScheduleCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&scheduleCourseRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}

	ans := slove(scheduleCourseRequest.TeacherCourseRelationShip)

	c.JSON(http.StatusOK, types.ScheduleCourseResponse{
		types.OK,
		ans,
	})
}

func slove(TeacherCourseRelationShip map[string][]string) map[string]string {
	var ans map[string]string
	return ans
}
