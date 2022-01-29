import requests

url = 'http://180.184.74.221:2000/api/v1/auth/login'
# 以字典的形式构造数据
data = {
    'Username': 'JudgeAdmin',
    'Password': 'JudgePassword2022'
}
# 与 get 请求一样，r 为响应对象
r = requests.post(url, data=data)
# 查看响应结果
print(r.text)

# CREATE TABLE TMember (
#   UserID BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
#   Nickname VARCHAR(32) NOT NULL,
#   Username VARCHAR(32) NOT NULL,
#   UserType int NOT NULL,
#   Password VARCHAR(32) NOT NULL
# ) DEFAULT CHARSET UTF8;

# insert into members (Nickname,Username,User_Type,Password) values ("root","JudgeAdmin",1,"JudgePassword2022")
