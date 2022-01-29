package controller

import (
	"course_selection/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建课程
func Course_create(c *gin.Context) {
}

// 获取课程对应老师
func Course_get(c *gin.Context) {
}

// 老师绑定课程
func Teacher_bind_course(c *gin.Context) {
}

// 解绑
func Teacher_unbind_course(c *gin.Context) {
}

// 获取老师所有能上的课
func Teacher_get_course(c *gin.Context) {
}

// 二分图匹配对应课程
func Course_schedule(c *gin.Context) {
	var data types.ScheduleCourseRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		return
	}
	res := match_course(data.TeacherCourseRelationShip)
	c.JSON(http.StatusOK, types.ScheduleCourseResponse{Code: types.OK, Data: res})
}

func match_course(data map[string][]string) map[string]string {
	t2c := make(map[string]string)
	c2t := make(map[string]string)
	for k := range data {
		if _, is_exist := t2c[k]; !is_exist {
			check := make(map[string]struct{})
			dfs(k, &data, &check, &t2c, &c2t)
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
