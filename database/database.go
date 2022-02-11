package database

import (
	"context"
	"course_selection/types"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//数据库配置
const (
	userName = "root"
	password = "12345678"
	// password = "bytedancecamp"
	ip     = "127.0.0.1"
	port   = "3306"
	dbName = "isuse"
)

var Db *gorm.DB
var Rdb *redis.Client

func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8mb4"
	dsn := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8mb4"}, "")
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		fmt.Println("open database fail")
		return
	}
	sqlDb, _ := Db.DB()
	// 设置空闲连接数
	// 数量 connections = ((core_count * 2) + effective_spindle_count)
	sqlDb.SetConnMaxIdleTime(4)
	// 最大连接数
	sqlDb.SetMaxOpenConns(100)
	// 连接复用连接时间
	sqlDb.SetConnMaxLifetime(time.Hour)
	fmt.Println("connnect success")
}

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",  // no password set
		DB:       1,   // use default DB
		PoolSize: 100, // 连接池大小
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := Rdb.Ping(ctx).Result()

	if err != nil {
		fmt.Println("open redis fail")
		return
	}

	// 删除 redis 缓存
	res, err := Rdb.FlushDB(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("delete redis:", res)

	// 读取学生选课和课程列表，构建缓存
	course_cnt := make(map[string]int)
	var courses []struct {
		CourseID string
		Cap      int
	}
	Db.Table("courses").Find(&courses)
	for i := range courses {
		course_cnt[courses[i].CourseID] = courses[i].Cap
	}

	var data []types.SCourse
	Db.Table("s_courses").Find(&data)
	for i := range data {
		course_cnt[data[i].CourseID] -= 1
		err := Rdb.HSetNX(ctx, data[i].UserID, data[i].CourseID, 0).Err()
		if err != nil {
			panic(err)
		}
	}

	for k, v := range course_cnt {
		Rdb.Set(ctx, k+"cnt", v, 0)
	}

	fmt.Println("open redis success")
}
