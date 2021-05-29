#服务器地址 http://39.100.93.194:8081

=================登录=================
#接口URL
http://39.100.93.194:8081/api/twocard/home/login
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
http://39.100.93.194:8081/api/twocard/home/logout
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
=================员工注册=================
#员工注册
#接口URL
http://39.100.93.194:8081/api/twocard/home/register

接口预留



=================保存用户信息=================
#保存用户信息
#接口URL
http://39.100.93.194:8081/api/twocard/home/save
#请求方式
POST
#请求 Content-Type：application/json
#请求参数

####下列参数任选  ！！！！！
{"Gender":1，"Mobile":"1232323213"，"Role":1,"Department":"1212","Post":"213213","Safetybelt":"3213121","Safetyhelmet":"23213"}




=================声纹注册=================
#声纹注册
#接口URL
http://39.100.93.194:8081/api/twocard/home/voiceprintregister

#请求方式
POST
#请求 Content-Type：application/json

#请求参数
{}

返回值
声纹未注册提示 "请上传声纹文件，注册声纹"
声纹已注册提示 "已注册声纹"


=================上传声纹录音文件=================
#上传声纹录音文件

#接口URL
http://39.100.93.194:8081/api/twocard/home/upload

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
http://39.100.93.194:8081/api/twocard/information/pagelist

POST
#请求 Content-Type: application/json

请求参数
{"Page":"1","Limit":"10"}  // 可以 添加参数 "VoicePrint":"已注册" 查看已注册声纹员工列表


=================同步个人信息=================
#同步个人信息
#接口URL
http://39.100.93.194:8081/api/twocard/synchronize/usersynchronize

POST
#请求 Content-Type: application/json

请求参数{}

#读取数据库有延迟  提示 "同步失败"  重新同步个人信息即可

#未注册声纹 提示 "未注册声纹，请注册声纹"


#返回值
{
  "code": 200,
  "msg": "同步个人信息,成功",
  "data": {
    "Address": "西二旗",
    "Department": "运维中心",
    "Gender": "M",
    "IdCard": "00001",
    "Major": "计算机",
    "Mobile": "17600001234",
    "Name": "admin1",
    "Role": 1,
    "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc",
    "VoicePrint": "已注册",
    "VoiceUrl": "/home/liangpiao/userfile/voiceprint.wav"
  }
}

=================同步声纹信息=================
#同步声纹信息
#接口URL
http://39.100.93.194:8081/api/twocard/synchronize/voicesynchronize

POST
#请求 Content-Type: application/json

请求参数{}

#读取数据库有延迟  提示 "同步失败"  重新同步个人信息即可

#未注册声纹 提示 "未注册声纹，请注册声纹"

#开始下载声纹音频


=================文字提示=================
#文字提示

#接口URL
http://39.100.93.194:8081/api/twocard/tip/texttip

#请求方式
POST

#请求 Content-Type：application/json

请求参数
{}

#返回值
{
  "code": 200,
  "msg": "员工编号admin1:文字提示", //当前用户名字
  "data": {
    "TextTip": "我国《民事诉讼法》的规定，证据有当事人的陈述、证人证言、视听资料等，录音资料则属于视听资料证据的一种"
  }
}
=================工作票根目录=================
#工作票列表

#接口URL
http://39.100.93.194:8081/api/twocard/work/root

#请求方式
POST

#请求 Content-Type：application/json

#请求参数
{"Page":"1","Limit":"10"}

#返回值
{
  "code": 200,
  "msg": "成功",
  "data": {
    "path": "/home/liangpiao/workfile/",
    "rows": [
      {
        "path": "华控工作票01"
      },
      {
        "path": "华控工作票02"
      },
      {
        "path": "华控工作票03"
      },
      {
        "path": "华控工作票04"
      }
    ]
  }
}


=================工作票列表=================
#工作票列表

#接口URL
http://39.100.93.194:8081/api/twocard/work/workpagelist

#请求方式
POST

#请求 Content-Type：application/json

#请求参数
{"Page":"1","Limit":"10"}

