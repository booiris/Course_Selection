package controller

import (
	"context"
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"

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
	ctx := context.Background()
	database.Rdb.Set(ctx, u.CourseID+"cnt", u.Cap, 0)
	c.JSON(http.StatusOK, types.CreateCourseResponse{
		Code: types.OK,
		Data: struct{ CourseID string }{CourseID: u.CourseID},
	})
}

// 获取课程对应老师
func Course_get(c *gin.Context) {
	var getCourseRequest types.GetCourseRequest
	//表单数据绑定，如果填的参数不规范，或者为空应该返回参数不合法
	if err := c.ShouldBind(&getCourseRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	var course types.TCourse
	database.Db.Table("Courses").First(&course, getCourseRequest.CourseID)
	if course == (types.TCourse{}) { //如果数据库中没有查到课程说明课程不存在
		c.JSON(http.StatusOK, types.GetCourseResponse{
			Code: types.CourseNotExisted,
			Data: course,
		})
	} else {
		c.JSON(http.StatusOK, types.GetCourseResponse{
			Code: types.OK,
			Data: course,
		})
	}
}

// 老师绑定课程
func Teacher_bind_course(c *gin.Context) {
	var bindCourseRequest types.BindCourseRequest
	//表单数据绑定
	if err := c.ShouldBind(&bindCourseRequest); err != nil {
		c.JSON(http.StatusInternalServerError, types.ResponseMeta{Code: types.ParamInvalid})
		return
	}
	//查询课程是否被绑定
	var course types.Course
	database.Db.First(&course, bindCourseRequest.CourseID)
	if course == (types.Course{}) { //如果数据库中没有查到课程说明课程不存在
		c.JSON(http.StatusOK, types.BindCourseResponse{
			Code: types.CourseNotExisted,
		})
		return
	}
	//不用判断老师id，因为types里给出了老师id不用做落户校验
	//如果查询到了课程
	if course.TeacherID == "" { //如果teacherid字段为空，说明可以绑定
		database.Db.Table("courses").Where("course_id=?", bindCourseRequest.CourseID).Update("teacher_id", bindCourseRequest.TeacherID)
		c.JSON(http.StatusOK, types.BindCourseResponse{Code: types.OK})
	} else { //否则返回课程已绑定过
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
	//查询课程是否被绑定
	var course types.Course
	database.Db.First(&course, unbindCourseRequest.CourseID)
	if course == (types.Course{}) { //如果数据库中没有查到课程说明课程不存在
		c.JSON(http.StatusOK, types.UnbindCourseResponse{
			Code: types.CourseNotExisted,
		})
		return
	}
	//同理 老师id应该不用判断，因为没有老师不存在的状态码
	//如果课程存在
	if course.TeacherID == "" { //teacher字段为空，说明该课程没有绑定
		c.JSON(http.StatusOK, types.UnbindCourseResponse{Code: types.CourseNotBind})
	} else { // 否则解绑
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
	// 同理 这里不用判断老师id是否存在
	var courses []*types.TCourse
	database.Db.Table("courses").Where("teacher_id", getTeacherCourseRequest.TeacherID).Find(&courses)
	c.JSON(http.StatusOK, types.GetTeacherCourseResponse{
		Code: types.OK,
		Data: struct{ CourseList []*types.TCourse }{CourseList: courses},
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
