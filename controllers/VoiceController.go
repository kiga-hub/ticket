package controllers

import (
	"Two-Card/enums"
	"Two-Card/models"
	"Two-Card/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"os"
	"path"
	"strconv"
)

//VoiceController .
type VoiceController struct {
	BaseController
}

// BaseQueryParam 用于查询的类
type BaseQueryParam struct {
	Sort   string `json:"Sort"`
	Order  string `json:"Order"`
	Offset int    `json:"Offset"`
	Limit  int    `json:"Limit"`
}

//Prepare 进入处理程序之前先验证 .
func (c *VoiceController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	c.checkAuthor()
}

//WorkVoicePageList .
func (c *VoiceController) WorkVoicePageList() {
	// 参数检测
	if !c.checkParams([]string{"Page", "Limit"}) {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}
	var params models.VoiceQueryParam
	limit, _ := strconv.Atoi(c.param["Limit"].(string))
	params.Limit = limit
	page, _ := strconv.Atoi(c.param["Page"].(string))
	params.Offset = limit * (page - 1)

	// if _, ok := c.param["CardType"]; ok {
	// 	params.TypeLike =int(c.param["CardType"].(float64))
	// }
	params.TypeLike = 0
	data, total := models.WorkVoicePageList(&params)

	result := make(map[string]interface{})
	result["total"] = total

	var rows []interface{}
	for _, row := range data {
		one := map[string]interface{}{
			"Id":         row.Id,
			"UserId":     row.UserId,
			"WorkId":     row.WorkId,
			"VoiceUrl":   row.VoiceUrl,
			"CardType":   row.CardType,
			"CardOption": row.CardOption,
		}
		rows = append(rows, one)
	}
	result["rows"] = rows
	c.jsonResult(enums.JRCodeSucc, "成功", result)
}

//OperateVoicePageList .
func (c *VoiceController) OperateVoicePageList() {
	// 参数检测
	if !c.checkParams([]string{"Page", "Limit"}) {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}
	var params models.VoiceQueryParam
	limit, _ := strconv.Atoi(c.param["Limit"].(string))
	params.Limit = limit
	page, _ := strconv.Atoi(c.param["Page"].(string))
	params.Offset = limit * (page - 1)

	// if _, ok := c.param["CardType"]; ok {
	// 	params.TypeLike =int(c.param["CardType"].(float64))
	// }
	params.TypeLike = 1
	data, total := models.WorkVoicePageList(&params)

	result := make(map[string]interface{})
	result["total"] = total

	var rows []interface{}
	for _, row := range data {
		one := map[string]interface{}{
			"Id":         row.Id,
			"UserId":     row.UserId,
			"OperateId":  row.OperateId,
			"VoiceUrl":   row.VoiceUrl,
			"CardType":   row.CardType,
			"CardOption": row.CardOption,
		}
		rows = append(rows, one)
	}
	result["rows"] = rows
	c.jsonResult(enums.JRCodeSucc, "成功", result)
}

//VoiceDelete .
func (c *VoiceController) VoiceDelete() {
	filepath := c.curUser.VoiceUrl
	var t string
	if exsit := pathExists(filepath); exsit == true {
		if ok := os.Remove(filepath); ok != nil {
			c.jsonResult(enums.JRCodeRequestError, "删除文件失败", "")
		}

		oM, _ := models.UserOneByName(c.curUser.Name)
		t = oM.VoiceUrl
		oM.VoicePrint = "未注册"
		oM.VoiceUrl = ""
		o := orm.NewOrm()
		if _, err := o.Update(oM); err != nil {
			c.jsonResult(enums.JRCodeRequestError, "删除录音文件失败", map[string]interface{}{
				"File_Path":  c.curUser.VoiceUrl,
				"File_Path2": t,
			})
		}
	}
	c.jsonResult(enums.JRCodeSucc, "删除录音文件，成功", map[string]interface{}{
		// "VoicePrint":"请重新注册声纹",
		"VoicePrint": filepath,
		"File_Path":  c.curUser.VoiceUrl,
		"File_Path2": t,
	})
}

