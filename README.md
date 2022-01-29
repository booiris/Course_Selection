# Course Selection

## 项目结构如下

```bash
.
├── controller
│   ├── auth.go
│   ├── course_teacher.go
│   ├── member.go
│   └── student.go
├── database
│   └── gorm.go
├── go.mod
├── go.sum
├── main.go
├── README.md
├── router
│   └── router.go
└── types
    └── types.go
```

### controller 模块

地址的响应函数的实现

### database 模块

初始化数据库连接，其中包含了一个数据库连接变量 Db ，可以通过这个变量对数据库进行操作。

### main 模块

在 main.go 文件中，负责启动程序。

### router 模块

绑定地址，为不同的地址绑定不同的响应函数。

### types 模块

其中定义了所使用的要的数据结构和常量。
