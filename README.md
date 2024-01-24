# 项目文档

- 后台框架：[beego](https://github.com/beego/bee) 

## 1. 基本介绍

### 1.1 项目介绍

> ticket是一个基于beego开发的两票管理系统。

## 2. 使用说明

```
golang版本 >= v1.18
```

### 2.1 server端

```bash
# 使用 go.mod

# 安装go依赖包
go list (go mod tidy)

# 本地调试
go run main.go run
```

### 2.2 swagger自动化API文档

#### 2.2.1 安装 bee

````
go get -u github.com/beego/bee/v2
go install github.com/beego/bee/v2
````

#### 2.2.2 生成API文档

````
bee -downdoc=true -gendoc=true
./iam run
````
执行上面的命令后，目录下会出现swagger文件夹，登录http://127.0.0.1/api/ticket/v1/swagger，即可查看swagger文档


## 3. 技术选型

- 后端：用`Beego`快速搭建基础restful风格API，`Beego`是一个go语言编写的Web框架。
- 数据库：采用`MySql`(8.0.20)版本，使用`orm`实现对数据库的基本操作。
- API文档：使用`Swagger`构建自动化文档。
- 配置文件：使用`beego`自带的`conf/app.conf`格式的配置文件。
- 日志：使用`beego/logs`实现日志记录。
