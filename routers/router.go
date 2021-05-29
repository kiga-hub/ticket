package routers

import (
	"Two-Card/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	ns :=
		beego.NewNamespace("/api/twocard",
			/*
			 *	请求校验
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
				beego.NSRouter("login", &controllers.UserController{}, "Post:Login"),                      //登录
				beego.NSRouter("logout", &controllers.UserController{}, "Post:Logout"),                    //退出
				beego.NSRouter("register", &controllers.UserController{}, "Post:Register"),                //user注册
				beego.NSRouter("save", &controllers.UserController{}, "Post:UserSave"),                    //保存用户信息
				beego.NSRouter("voiceprintregister", &controllers.UserController{}, "Post:VoiceRegister"), //声纹注册
				beego.NSRouter("upload", &controllers.UserController{}, "Post:VoiceUpload"),               //声纹注册 上传文件
			),
			beego.NSNamespace("information",
				beego.NSRouter("pagelist", &controllers.UserController{}, "Post:PageList"), //显示全部员工信息
			),
			beego.NSNamespace("synchronize",
				beego.NSRouter("usersynchronize", &controllers.UserController{}, "Post:UserSynchronize"),   //同步个人信息
				beego.NSRouter("voicesynchronize", &controllers.UserController{}, "Post:VoiceSynchronize"), //声纹同步 下载声纹音频
			),
			beego.NSNamespace("tip",
				beego.NSRouter("texttip", &controllers.VoiceController{}, "Post:TextTip"), //文字提示
			),
			beego.NSNamespace("workrecord",
				beego.NSRouter("upload", &controllers.VoiceController{}, "Post:WorkRecordUpload"),         //工作票录音上传
				beego.NSRouter("ordinaryupload", &controllers.VoiceController{}, "Post:WorkOrdinaryRecordUpload"),         //工作票普通录音上传
				beego.NSRouter("commandwordupload", &controllers.VoiceController{}, "Post:WorkCommandeWordRecordUpload"),         //工作票命令词录音上传
				beego.NSRouter("voicepagelist", &controllers.VoiceController{}, "Post:WorkVoicePageList"), //工作票录音List
			),
			beego.NSNamespace("operaterecord",
				beego.NSRouter("upload", &controllers.VoiceController{}, "Post:OperateRecordUpload"),         //操作票录音上传
				beego.NSRouter("ordinaryupload", &controllers.VoiceController{}, "Post:OperateOrdinaryRecordUpload"),         //操作票普通录音上传
				beego.NSRouter("commandwordupload", &controllers.VoiceController{}, "Post:OperateCommandWordRecordUpload"),         //操作票命令词录音上传
				beego.NSRouter("voicepagelist", &controllers.VoiceController{}, "Post:OperateVoicePageList"), //操作票录音List
			),
			beego.NSNamespace("work",
				beego.NSRouter("root", &controllers.WorkController{}, "Post:WorkPageList"),               //工作票根目录列表
				beego.NSRouter("show", &controllers.WorkController{}, "Post:ShowWorkFile"),               //工作票下载
				beego.NSRouter("worksynchronize", &controllers.WorkController{}, "Post:WorkSynchronize"), // 工作票同步信息 下载两票内容
				beego.NSRouter("workpagelist", &controllers.WorkController{}, "Post:PageList"),           //工作票列表
			),
			beego.NSNamespace("operate",
				beego.NSRouter("root", &controllers.OperateController{}, "Post:OperatePageList"), //操作票根目录列表
				beego.NSRouter("show", &controllers.OperateController{}, "Post:ShowOperateFile"),
				beego.NSRouter("operatesynchronize", &controllers.OperateController{}, "Post:OperateSynchronize"), // 操作票同步
				beego.NSRouter("operatepagelist", &controllers.OperateController{}, "Post:PageList"),              //操作票列表

			),
		)
	beego.AddNamespace(ns)
}
