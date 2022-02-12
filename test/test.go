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
	"time"
)

const (
	thread     int    = 3000 //访问线程个数
	q_cnt      int    = 10   //每个线程发出请求次数
	host       string = "http://180.184.74.221:"
	group      string = "/api/v1"
	port       string = "2000"
	course_cnt int    = 300
	// dir    string = "/course/create"
)

var student [thread + 1][]string
var wg sync.WaitGroup
var globalTransport *http.Transport
var client http.Client

func client_init() {
	globalTransport = &http.Transport{
		DisableKeepAlives: true,
	}
	client = http.Client{
		// Timeout:   30,
		Transport: globalTransport,
	}
}

func create_course() {
	dir := "/course/create"
	for i := 0; i < course_cnt; i++ {
		name := "C" + strconv.Itoa(i)
		cap := strconv.Itoa(rand.Intn(300) + 1)
		params := url.Values{"Name": {name}, "Cap": {cap}}
		resp, err := client.PostForm(host+port+group+dir, params)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		res, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(res))

	}
}

func book_course(student_id int) {
	defer wg.Done()
	// 检验过程中检查数据一致性
	check := make(chan bool)
	over := make(chan struct{})
	defer close(check)
	defer close(over)
	go get_student_courses_chan(student_id, check, over)

	dir := "/student/book_course"
	id := strconv.Itoa(student_id)
	for i := 0; i < q_cnt; i++ {
		course_id := strconv.Itoa(rand.Intn(course_cnt) + 1)
		params := url.Values{"StudentID": {id}, "CourseID": {course_id}}
		resp, err := client.PostForm(host+port+group+dir, params)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		res, _ := ioutil.ReadAll(resp.Body)
		var data types.BookCourseResponse
		json.Unmarshal(res, &data)

		if data.Code == types.OK {
			student[student_id] = append(student[student_id], course_id)
		}

		key := rand.Intn(thread)
		if key == 0 {
			check <- true
			<-over
		}
	}
	check <- false
}

func get_student_courses_chan(student_id int, check chan bool, over chan struct{}) {

	id := strconv.Itoa(student_id)
	dir := "/student/course"
	params := "?StudentID=" + id
	for {
		flag := <-check
		if !flag {
			break
		}
		resp, err := client.Get(host + port + group + dir + params)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var data types.GetStudentCourseResponse
		json.Unmarshal(body, &data)

		temp := make([]int, len(student[student_id]))
		for i := range temp {
			temp[i], _ = strconv.Atoi(student[student_id][i])
		}
		sort.Ints(temp)
		for i := range temp {
			if strconv.Itoa(temp[i]) != data.Data.CourseList[i].CourseID {
				fmt.Println(id+":wrong\n", string(body)+"\n", temp)
				break
			}
		}

		over <- struct{}{}
	}
}

func get_student_courses(student_id int) {
	id := strconv.Itoa(student_id)
	dir := "/student/course"
	params := "?StudentID=" + id
	resp, err := client.Get(host + port + group + dir + params)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var data types.GetStudentCourseResponse
	json.Unmarshal(body, &data)
	if len(data.Data.CourseList) != len(student[student_id]) {
		fmt.Println(id + ":wrong")
		return
	}

	temp := make([]int, len(student[student_id]))
	for i := range temp {
		temp[i], _ = strconv.Atoi(student[student_id][i])
	}
	sort.Ints(temp)
	for i := range temp {
		if strconv.Itoa(temp[i]) != data.Data.CourseList[i].CourseID {
			fmt.Println(id + ":wrong")
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
	wg.Add(thread)
	for i := 1; i <= thread; i++ {
		go book_course(i)
	}
	wg.Wait()
	fmt.Println("检查最终数据一致性")
	time.Sleep(time.Duration(500) * time.Millisecond)
	for i := 1; i <= thread; i++ {
		get_student_courses(i)
	}
}
