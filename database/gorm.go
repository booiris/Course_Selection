package database

import (
	"context"
	"course_selection/types"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//数据库配置
const (
	userName = "root"
	password = "bytedancecamp"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "main"
)

var Db *gorm.DB
var Rdb *redis.Client

func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8mb4"
	dsn := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8mb4&parseTime=True"}, "")
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
	sqlDb.SetConnMaxIdleTime(10)
	// 最大连接数
	sqlDb.SetMaxOpenConns(100)
	// 连接复用连接时间
	sqlDb.SetConnMaxLifetime(time.Hour)
	fmt.Println("connnect success")
}

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
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

func CheckCourse(course_id string) bool {
	ctx := context.Background()
	is_exist := Rdb.Get(ctx, "course"+course_id)
	if is_exist.Err() == redis.Nil {
		var res types.Course
		Db.Model(types.Course{}).Where("course_id=?", course_id).Find(&res)
		if res == (types.Course{}) {
			return false
		} else {
			Rdb.Set(ctx, "course"+course_id, 0, 0)
			return true
		}
	} else {
		return true
	}
}

func FindUserType(user_id string) types.UserType {
	ctx := context.Background()
	redis_res := Rdb.Get(ctx, "usertype"+user_id)
	if redis_res.Err() == redis.Nil {
		var res types.Member
		Db.Model(types.Member{}).Unscoped().Where("user_id=?", user_id).Find(&res)
		var ans int
		if res == (types.Member{}) {
			ans = 0
		} else if res.Deleted.Valid {
			ans = -1
		} else {
			ans = int(res.UserType)
		}
		Rdb.Set(ctx, "usertype"+user_id, ans, 0)
		return types.UserType(ans)
	} else {
		fmt.Println("1233222")
		res, _ := strconv.Atoi(redis_res.Val())
		return types.UserType(res)
	}
}

func SyncMysql() {
	ctx := context.Background()
	cnt := 0
	for {
		query := Rdb.BLPop(ctx, 0, "Sync_mysql")
		temp := strings.Split(query.Val()[1], ",")
		create_data := types.SCourse{
			UserID:   temp[0],
			CourseID: temp[1],
		}
		Db.Table("s_courses").Select("user_id", "course_id").Create(&create_data)
		cnt++
		if cnt > 100 {
			cnt = 0
			time.Sleep(200 * time.Millisecond)
		}
	}
}
