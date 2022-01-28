package router

import (
	"course_selection/control"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create", control.Member_create)
	g.GET("/member", control.Member_get)
	g.GET("/member/list", control.Member_get_list)
	g.POST("/member/update", control.Member_update)
	g.POST("/member/delete", control.Member_delete)

	// 登录

	g.POST("/auth/login", control.Login)
	g.POST("/auth/logout", control.Logout)
	g.GET("/auth/whoami", control.Whoami)

	// 排课
	g.POST("/course/create")
	g.GET("/course/get")

	g.POST("/teacher/bind_course")
	g.POST("/teacher/unbind_course")
	g.GET("/teacher/get_course")
	g.POST("/course/schedule")

	// 抢课
	g.POST("/student/book_course")
	g.GET("/student/course")

}
