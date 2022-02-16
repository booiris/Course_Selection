package main

import (
	"course_selection/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
)

// 参数
const (
	thread      int    = 2000 //访问线程个数
	student_num int    = 1000
	q_cnt       int    = 5 //每个线程发出请求次数
	host        string = "http://180.184.74.221:"
	group       string = "/api/v1"
	port        string = "80"
	course_cnt  int    = 300
)

// 整个压测过程是这样的

// 需要一个数组记录得到的结果，当抢课得到 OK 时，将数据写入数组，通过这个数组检查一致性
// 1. 首先调用 create_course 创建课程 （记得服务器运行前清空课程表和选课表）
// 2. 然后调用 book_course 开始抢课,抢课过程中检验数据一致性，检验方法如下：
//  	a. 抢课
//		b. 随机生成一个数，如果为0开始检查数据一致性，抢课阻塞，开始进行课程查询
// 		c. 查询课程，检查一致性，检查完毕后回到抢课协程
// 3. 等待 0.5s 等待mysql写入数据
// 4. 最后调用 get_student_courses 查询数据最终一致性

// 记录抢课的结果
var student [thread + 1][]string

// 协程锁，等待所有抢课协程结束再进行接下来步骤
var wg sync.WaitGroup

// 初始化http连接的 不用管
var globalTransport *http.Transport

// 并发量高的时候有时候连接会出错，记录这个错误，得查一下为什么
var fail_cnt int = 0
var requset_cnt int = 0

// 初始化http连接的 不用管
func client_init() {
	globalTransport = &http.Transport{
		DisableKeepAlives: true,
	}
}

// 创建课程
func create_course() {
	// 初始化 http 连接
	client := http.Client{
		// Timeout:   30,
		Transport: globalTransport,
	}
	dir := "/course/create"
	for i := 0; i < course_cnt; i++ {
		// 构造 课程名和容量
		name := "C" + strconv.Itoa(i)
		cap := strconv.Itoa(rand.Intn(260) + 40)
		params := url.Values{"Name": {name}, "Cap": {cap}}

		// 发送 post 请求
		resp, err := client.PostForm(host+port+group+dir, params)
		requset_cnt++
		if err != nil {
			panic(err)
		}

		// 关闭 http 连接
		defer resp.Body.Close()

		// 读取返回结果
		res, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(res))
	}
}

// 抢课函数 每一个协程对应一个学生
func book_course(student_id int) {
	// 初始化 http 连接
	client := http.Client{
		// Timeout:   30,
		Transport: globalTransport,
	}
	// 协程结束释放锁
	defer wg.Done()
	defer fmt.Println(student_id, " over")

	// 检验过程中检查数据一致性
	// 加入两个通道进行协程通信
	check := make(chan bool)
	over := make(chan struct{})
	defer close(check)
	defer close(over)
	// 新建一个查询协程，查询该学生的课程
	go get_student_courses_chan(student_id, check, over)

	dir := "/student/book_course"
	// 将学生id转化为字符串
	id := strconv.Itoa(student_id % student_num)
	for i := 0; i < q_cnt; i++ {

		// 构造课程 id
		course_id := strconv.Itoa(rand.Intn(course_cnt) + 1)
		params := url.Values{"StudentID": {id}, "CourseID": {course_id}}

		// 发送 post 请求
		resp, err := client.PostForm(host+port+group+dir, params)
		requset_cnt++
		if err != nil {
			fail_cnt++
			continue
		}

		// 关闭 http 连接
		defer resp.Body.Close()

		// 读取返回结果
		res, _ := ioutil.ReadAll(resp.Body)
		var data types.BookCourseResponse
		// 反序列化将字符串转换为对应的数据结构
		json.Unmarshal(res, &data)

		// 如果返回的是 OK ，将抢课结果记录
		if data.Code == types.OK {
			student[student_id] = append(student[student_id], course_id)
		}

		// 随机生成一个数，如果为 0 进行一致性查询
		key := rand.Intn(100)
		if key == 0 {
			check <- true
			<-over
		}
	}
	// 抢课结束时，写入false，关闭查询子协程
	check <- false
	<-over
}

func get_student_courses_chan(student_id int, check chan bool, over chan struct{}) {
	// 初始化 http 连接
	client := http.Client{
		// Timeout:   30,
		Transport: globalTransport,
	}
	id := strconv.Itoa(student_id % student_num)
	dir := "/student/course"
	params := "?StudentID=" + id
	for {
		flag := <-check
		// 当传入的为 false 时，抢课结束，结束循环，关闭协程
		if !flag {
			break
		}

		// 发送 get 请求
		resp, err := client.Get(host + port + group + dir + params)
		requset_cnt++
		if err != nil {
			fail_cnt++
			// 查询结束，抢课继续
			over <- struct{}{}
			continue
		}

		// 关闭 http 连接
		defer resp.Body.Close()

		// 读取返回结果
		body, _ := ioutil.ReadAll(resp.Body)
		var data types.GetStudentCourseResponse
		json.Unmarshal(body, &data)

		// 是这样检验数据一致性的
		// 如果本地存在记录 学生 1 已选 3，而 get 得到的数据中没有课程 3，证明一致性出错

		temp := make(map[string]struct{})
		for i := range data.Data.CourseList {
			temp[data.Data.CourseList[i].CourseID] = struct{}{}
		}

		for i := range student[student_id] {
			if _, is_exist := temp[student[student_id][i]]; !is_exist {
				fmt.Println(id+":wrong\n", string(body)+"\n", temp)

				// 查询结束，抢课继续
				over <- struct{}{}
				continue
			}
		}

		// 查询结束，抢课继续
		over <- struct{}{}
	}
	over <- struct{}{}
}

// 步骤和上面差不多，只不过少了通道通信的过程
func get_student_courses(student_id int) {
	// 初始化 http 连接
	client := http.Client{
		// Timeout:   30,
		Transport: globalTransport,
	}
	id := strconv.Itoa(student_id % student_num)
	dir := "/student/course"
	params := "?StudentID=" + id
	resp, err := client.Get(host + port + group + dir + params)
	requset_cnt++
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var data types.GetStudentCourseResponse
	json.Unmarshal(body, &data)
	if len(data.Data.CourseList) != len(student[student_id]) {
		fmt.Println(id + ":final len wrong")
		fmt.Println(string(body))
		fmt.Println(student[student_id])
		return
	}

	temp := make([]int, len(student[student_id]))
	for i := range temp {
		temp[i], _ = strconv.Atoi(student[student_id][i])
	}
	sort.Ints(temp)
	for i := range temp {
		if strconv.Itoa(temp[i]) != data.Data.CourseList[i].CourseID {
			fmt.Println(id + ":final wrong")
			fmt.Println(string(body))
			fmt.Println(temp)
			return
		}
	}
	//fmt.Println(string(body)+"\n", temp)
}

func main() {
	client_init()
	create_course()
	// wg.Add(thread)
	// for i := 1; i <= thread; i++ {
	// 	go book_course(i)
	// }
	// wg.Wait()
	// fmt.Println("检查最终数据一致性")
	// fmt.Println("fail:", fail_cnt)
	// fmt.Println("total:", requset_cnt)
	// time.Sleep(time.Duration(500) * time.Millisecond)
	// for i := 1; i <= thread; i++ {
	// 	get_student_courses(i)
	// }
}
