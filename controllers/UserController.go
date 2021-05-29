package controllers

import (
	"Two-Card/enums"
	"Two-Card/models"
	"Two-Card/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"path"
	"strconv"
)

//UserController .
type UserController struct {
	BaseController
}

//Prepare 进入处理程序之前先验证 .
func (c *UserController) Prepare() {
	//先执行
	c.BaseController.Prepare()
}

//Login .登录
func (c *UserController) Login() {
	if !c.checkParams([]string{
		"Name",
		"Pwd",
	}) {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}
	c.param["Pwd"] = utils.String2md5(c.param["Pwd"].(string))
	user, err := models.UserOneByCheck(c.param["Name"].(string), c.param["Pwd"].(string))

	if user == nil || err != nil {
		c.jsonResult(enums.JRCodeRequestError, "用户名或者密码错误", err)
	}
	// 超时时间
	user.Expires = 12 * 3600
	user.Token = c.getTokenStr(user)
	c.setUser2Cache(user)

	c.curUser = *user

	//保存用户信息到缓存
	// c.setUser2Session(user)

	c.jsonResult(enums.JRCodeSucc, "登录成功", map[string]interface{}{
		"Token":   user.Token,
		"UsserId": user.UserId,
		"Name":    user.Name,
	})
}

//Logout .退出
func (c *UserController) Logout() {
	user, err := models.UserOne(c.curUser.UserId)

	if user == nil || err != nil {
		c.jsonResult(enums.JRCodeParamError, "用户名或者密码错误", err)
	}

	utils.DelCache(c.curUser.Token)
	//c.curUser = *user

	c.jsonResult(enums.JRCodeSucc, "退出成功", map[string]interface{}{
		"Name": user.Name,
	})
}

//Register User . 员工 注册
func (c *UserController) Register() {
	if !c.checkParams([]string{"Name", "Pwd", "Mobile", "Address", "Major", "Role",
		"Department", "Gender"}) {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}
	var user models.User
	u, _ := models.UserOneByName(c.param["Name"].(string))
	if u != nil {
		c.jsonResult(enums.JRCodeRequestError, "用户名已经注册", map[string]interface{}{
			"Name":   u.Name,
			"Moblie": u.Mobile,
		})
	}
	m, _ := models.UserOneByMobile(c.param["Mobile"].(string))
	if m != nil {
		c.jsonResult(enums.JRCodeRequestError, "手机号已经注册", map[string]interface{}{
			"Name":   u.Name,
			"Moblie": u.Mobile,
		})
	}
	user.UserId = utils.Uuid()
	if _, ok := c.param["Name"]; ok {
		user.Name = c.param["Name"].(string)
	}
	user.Address = c.param["Address"].(string)
	user.Mobile = c.param["Mobile"].(string)
	user.Pwd = utils.String2md5(c.param["Pwd"].(string))
	user.Major = c.param["Major"].(string)
	user.Role = int(c.param["Role"].(float64))
	user.Department = c.param["Department"].(string)
	user.Gender = int(c.param["Gender"].(float64))
	user.IdCard = utils.Uuid()

	if _, err := orm.NewOrm().Insert(&user); err != nil {
		c.jsonResult(enums.JRCodeRequestError, "用户注册失败", err)
	}
	user.Expires = 12 * 3600
	user.Token = c.getTokenStr(&user)
	c.setUser2Cache(&user)

	c.curUser = user
	c.jsonResult(enums.JRCodeSucc, "注册成功", map[string]interface{}{
		"Token": user.Token,
		"Name":  user.Name,
	})
}

//UserSave .
func (c *UserController) UserSave() {
	name := c.curUser.Name
	oM, err := models.UserOneByName(name)

	if oM == nil || err != nil {
		c.jsonResult(enums.JRCodeParamError, "数据无效，请刷新后重试", "")
	}
	if _, ok := c.param["Gender"]; ok {
		oM.Gender = int(c.param["Gender"].(float64))
	}
	if _, ok := c.param["Mobile"]; ok {
		oM.Mobile = c.param["Mobile"].(string)
	}
	if _, ok := c.param["Role"]; ok {
		oM.Role = int(c.param["Role"].(float64))
	}
	if _, ok := c.param["Department"]; ok {
		oM.Department = c.param["Department"].(string)
	}
	if _, ok := c.param["Post"]; ok {
		oM.Post = c.param["Post"].(string)
	}
	if _, ok := c.param["Safetybelt"]; ok {
		oM.Safetybelt = c.param["Safetybelt"].(string)
	}
	if _, ok := c.param["Safetyhelmet"]; ok {
		oM.Safetyhelmet = c.param["Safetyhelmet"].(string)
	}

	o := orm.NewOrm()
	if _, err := o.Update(oM); err != nil {
		c.jsonResult(enums.JRCodeRequestError, "保存失败", name)
	} else {
		oM.Token = c.curUser.Token
		oM.Expires = c.curUser.Expires
		c.setUser2Cache(oM)
		c.jsonResult(enums.JRCodeSucc, "保存成功", name)
	}

}

