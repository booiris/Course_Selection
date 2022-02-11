package router

import (
	"course_selection/controller"
	"course_selection/types"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {

	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create", controller.Member_create)
	g.GET("/member", controller.Member_get)
	g.GET("/member/list", controller.Member_get_list)
	g.POST("/member/update", controller.Member_update)
	g.POST("/member/delete", controller.Member_delete)

	// 登录

	g.POST("/auth/login", controller.Login)
	g.POST("/auth/logout", controller.Logout)
	g.GET("/auth/whoami", controller.Whoami)

	// 排课
	g.POST("/course/create", controller.Course_create)
	g.GET("/course/get", controller.Course_get)

	g.POST("/teacher/bind_course", controller.Teacher_bind_course)
	g.POST("/teacher/unbind_course", controller.Teacher_unbind_course)
	g.GET("/teacher/get_course", controller.Teacher_get_course)
	g.POST("/course/schedule", controller.Course_schedule)

	// 抢课
	g.POST("/student/book_course", controller.Student_book_course)
	g.GET("/student/course", controller.Student_course)

	types.Count = 10
	g.GET("/count", controller.Count)

}
