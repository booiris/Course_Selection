package auth

import (
	"encoding/json"
	"fmt"
	"helloWorld/dataBase"
	"helloWorld/rule"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func Login(c *gin.Context) {
	/**
	 * 读取 JSON 获取账号密码
	 */
	body := c.Request.Body
	read, _ := ioutil.ReadAll(body)
	var request rule.CreateMemberRequest
	errC := json.Unmarshal(read, &request)
	if errC != nil {
		println("json read over")
		return
	}

	/**
	 * 获取数据库连接
	 */
	db, err := dataBase.CreateConnection()
	if err != nil {
		return
	}

	/*
	 * 提取用户 id 并发至cookie
	 */
	sqlStr := "SELECT UserID From member WHERE Username=? AND UserPassword=?"
	var id string
	errQ := db.QueryRow(sqlStr, request.Username, request.Password).Scan(&id)
	if errQ != nil {
		c.String(http.StatusOK, fmt.Sprintf("cann't login"))
	} else {
		c.String(http.StatusOK, fmt.Sprintf("Login Success. The id is %s", id))
	}

	c.SetCookie("user_cookie", id, 1000, "/", "localhost", false, true)
	/*
	 * 断开连接
	 */

	dataBase.CloseConnection(db)
}