//VoiceRegister . 员工声纹注册
func (c *UserController) VoiceRegister() {
	oM, err := models.UserOneByName(c.curUser.Name)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, "同步失败", "")
	}

	voiceprint := oM.VoicePrint
	if voiceprint == "已注册" {
		c.jsonResult(enums.JRCodeRequestError, "已注册声纹", "")
	}

	c.jsonResult(enums.JRCodeSucc, "请上传声纹文件，注册声纹", "")
}

//PageList . 显示全部员工信息
func (c *UserController) PageList() {
	var params models.UserQueryParam
	// 参数检测
	if !c.checkParams([]string{"Page", "Limit"}) {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}
	limit, _ := strconv.Atoi(c.param["Limit"].(string))
	params.Limit = limit
	page, _ := strconv.Atoi(c.param["Page"].(string))
	params.Offset = limit * (page - 1)

	if _, ok := c.param["VoicePrint"]; ok {
		params.SearchLike = c.param["VoicePrint"].(string)
	}

	data, total := models.UserPageList(&params, &c.curUser)

	result := make(map[string]interface{})
	result["total"] = total
	var rows []interface{}
	//O := orm.NewOrm()
	for _, row := range data {
		// O.LoadRelated(row, "Organizations")
		//O.LoadRelated(row, "Name")
		one := map[string]interface{}{
			"UserId":     row.UserId,
			"Name":       row.Name,
			"Mobile":     row.Mobile,
			"Address":    row.Address,
			"VoicePrint": row.VoicePrint,
			"Major":      row.Major,
			"Role":       row.Role,
			"Department": row.Department,
			"Gender":     row.Gender,
			"IdCard":     row.IdCard,
			"VoiceUrl":   row.VoiceUrl,
		}
		rows = append(rows, one)
	}
	result["rows"] = rows
	c.jsonResult(enums.JRCodeSucc, "成功", result)
}

//UserSynchronize .同步信息
func (c *UserController) UserSynchronize() {
	oM, err := models.UserOneByName(c.curUser.Name)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, "同步失败", "")
	}
	// voiceprint := oM.VoicePrint
	// if voiceprint == "未注册" {
	// 	c.jsonResult(enums.JRCodeRequestError, "未注册声纹，请注册声纹", "")
	// }
	result := map[string]interface{}{
		"UserId":       oM.UserId,
		"Name":         oM.Name,
		"Mobile":       oM.Mobile,
		"Address":      oM.Address,
		"VoicePrint":   oM.VoicePrint,
		"Major":        oM.Major,
		"Role":         oM.Role,
		"Department":   oM.Department,
		"Gender":       oM.Gender,
		"IdCard":       oM.IdCard,
		"VoiceUrl":     oM.VoiceUrl,
		"Post":         oM.Post,
		"SafetyBelt":   oM.Safetybelt,
		"SafetyHelmet": oM.Safetyhelmet,
	}
	c.jsonResult(enums.JRCodeSucc, "同步个人信息,成功", result)

}

//VoiceSynchronize .声纹同步
func (c *UserController) VoiceSynchronize() {
	oM, err := models.UserOneByName(c.curUser.Name)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, "同步失败", "")
	}
	voiceprint := oM.VoicePrint
	if voiceprint == "未注册" {
		c.jsonResult(enums.JRCodeRequestError, "未注册声纹，请注册声纹", "")
	}
	filepath := c.curUser.VoiceUrl
	name := path.Base(filepath)
	c.Ctx.Output.Download(filepath, name)
	c.jsonResult(enums.JRCodeSucc, "声纹同步成功", map[string]interface{}{
		"Url":  c.curUser.VoiceUrl,
		"Name": name,
	})
}

//VoiceUpload .
func (c *UserController) VoiceUpload() {

	f, h, _ := c.GetFile("file") //获取上传的文件
	// ext := path.Ext(h.Filename)
	// //验证后缀名是否符合要求
	// var AllowExtMap map[string]bool = map[string]bool{
	// 	".wav": true,
	// 	".mp3": true,
	// }
	// if _, ok := AllowExtMap[ext]; !ok {
	// 	c.jsonResult(enums.JRCodeRequestError, "后缀名不符合,上传文件失败，请重新上传", "")
	// 	return
	// }

	fileName := h.Filename
	fpath := c.dataRoot + c.curUser.UserId+"_"+ fileName
	defer f.Close() //关闭上传的文件
	err := c.SaveToFile("file", fpath)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, fmt.Sprintf("%v", err), "")
	}

	oM, err := models.UserOneByName(c.curUser.Name)
	oM.VoiceUrl = fpath
	oM.VoicePrint = "已注册"
	o := orm.NewOrm()
	if _, err := o.Update(oM); err != nil {
		c.jsonResult(enums.JRCodeRequestError, "上传文件失败", map[string]interface{}{
			"Result": "无法匹配用户",
			"Error":err,
		})
	}

	c.jsonResult(enums.JRCodeSucc, "上传文件成功", map[string]interface{}{
		"Result": "上传文件成功",
	})
}
