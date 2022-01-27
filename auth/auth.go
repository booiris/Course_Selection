package auth

import (
	"course_selection/database"
	"course_selection/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var data types.LoginRequest
	if err := c.ShouldBind(&data); err != nil {
		fmt.Println(err) //TODO 删除错误提示
		c.JSON(http.StatusBadRequest, types.LoginResponse{Code: types.WrongPassword})
		return
	}
	var res types.LoginRequest
	database.Db.Table("TMember").Where(&data).Find(&res)
	if res == (types.LoginRequest{}) {
		c.JSON(http.StatusOK, types.LoginResponse{Code: types.WrongPassword})
	} else {
		c.JSON(http.StatusOK, types.LoginResponse{Code: types.OK})
	}
	//TODO 添加cookie
}

func Logout(c *gin.Context) {

}
func Whoami(c *gin.Context) {

}
