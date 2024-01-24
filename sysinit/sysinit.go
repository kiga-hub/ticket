package sysinit

import (
	"ticket/utils"
	"time"

	"github.com/astaxie/beego"
)

func init() {
	time.Local, _ = time.LoadLocation("Asia/Shanghai")
	//启用Session
	beego.BConfig.WebConfig.Session.SessionOn = true
	//初始化日志
	//utils.InitLogs()
	//初始化缓存
	utils.InitCache()
	//初始化数据库
	InitDatabase()
}
