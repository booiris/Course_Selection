import requests

url = 'http://180.184.74.221:2000/api/v1/auth/login'
# 以字典的形式构造数据
data = {
    'Username': 'booiris',
    'Passwd': '12asd3'
}
# 与 get 请求一样，r 为响应对象
r = requests.post(url, data=data)
# 查看响应结果
print(r.text)