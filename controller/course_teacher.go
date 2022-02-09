package controller

import (
	"course_selection/types"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"course_selection/database"
	"github.com/gin-gonic/gin"
)

// 创建课程
func Course_create(c *gin.Context) {
	var createCourseRequest types.CreateCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&createCourseRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	database.Db.Table("courses").Create(&createCourseRequest)
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
		database.Db.Table("courses").Where("course_id=?", bindCourseRequest.CourseID).Update("teacher_id", bindCourseRequest.TeacherID)
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
		database.Db.Table("courses").Where("course_id=?", unbindCourseRequest.CourseID).Update("teacher_id", "")
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
	var data types.ScheduleCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	res := match_course(&data.TeacherCourseRelationShip)
	c.JSON(http.StatusOK, types.ScheduleCourseResponse{Code: types.OK, Data: res})
}

func match_course(data *map[string][]string) map[string]string {
	t2c := make(map[string]string)
	c2t := make(map[string]string)
	for k := range *data {
		if _, is_exist := t2c[k]; !is_exist {
			check := make(map[string]struct{})
			dfs(k, data, &check, &t2c, &c2t)
		}
	}
	return t2c
}

func dfs(now string, data *map[string][]string, check *map[string]struct{}, t2c *map[string]string, c2t *map[string]string) bool {
	for _, i := range (*data)[now] {
		if _, is_exist := (*check)[i]; !is_exist {
			(*check)[i] = struct{}{}
			if nxt, is_exist := (*c2t)[i]; !is_exist || dfs(nxt, data, check, t2c, c2t) {
				(*c2t)[i] = now
				(*t2c)[now] = i
				return true
			}
		}
	}
	return false
}
