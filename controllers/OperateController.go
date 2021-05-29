package controllers

import (
	"Two-Card/enums"
	"Two-Card/models"
	"Two-Card/utils"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

//OperateController .
type OperateController struct {
	BaseController
}

//Prepare .
func (c *OperateController) Prepare() {
	c.BaseController.Prepare()
	c.checkAuthor()
}

//OperatePageList .
func (c *OperateController) OperatePageList() {
	// 参数检测
	if !c.checkParams([]string{"Page", "Limit"}) {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}
	var params models.OperateQueryParam
	limit, _ := strconv.Atoi(c.param["Limit"].(string))
	params.Limit = limit
	page, _ := strconv.Atoi(c.param["Page"].(string))
	params.Offset = limit * (page - 1)

	operate := c.BaseController.operateRoot
	result := make(map[string]interface{})
	var rows []interface{}

	result["path"] = operate
	da := utils.GetFolders(operate)
	for _, i := range da {
		one := map[string]interface{}{
			"path": i,
		}
		rows = append(rows, one)
	}
	result["rows"] = rows
	c.jsonResult(enums.JRCodeSucc, "成功", result)
}

//PageList .
func (c *OperateController) PageList() {
	var params models.OperateQueryParam
	// 参数检测
	if !c.checkParams([]string{"Page", "Limit"}) {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}
	limit, _ := strconv.Atoi(c.param["Limit"].(string))
	params.Limit = limit
	page, _ := strconv.Atoi(c.param["Page"].(string))
	params.Offset = limit * (page - 1)

	if _, ok := c.param["OperatecardtaskUrl"]; ok {
		params.SearchLike = c.param["OperatecardtaskUrl"].(string)
	}
	if _, ok := c.param["Catalog"]; ok {
		params.CatalogLike = c.param["Catalog"].(string)
	}
	params.IdLike = c.curUser.UserId
	data, total := models.OperatePageList(&params, &c.curUser)

	result := make(map[string]interface{})
	result["total"] = total

	var rows []interface{}
	for _, row := range data {
		one := map[string]interface{}{
			"Id":                 row.Id,
			"UserId":             row.UserId,
			"OperatecardtaskUrl": row.OperatecardtaskUrl,
			"GetOperateTime":     row.GetoperateTime,
			"OperatestartTime":   row.OperatestartTime,
			"OperateendTime":     row.OperateendTime,
			"Operatesigner":      row.Operatesigner,
			"Operatemanager":     row.Operatemanager,
			"Command":row.Command,
			"Catalog":            row.Catalog,
		}
		rows = append(rows, one)
	}
	result["rows"] = rows
	c.jsonResult(enums.JRCodeSucc, "成功", result)
}

//ShowOperateFile 操作票下载
func (c *OperateController) ShowOperateFile() {
	// 参数检测
	if !c.checkParams([]string{"OperateId"}) {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}

	oM, err := models.UserOneByName(c.curUser.Name)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, "同步失败", "")
	}
	voiceprint := oM.VoicePrint
	if voiceprint == "未注册" {
		c.jsonResult(enums.JRCodeRequestError, "未注册声纹，请注册声纹", "")
	}

	operate, _ := models.OperateOne(c.param["OperateId"].(string))
	filepath := operate.OperatecardtaskUrl
	name := path.Base(filepath)

	fmt.Println(name)
	c.Ctx.Output.Download(filepath, name)
}

//OperateSynchronize 同步信息 下载两票内容
func (c *OperateController) OperateSynchronize() {
	oM, err := models.UserOneByName(c.curUser.Name)
	if err != nil {
		c.jsonResult(enums.JRCodeParamError, "同步失败", "")
	}
	voiceprint := oM.VoicePrint
	if voiceprint == "未注册" {
		c.jsonResult(enums.JRCodeRequestError, "未注册声纹，请注册声纹", "")
	}
	operate, total := models.OperateByUserID(oM.UserId)
	if total == 0 {
		c.jsonResult(enums.JRCodeRequestError, "无法找到操作票", map[string]interface{}{
			"UserId": oM.UserId,
		})
	}

	result := make(map[string]interface{})
	var rows []interface{}

	for _, row := range operate {
		filepath := row.OperatecardtaskUrl
		name := path.Base(filepath)
		//c.Ctx.Output.Download(filepath,name)

		utils.CopyFile(c.tmpRoot+"/"+name, filepath)

		one := map[string]interface{}{
			"Id":          row.Id,
			"UserId":      row.UserId,
			"WorkpathUrl": row.OperatecardtaskUrl,
			"FileName":    name,
		}
		rows = append(rows, one)
	}
	result["rows"] = rows

	f, err := os.Open(c.tmpRoot)
	if err != nil {
		fmt.Println(err)
		c.jsonResult(enums.JRCodeRequestError, "操作票同步失败", result)
	}
	defer f.Close()
	var files = []*os.File{f}
	zipfile := c.homeRoot + "/" + c.curUser.UserId + "_operate.zip"
	err = utils.Compress(files, zipfile)
	if err != nil {
		fmt.Println(err)
		c.jsonResult(enums.JRCodeRequestError, "操作票同步失败", result)
	}
	tmpfile := c.BaseController.tmpRoot + "/"
	fileInfoList, err := ioutil.ReadDir(tmpfile)

	fmt.Println(len(fileInfoList))
	for i := range fileInfoList {
		fmt.Println(fileInfoList[i].Name()) //打印当前文件或目录下的文件或目录名
		os.Chmod(fileInfoList[i].Name(), 0777)
		os.Remove(tmpfile + fileInfoList[i].Name())
	}
	downfile := c.curUser.UserId + "_operatefile.zip"
	c.Ctx.Output.Download(zipfile, downfile)

}
