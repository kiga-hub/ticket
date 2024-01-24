package routers

import (
	"ticket/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	ns :=
		beego.NewNamespace("/api/ticket",
			/*
			 *	request verification
			 */
			beego.NSCond(func(ctx *context.Context) bool {
				return true
				/*
					if ctx.Input.Domain() == "api.beego.me" {
						return true
					}
					return false
				*/
			}),
			beego.NSNamespace("home",
				beego.NSRouter("login", &controllers.UserController{}, "Post:Login"),       // login
				beego.NSRouter("logout", &controllers.UserController{}, "Post:Logout"),     // logout
				beego.NSRouter("register", &controllers.UserController{}, "Post:Register"), // user register
				beego.NSRouter("save", &controllers.UserController{}, "Post:UserSave"),     // user save
				beego.NSRouter("upload", &controllers.UserController{}, "Post:Upload"),     // upload file
				beego.NSRouter("download", &controllers.UserController{}, "Post:Donwload"), // download file

			),
			beego.NSNamespace("info",
				beego.NSRouter("pagelist", &controllers.UserController{}, "Post:PageList"), // Display all employee information
			),
			beego.NSNamespace("tip",
				beego.NSRouter("texttip", &controllers.UserController{}, "Post:TextTip"), // tips
			),
		)
	beego.AddNamespace(ns)
}
