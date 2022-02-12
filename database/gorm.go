package database

import (
	"context"
	"course_selection/types"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

//数据库配置
const (
	userName = "root"
	password = "bytedancecamp"
	//bytedancecamp
	ip     = "180.184.74.221"
	port   = "3306"
	dbName = "jiar_user"
)

var Db *gorm.DB
var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "180.184.74.221:6379",
		Password: "zjredis", // no password set
		DB:       0,         // use default DB
		PoolSize: 100,       // 连接池大小
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

	// 读取课程列表，构建缓存
	course_cnt := make(map[string]int)
	var courses []types.Course
	Db.Table("courses").Find(&courses)
	for i := range courses {
		course_cnt[courses[i].CourseID] = courses[i].Cap
	}

	//var data []types.SCourse
	//Db.Table("s_courses").Find(&data)
	//for i := range data {
	//	course_cnt[data[i].CourseID] -= 1
	//	err := Rdb.HSetNX(ctx, data[i].UserID, data[i].CourseID, 0).Err()
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	for k, v := range course_cnt {
		fmt.Println(k, v)
		Rdb.Set(ctx, k+"cnt", v, 0)
	}

	fmt.Println("open redis success")
}

func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8mb4"
	dsn := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8mb4"}, "")
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	Db.AutoMigrate(&types.Member{})
	Db.AutoMigrate(&types.Course{})
	Db.AutoMigrate(&types.SCourse{})
	if err != nil {
		fmt.Println("open database fail")
		return
	}
	fmt.Println("connnect success")
}