#返回值
{
  "code": 200,
  "msg": "成功",
  "data": {
    "rows": [
      {
        "GetworkTime": "2020-04-30T19:23:53-07:00",
        "Id": "3",
        "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc",
        "WorkcardtaskUrl": "/home/liangpiao/workfile/work2.docx",
        "WorkendTime": "2020-04-30T19:24:00-07:00",
        "Workmanager": "admin1",
        "Worksigner": "admin4",
        "WorkstartTime": "2020-04-30T19:23:57-07:00"
      }
    ],
    "total":1
  }
}

=================工作票同步下载文档=================
#工作票同步下载文档

#接口URL
http://39.100.93.194:8081/api/twocard/work/worksynchronize

#请求方式
POST

#请求 Content-Type：application/json

#请求参数
{}

下载工作票
下载格式为  usereid_workfile.zip
打包下载与该用户相关的所有工作票内容

=================操作票根目录=================
#操作票列表

#接口URL
http://39.100.93.194:8081/api/twocard/operate/root

#请求方式
POST

#请求 Content-Type：application/json

#请求参数
{"Page":"1","Limit":"10"}

#返回值
{
  "code": 200,
  "msg": "成功",
  "data": {
    "path": "/home/liangpiao/operatefile/",
    "rows": [
      {
        "path": "华控操作票01"
      },
      {
        "path": "华控操作票02"
      },
      {
        "path": "华控操作票03"
      },
      {
        "path": "华控操作票04"
      }
    ]
  }
}



=================操作票列表=================
#操作票列表

#接口URL
http://39.100.93.194:8081/api/twocard/operate/operatepagelist

#请求方式
POST

#请求 Content-Type：application/json

#请求参数
{"Page":"1","Limit":"10"}

#返回值
{
  "code": 200,
  "msg": "成功",
  "data": {
    "rows": [
      {
        "GetOperateTime": "2020-05-07T17:54:01-07:00",
        "Id": "1",
        "OperatecardtaskUrl": "/home/liangpiao/operatefile/operate1.docx",
        "OperateendTime": "2020-05-07T17:54:06-07:00",
        "Operatemanager": "admin1",
        "Operatesigner": "admin4",
        "OperatestartTime": "2020-05-07T17:54:03-07:00",
        "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc"
      }
    ],
    "total": 1
  }
}


=================操作票同步下载文档=================
#操作票同步下载文档

#接口URL
http://39.100.93.194:8081/api/twocard/operate/operatesynchronize

#请求方式
POST

#请求 Content-Type：application/json

#请求参数
{}

下载操作票
下载格式为  usereid_operate.zip
打包下载与该用户相关的所有操作票内容

=================工作票录音上传=================
#工作票录音上传
#上传文件的同时  上传json 参数   

#接口URL
http://39.100.93.194:8081/api/twocard/workrecord/upload?WorkId=1&CardOption=This is Testq2121e

#Params 选项添加 KEY            VALUE
                WorkId         1
                CardOption     工作票内容

#请求方式
POST

#请求 Content-Type: multipart/form-data

#(file 选择文件)
Key      VALUE
file     "workrecord.wav"


{
  "code": 200,
  "msg": "上传成功",
  "data": {
    "Id": "525aa3fc-b807-4928-8447-de449dda7b13",
    "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc"
  }
}
=================工作票普通录音上传=================
#工作票录音上传
#上传文件的同时  上传json 参数   

#接口URL
http://39.100.93.194:8081/api/twocard/workrecord/ordinaryupload?WorkId=1&CardOption=普通录音上传

#Params 选项添加 KEY            VALUE
                WorkId         1
                CardOption     工作票内容

#请求方式
POST

#请求 Content-Type: multipart/form-data

#(file 选择文件)
Key      VALUE
file     "workrecord.wav"


{
  "code": 200,
  "msg": "上传成功",
  "data": {
    "Id": "525aa3fc-b807-4928-8447-de449dda7b13",
    "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc"
  }
}
=================工作票命令词录音上传=================
#工作票录音上传
#上传文件的同时  上传json 参数   

