package controllers

import (
	"fmt"
	"path"
	"strconv"
	"ticket/enums"
	"ticket/models"
	"ticket/utils"

	"github.com/astaxie/beego/orm"
)

// UserController .
type UserController struct {
	BaseController
}

// Prepare -
func (c *UserController) Prepare() {
	// prepare
	c.BaseController.Prepare()
}

// Login -
func (c *UserController) Login() {
	if !c.checkParams([]string{
		"Name",
		"Pwd",
	}) {
		c.jsonResult(enums.JRCodeParamError, "Invalid parameter, please refresh and try again", "")
	}
	c.param["Pwd"] = utils.String2md5(c.param["Pwd"].(string))
	user, err := models.UserOneByCheck(c.param["Name"].(string), c.param["Pwd"].(string))

	if user == nil || err != nil {
		c.jsonResult(enums.JRCodeRequestError, "Username or password is incorrect", err)
	}
	// expire time
	user.Expires = 12 * 3600
	user.Token = c.getTokenStr(user)
	c.setUser2Cache(user)

	c.curUser = *user

	// save session
	// c.setUser2Session(user)

	c.jsonResult(enums.JRCodeSucc, "Login successful", map[string]interface{}{
		"Token":   user.Token,
		"UsserId": user.UserId,
		"Name":    user.Name,
	})
}

// Logout -
func (c *UserController) Logout() {
	user, err := models.UserOne(c.curUser.UserId)

	if user == nil || err != nil {
		c.jsonResult(enums.JRCodeParamError, "Username or password is incorrect", err)
	}

	utils.DelCache(c.curUser.Token)
	//c.curUser = *user

	c.jsonResult(enums.JRCodeSucc, "Logout successful", map[string]interface{}{
		"Name": user.Name,
	})
}

// Register User -
func (c *UserController) Register() {
	if !c.checkParams([]string{"Name", "Pwd", "Mobile", "Address", "Major", "Role",
		"Department", "Gender"}) {
		c.jsonResult(enums.JRCodeParamError, "Invalid parameter, please refresh and try again", "")
	}
	var user models.User
	u, _ := models.UserOneByName(c.param["Name"].(string))
	if u != nil {
		c.jsonResult(enums.JRCodeRequestError, "Username has already been registered", map[string]interface{}{
			"Name":   u.Name,
			"Moblie": u.Mobile,
		})
	}
	m, _ := models.UserOneByMobile(c.param["Mobile"].(string))
	if m != nil {
		c.jsonResult(enums.JRCodeRequestError, "The phone number has already been registered", map[string]interface{}{
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
		c.jsonResult(enums.JRCodeRequestError, "User registration failed", err)
	}
	user.Expires = 12 * 3600
	user.Token = c.getTokenStr(&user)
	c.setUser2Cache(&user)

	c.curUser = user
	c.jsonResult(enums.JRCodeSucc, "Registration successful", map[string]interface{}{
		"Token": user.Token,
		"Name":  user.Name,
	})
}

// UserSave .
func (c *UserController) UserSave() {
	name := c.curUser.Name
	oM, err := models.UserOneByName(name)

	if oM == nil || err != nil {
		c.jsonResult(enums.JRCodeParamError, "Data is invalid, please refresh and try again", "")
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
		c.jsonResult(enums.JRCodeRequestError, "Save failed", name)
	} else {
		oM.Token = c.curUser.Token
		oM.Expires = c.curUser.Expires
		c.setUser2Cache(oM)
		c.jsonResult(enums.JRCodeSucc, "Save successful", name)
	}

}

// PageList -
func (c *UserController) PageList() {
	var params models.UserQueryParam

	if !c.checkParams([]string{"Page", "Limit"}) {
		c.jsonResult(enums.JRCodeParamError, "Invalid parameter, please refresh and try again", "")
	}
	limit, _ := strconv.Atoi(c.param["Limit"].(string))
	params.Limit = limit
	page, _ := strconv.Atoi(c.param["Page"].(string))
	params.Offset = limit * (page - 1)

	if _, ok := c.param["Url"]; ok {
		params.SearchLike = c.param["Url"].(string)
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
			"Major":      row.Major,
			"Role":       row.Role,
			"Department": row.Department,
			"Gender":     row.Gender,
			"IdCard":     row.IdCard,
		}
		rows = append(rows, one)
	}
	result["rows"] = rows
	c.jsonResult(enums.JRCodeSucc, "Successful", result)
}

// Donwload -
func (c *UserController) Donwload() {
	_, err := models.UserOneByName(c.curUser.Name)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, "Donwload failed", "")
	}
	filepath := c.curUser.Url
	name := path.Base(filepath)
	c.Ctx.Output.Download(filepath, name)
	c.jsonResult(enums.JRCodeSucc, "Download successful", map[string]interface{}{
		"Url":  c.curUser.Url,
		"Name": name,
	})
}

// Upload .
func (c *UserController) Upload() {

	f, h, _ := c.GetFile("file")
	// ext := path.Ext(h.Filename)
	//
	// var AllowExtMap map[string]bool = map[string]bool{
	// 	".wav": true,
	// 	".mp3": true,
	// }
	// if _, ok := AllowExtMap[ext]; !ok {
	// 	c.jsonResult(enums.JRCodeRequestError, "The suffix does not match, the file upload failed, please re-upload", "")
	// 	return
	// }

	fileName := h.Filename
	fpath := c.dataRoot + c.curUser.UserId + "_" + fileName
	defer f.Close()
	err := c.SaveToFile("file", fpath)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, fmt.Sprintf("%v", err), "")
	}

	oM, err := models.UserOneByName(c.curUser.Name)
	oM.Url = fpath
	o := orm.NewOrm()
	if _, err := o.Update(oM); err != nil {
		c.jsonResult(enums.JRCodeRequestError, "File upload failed", map[string]interface{}{
			"Result": "Unable to match user",
			"Error":  err,
		})
	}

	c.jsonResult(enums.JRCodeSucc, "File upload successful", map[string]interface{}{
		"Result": "successful",
	})
}

// TextTip .
func (c *UserController) TextTip() {
	result := map[string]interface{}{}
	result["TextTip"] = `Text prompt: Employee number + text prompt, for example: Employee number 001: text prompt`
	usertip := fmt.Sprintf("Employee number%s: text prompt", c.curUser.Name)
	c.jsonResult(enums.JRCodeSucc, usertip, result)
}