//pathExists 文件是否存在方法定义
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func (c *VoiceController) WorkCommandeWordRecordUpload() {
	if c.GetString("CardOption") == "" || c.GetString("WorkId") == "" {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}

	var cardoption string
	c.Ctx.Input.Bind(&cardoption, "CardOption")

	var workid string
	c.Ctx.Input.Bind(&workid, "WorkId")

	f, h, _ := c.GetFile("file") //获取上传的文件
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".wav": true,
		".mp3": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		c.jsonResult(enums.JRCodeRequestError, "后缀名不符合,上传文件失败，请重新上传", "")
		return
	}

	fileName := h.Filename
	fpath := c.workcommandword + fileName
	defer f.Close() //关闭上传的文件
	err := c.SaveToFile("file", fpath)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, fmt.Sprintf("%v", err), "")
	}

	var workrecord models.Voice
	workrecord.Id = utils.Uuid()
	workrecord.UserId = c.curUser.UserId
	workrecord.CardType = 0
	workrecord.VoiceUrl = fpath
	workrecord.CardOption = cardoption
	workrecord.WorkId = workid
	if _, err := orm.NewOrm().Insert(&workrecord); err != nil {
		c.jsonResult(enums.JRCodeRequestError, "上传工作票音频文件失败", err)
	}

	c.jsonResult(enums.JRCodeSucc, "上传成功", map[string]interface{}{
		"Id":     workrecord.Id,
		"UserId": c.curUser.UserId,
	})
}

func (c *VoiceController) WorkOrdinaryRecordUpload() {
	if c.GetString("CardOption") == "" || c.GetString("WorkId") == "" {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}

	var cardoption string
	c.Ctx.Input.Bind(&cardoption, "CardOption")

	var workid string
	c.Ctx.Input.Bind(&workid, "WorkId")

	f, h, _ := c.GetFile("file") //获取上传的文件
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".wav": true,
		".mp3": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		c.jsonResult(enums.JRCodeRequestError, "后缀名不符合,上传文件失败，请重新上传", "")
		return
	}

	fileName := h.Filename
	fpath := c.workordinary + fileName
	defer f.Close() //关闭上传的文件
	err := c.SaveToFile("file", fpath)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, fmt.Sprintf("%v", err), "")
	}

	var workrecord models.Voice
	workrecord.Id = utils.Uuid()
	workrecord.UserId = c.curUser.UserId
	workrecord.CardType = 0
	workrecord.VoiceUrl = fpath
	workrecord.CardOption = cardoption
	workrecord.WorkId = workid
	if _, err := orm.NewOrm().Insert(&workrecord); err != nil {
		c.jsonResult(enums.JRCodeRequestError, "上传工作票音频文件失败", err)
	}

	c.jsonResult(enums.JRCodeSucc, "上传成功", map[string]interface{}{
		"Id":     workrecord.Id,
		"UserId": c.curUser.UserId,
	})
}

//WorkRecordUpload .
func (c *VoiceController) WorkRecordUpload() {

	if c.GetString("CardOption") == "" || c.GetString("WorkId") == "" {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}

	var cardoption string
	c.Ctx.Input.Bind(&cardoption, "CardOption")

	var workid string
	c.Ctx.Input.Bind(&workid, "WorkId")

	f, h, _ := c.GetFile("file") //获取上传的文件
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".wav": true,
		".mp3": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		c.jsonResult(enums.JRCodeRequestError, "后缀名不符合,上传文件失败，请重新上传", "")
		return
	}

	fileName := h.Filename
	fpath := c.workvoiceRoot + fileName
	defer f.Close() //关闭上传的文件
	err := c.SaveToFile("file", fpath)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, fmt.Sprintf("%v", err), "")
	}

	var workrecord models.Voice
	workrecord.Id = utils.Uuid()
	workrecord.UserId = c.curUser.UserId
	workrecord.CardType = 0
	workrecord.VoiceUrl = fpath
	workrecord.CardOption = cardoption
	workrecord.WorkId = workid
	if _, err := orm.NewOrm().Insert(&workrecord); err != nil {
		c.jsonResult(enums.JRCodeRequestError, "上传工作票音频文件失败", err)
	}

	c.jsonResult(enums.JRCodeSucc, "上传成功", map[string]interface{}{
		"Id":     workrecord.Id,
		"UserId": c.curUser.UserId,
	})

}
func (c *VoiceController) OperateCommandWordRecordUpload() {
	if c.GetString("CardOption") == "" || c.GetString("OperateId") == "" {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}

	var cardoption string
	c.Ctx.Input.Bind(&cardoption, "CardOption")

	var operateid string
	c.Ctx.Input.Bind(&operateid, "OperateId")

	f, h, _ := c.GetFile("file") //获取上传的文件
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".wav": true,
		".mp3": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		c.jsonResult(enums.JRCodeRequestError, "后缀名不符合,上传文件失败，请重新上传", "")
		return
	}
	fileName := h.Filename
	fpath := c.operatecommandword + fileName
	defer f.Close() //关闭上传的文件
	err := c.SaveToFile("file", fpath)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, fmt.Sprintf("%v", err), "")
	}
	var operaterecord models.Voice
	operaterecord.Id = utils.Uuid()
	operaterecord.UserId = c.curUser.UserId
	operaterecord.CardType = 1
	operaterecord.VoiceUrl = fpath
	operaterecord.CardOption = cardoption
	operaterecord.OperateId = operateid
	if _, err := orm.NewOrm().Insert(&operaterecord); err != nil {
		c.jsonResult(enums.JRCodeRequestError, "上传操作票音频文件失败", err)
	}

	c.jsonResult(enums.JRCodeSucc, "上传成功", map[string]interface{}{
		"Id":     operaterecord.Id,
		"UserId": c.curUser.UserId,
	})

}
func (c *VoiceController) OperateOrdinaryRecordUpload() {
	if c.GetString("CardOption") == "" || c.GetString("OperateId") == "" {
		c.jsonResult(enums.JRCodeRequestError, "参数无效，请刷新后重试", "")
	}

	var cardoption string
	c.Ctx.Input.Bind(&cardoption, "CardOption")

	var operateid string
	c.Ctx.Input.Bind(&operateid, "OperateId")

	f, h, _ := c.GetFile("file") //获取上传的文件
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".wav": true,
		".mp3": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		c.jsonResult(enums.JRCodeRequestError, "后缀名不符合,上传文件失败，请重新上传", "")
		return
	}
	fileName := h.Filename
	fpath := c.operateordinary + fileName
	defer f.Close() //关闭上传的文件
	err := c.SaveToFile("file", fpath)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, fmt.Sprintf("%v", err), "")
	}
	var operaterecord models.Voice
	operaterecord.Id = utils.Uuid()
	operaterecord.UserId = c.curUser.UserId
	operaterecord.CardType = 1
	operaterecord.VoiceUrl = fpath
	operaterecord.CardOption = cardoption
	operaterecord.OperateId = operateid
	if _, err := orm.NewOrm().Insert(&operaterecord); err != nil {
		c.jsonResult(enums.JRCodeRequestError, "上传操作票音频文件失败", err)
	}

	c.jsonResult(enums.JRCodeSucc, "上传成功", map[string]interface{}{
		"Id":     operaterecord.Id,
		"UserId": c.curUser.UserId,
	})

}