#接口URL
http://39.100.93.194:8081/api/twocard/workrecord/commandwordupload?WorkId=1&CardOption=命令词录音上传

#Params 选项添加 KEY            VALUE
                WorkId         1
                CardOption     工作票内容

#请求方式
POST

#请求 Content-Type: multipart/form-data

#(file 选择文件)
Key      VALUE
file     "workrecord.wav"


{
  "code": 200,
  "msg": "上传成功",
  "data": {
    "Id": "525aa3fc-b807-4928-8447-de449dda7b13",
    "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc"
  }
}
=================操作票录音上传=================
#操作票录音上传
#上传文件的同时  上传json 参数   

#接口URL
http://39.100.93.194:8081/api/twocard/operaterecord/upload?OperateId=1&CardOption=This is Testq2121e

#Params 选项添加 KEY            VALUE
                OperateId         1
                CardOption     操作票内容

#请求方式
POST

#请求 Content-Type: multipart/form-data

#(file 选择文件)
Key      VALUE
file     "workrecord.wav"


{
  "code": 200,
  "msg": "上传成功",
  "data": {
    "Id": "525aa3fc-b807-4928-8447-de449dda7b13",
    "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc"
  }
}
=================操作票普通录音上传=================
#操作票录音上传
#上传文件的同时  上传json 参数   

#接口URL
http://39.100.93.194:8081/api/twocard/operaterecord/ordinaryupload?OperateId=1&CardOption=操作票普通录音

#Params 选项添加 KEY            VALUE
                OperateId         1
                CardOption     操作票内容

#请求方式
POST

#请求 Content-Type: multipart/form-data

#(file 选择文件)
Key      VALUE
file     "workrecord.wav"


{
  "code": 200,
  "msg": "上传成功",
  "data": {
    "Id": "525aa3fc-b807-4928-8447-de449dda7b13",
    "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc"
  }
}
=================操作票命令词录音上传=================
#操作票录音上传
#上传文件的同时  上传json 参数   

#接口URL
http://39.100.93.194:8081/api/twocard/operaterecord/commandwordupload?OperateId=1&CardOption=操作票命令词录音

#Params 选项添加 KEY            VALUE
                OperateId         1
                CardOption     操作票内容

#请求方式
POST

#请求 Content-Type: multipart/form-data

#(file 选择文件)
Key      VALUE
file     "workrecord.wav"


{
  "code": 200,
  "msg": "上传成功",
  "data": {
    "Id": "525aa3fc-b807-4928-8447-de449dda7b13",
    "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc"
  }
}
=================工作票录音列表=================
#工作票录音列表 

#接口URL
http://39.100.93.194:8081/api/twocard/workrecord/voicepagelist

#请求方式
POST
#请求 Content-Type：application/json

#请求参数
{"Page":"1","Limit":"10"}

{
  "code": 200,
  "msg": "成功",
  "data": {
    "rows": [
      {
        "CardOption": "This is Test",
        "CardType": 0,
        "Id": "0cc8927e-f4c6-440b-a1c7-6cf79d744791",
        "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc",
        "VoiceUrl": "/home/liangpiao/voicefile/workrecord.wav",
        "WorkId": "1"
      },
    ],
    "total": 1
  }
}
=================操作票录音列表=================
#操作票录音列表 

#接口URL
http://39.100.93.194:8081/api/twocard/operaterecord/voicepagelist

#请求方式
POST
#请求 Content-Type：application/json

#请求参数
{"Page":"1","Limit":"10"}

{
  "code": 200,
  "msg": "成功",
  "data": {
    "rows": [
      {
        "CardOption": "feanflwfowafnwafnafwl",
        "CardType": 1,
        "Id": "525aa3fc-b807-4928-8447-de449dda7b13",
        "OperateId": "1",
        "UserId": "00ebe101-c2a9-4221-9472-ed880dc686cc",
        "VoiceUrl": "/home/liangpiao/voicefile/workrecord.wav"
      }
    ],
    "total": 1
  }
}

End
// TODO 
fix bug

