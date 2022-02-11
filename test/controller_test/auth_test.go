package controller_test

import (
	"course_selection/globals"
	"course_selection/test"
	"course_selection/types"
	"encoding/json"
	"fmt"
	"testing"
)

func TestLogin(t *testing.T) {

	// 初始化请求地址和请求参数
	uri := "http://127.0.1:1319/api/v1/auth/login"

	param := make(map[string]interface{})
	param["Username"] = "user1"
	param["Password"] = "JudgePassword2022"

	// 发起post请求，以表单形式传递参数
	body := test.PostJson(uri, param, globals.G)
	fmt.Printf("response:%v\n", string(body))

	// 解析响应，判断响应是否与预期一致
	response := &types.LoginResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		t.Errorf("解析响应出错，err:%v\n", err)
	}
	if response.Data.UserID != "4" {
		t.Errorf("响应数据不符，userID:%v\n", response.Data.UserID)
	}
}
