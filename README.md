# Course Selection

## 项目结构如下

```bash
.
├── auth
│   └── auth.go
├── database
│   └── gorm.go
├── go.mod
├── go.sum
├── main.go
├── README.md
├── router
│   └── router.go
├── test.py
└── types
    └── types.go
```

### main 模块

在 main.go 文件中，负责启动程序。

### types 模块

其中定义了所使用的要的数据结构和常量。

### router 模块

绑定地址，为不同的地址绑定不同的响应函数。

### database 模块

初始化数据库连接，其中包含了一个数据库连接变量 Db ，可以通过这个变量对数据库进行操作。

### auth 模块

登录验证模块

### test.py 文件

向主机发送不同请求进行测试。

### 数据库相关

#### 数据库创建

```sql
CREATE TABLE TMember (
  UserID BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
  Nickname VARCHAR(32) NOT NULL,
  Username VARCHAR(32) NOT NULL,
  UserType int NOT NULL,
  Password VARCHAR(32) NOT NULL
) DEFAULT CHARSET UTF8;
```

#### 数据库添加管理员数据

```sql
insert into TMember (Nickname,Username,UserType,Password) values ("root","JudgeAdmin",1,"JudgePassword2022")
```