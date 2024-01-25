
# server address: http://127.0.0.1:8081

=================login=================
#接口URL
http://127.0.0.1:8081/api/ticket/home/login
#请求方式
POST
#请求 Content-Type：application/json

#请求参数
{"Name":"admin1","Pwd":"123456"}  用户名 和密码

#返回值
{
  "code": 200,
  "msg": "登录成功",
  "data": {
    "Name": "admin1",
    "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODgwMjM0MjIsInN1YiI6IntcIlVzZXJJZFwiOlwiMDBlYmUxMDEtYzJhOS00MjIxLTk0NzItZWQ4ODBkYzY4NmNjXCJ9In0.U8qsFNQKStyKTZU4-GH2XUa7wULyb3Fpg-Sg-1BVTR4",
    "UsserId": "00ebe101-c2a9-4221-9472-ed880dc686cc"
  }
}

=================退出=================
#接口URL
http://127.0.0.1:8081/api/ticket/home/logout
#请求方式
POST
#请求 Content-Type：application/json

#请求参数
{}

#返回值
{
  "code": 200,
  "msg": "退出成功",
  "data": {
    "Name": "admin1"
  }
}
#重复退出 返回错误值
{
  "code": 200,
  "msg": "用户名或者密码错误",
  "data": {}
}


用户退出后无法调用后续接口
#返回值
{
  "code": 200,
  "msg": "未登录",
  "data": {
    "error": "用户信息未找到"
  }
}
=================用户注册=================
#用户注册
#接口URL
http://127.0.0.1:8081/api/ticket/home/register

接口预留



=================保存用户信息=================
#保存用户信息
#接口URL
http://127.0.0.1:8081/api/ticket/home/save
#请求方式
POST
#请求 Content-Type：application/json
#请求参数

####下列参数任选  ！！！！！
{"Gender":1，"Mobile":"1232323213"，"Role":1,"Department":"1212","Post":"213213","Safetybelt":"3213121","Safetyhelmet":"23213"}




=================上传文件=================
#上传文件

#接口URL
http://127.0.0.1:8081/api/ticket/home/upload

#请求方式
POST

#请求 Content-Type: multipart/form-data

#(file 选择文件)
Key      VALUE
file     "1.wav"

#返回值
{
  "code": 200,
  "msg": "上传文件成功",
  "data": {
    "Result": "上传文件成功"
  }
}

#上传文件不是.wav  或者 .mp3 
#返回参数错误
{
  "code": 200,
  "msg": "后缀名不符合,上传文件失败，请重新上传",
  "data": ""
}

=================载入所有员工信息=================
#载入所有员工信息
#接口URL
http://127.0.0.1:8081/api/ticket/info/pagelist

POST
#请求 Content-Type: application/json

请求参数
{"Page":"1","Limit":"10"}  
