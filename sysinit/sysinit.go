package sysinit

import (
	"ticket/utils"
	"time"

	"github.com/astaxie/beego"
)

func init() {
	time.Local, _ = time.LoadLocation("Asia/Shanghai")
	// Session on
	beego.BConfig.WebConfig.Session.SessionOn = true
	// init logs
	//utils.InitLogs()
	// init cache
	utils.InitCache()
	// init database
	InitDatabase()
}