//OperateRecordUpload .
func (c *VoiceController) OperateRecordUpload() {

	if c.GetString("CardOption") == "" || c.GetString("OperateId") == "" {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}

	var cardoption string
	c.Ctx.Input.Bind(&cardoption, "CardOption")

	var operateid string
	c.Ctx.Input.Bind(&operateid, "OperateId")

	f, h, _ := c.GetFile("file") //获取上传的文件
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".wav": true,
		".mp3": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		c.jsonResult(enums.JRCodeRequestError, "后缀名不符合,上传文件失败，请重新上传", "")
		return
	}

	fileName := h.Filename
	fpath := c.operatevoiceRoot + fileName
	defer f.Close() //关闭上传的文件
	err := c.SaveToFile("file", fpath)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, fmt.Sprintf("%v", err), "")
	}

	var operaterecord models.Voice
	operaterecord.Id = utils.Uuid()
	operaterecord.UserId = c.curUser.UserId
	operaterecord.CardType = 1
	operaterecord.VoiceUrl = fpath
	operaterecord.CardOption = cardoption
	operaterecord.OperateId = operateid
	if _, err := orm.NewOrm().Insert(&operaterecord); err != nil {
		c.jsonResult(enums.JRCodeRequestError, "上传操作票音频文件失败", err)
	}

	c.jsonResult(enums.JRCodeSucc, "上传成功", map[string]interface{}{
		"Id":     operaterecord.Id,
		"UserId": c.curUser.UserId,
	})

}

//TextTip .
func (c *VoiceController) TextTip() {
	result := map[string]interface{}{}
	result["TextTip"] = `操作设备人员必须掌握和严格执行国家有关安全政策法规和本企业所制定的安全规章制度，掌握本岗位安全技术操作规程，熟悉设备性能及处理应急情况的措施。操作机电设备必须掌握设备性能、操作方法、安全注意事项，经考试合格后方可独立操作，不许擅自乱开、乱动机电设备以及仪表器具。工作中不准擅自离开工作岗位，未经领导许可不得将设备交给他人操作。运转的机械设备要设安全罩，设备检修完后要做到工完、料净、场地清，安全设施恢复原状。检修时在必要的地方要挂警示牌及安全锁。 起重时必须由一人指挥，一人监护，杜绝多人指挥，紧急情况需起吊时除外。`
	usertip := fmt.Sprintf("员工编号%s:文字提示", c.curUser.Name)
	c.jsonResult(enums.JRCodeSucc, usertip, result)
}

//VoiceDownLoad .
func (c *VoiceController) VoiceDownLoad() {
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
