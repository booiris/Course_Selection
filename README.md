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
│   └── database.go
├── go.mod
├── go.sum
├── main.go
├── README.md
├── router
│   └── router.go
├── test
│   ├── analysis.py
│   └── test.go
└── types
    └── types.go
```

### controller 模块

地址的响应函数的实现

### database 模块

初始化数据库连接，其中包含了一个数据库连接变量 Db 和一个redis连接变量 Rdb，可以通过这些变量对数据库和redis进行操作。

### main 模块

在 main.go 文件中，负责启动程序。

### router 模块

绑定地址，为不同的地址绑定不同的响应函数。

### test 模块

自己编写的压测函数，用于检验数据一致性。

### types 模块

其中定义了所使用的要的数据结构和常量。

## 使用方法

安装依赖

```
go install
```

程序运行
```
go run .
```