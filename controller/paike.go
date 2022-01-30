package controller

import (
	"course_selection/database"
	"course_selection/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建课程
func Course_create(c *gin.Context) {
	var data types.CreateCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	database.Db.Table("courses").Create(&data)
	var res struct{ CourseID string }
	database.Db.Table("courses").Where(map[string]interface{}{"name": data.Name}).Find(&res)
	c.JSON(http.StatusOK, types.CreateCourseResponse{Code: types.OK, Data: struct{ CourseID string }{res.CourseID}})
}

// 获取课程对应老师
func Course_get(c *gin.Context) {
	var data types.GetCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	var res types.TCourse
	database.Db.Table("courses").Where(&data).Find(&res)
	c.JSON(http.StatusOK, types.GetCourseResponse{Code: types.OK, Data: res})
}

// 老师绑定课程
func Teacher_bind_course(c *gin.Context) {
	var data types.BindCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	//TODO 已绑定判断
	database.Db.Table("courses").Where("course_id=?", data.CourseID).Update("teacher_id", data.TeacherID)
	c.JSON(http.StatusOK, types.BindCourseResponse{Code: types.OK})
}

// 解绑
func Teacher_unbind_course(c *gin.Context) {
	var data types.UnbindCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	//TODO 未绑定判断
	database.Db.Table("courses").Where(&data).Update("teacher_id", "")
	c.JSON(http.StatusOK, types.UnbindCourseResponse{Code: types.OK})
}

// 获取老师所有能上的课
func Teacher_get_course(c *gin.Context) {
	var data types.GetTeacherCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	var res []*types.TCourse
	database.Db.Table("courses").Where(&data).Find(&res)
	var temp struct {
		CourseList []*types.TCourse
	}
	temp.CourseList = res
	c.JSON(http.StatusOK, types.GetTeacherCourseResponse{Code: types.OK, Data: temp})
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
