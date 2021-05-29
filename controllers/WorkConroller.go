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

//WorkController .
type WorkController struct {
	BaseController
}

//Prepare .
func (c *WorkController) Prepare() {
	c.BaseController.Prepare()
	c.checkAuthor()
}

//WorkPageList .
func (c *WorkController) WorkPageList() {
	// 参数检测
	if !c.checkParams([]string{"Page", "Limit"}) {
		c.jsonResult(enums.JRCodeRequestError, "参数无效，请刷新后重试", "")
	}
	var params models.WorkQueryParam
	limit, _ := strconv.Atoi(c.param["Limit"].(string))
	params.Limit = limit
	page, _ := strconv.Atoi(c.param["Page"].(string))
	params.Offset = limit * (page - 1)

	workfile := c.BaseController.workRoot
	result := make(map[string]interface{})
	var rows []interface{}

	result["path"] = workfile
	da := utils.GetFolders(workfile)
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
func (c *WorkController) PageList() {
	// 参数检测
	if !c.checkParams([]string{"Page", "Limit"}) {
		c.jsonResult(enums.JRCodeParamError, "参数无效，请刷新后重试", "")
	}
	var params models.WorkQueryParam
	limit, _ := strconv.Atoi(c.param["Limit"].(string))
	params.Limit = limit
	page, _ := strconv.Atoi(c.param["Page"].(string))
	params.Offset = limit * (page - 1)

	if _, ok := c.param["WorkcardtaskUrl"]; ok {
		params.SearchLike = c.param["WorkcardtaskUrl"].(string)
	}

	if _, ok := c.param["Catalog"]; ok {
		params.CatalogLike = c.param["Catalog"].(string)
	}
	params.IdLike = c.curUser.UserId
	data, total := models.WorkPageList(&params, &c.curUser)

	result := make(map[string]interface{})
	result["total"] = total

	var rows []interface{}
	for _, row := range data {
		one := map[string]interface{}{
			"Id":              row.Id,
			"UserId":          row.UserId,
			"WorkcardtaskUrl": row.WorkcardtaskUrl,
			"GetworkTime":     row.GetworkTime,
			"WorkstartTime":   row.WorkstartTime,
			"WorkendTime":     row.WorkendTime,
			"Worksigner":      row.Worksigner,
			"Workmanager":     row.Workmanager,
			"Catalog":         row.Catalog,
		}
		rows = append(rows, one)
	}
	result["rows"] = rows
	c.jsonResult(enums.JRCodeSucc, "成功", result)
}

//ShowWorkFile 工作票下载
func (c *WorkController) ShowWorkFile() {
	// 参数检测
	if !c.checkParams([]string{"WorkId"}) {
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

	work, _ := models.WorkOne(c.param["WorkId"].(string))
	filepath := work.WorkcardtaskUrl
	name := path.Base(filepath)

	//fmt.Println(name)
	c.Ctx.Output.Download(filepath, name)
}

//WorkSynchronize 同步信息 下载两票内容
func (c *WorkController) WorkSynchronize() {

	oM, err := models.UserOneByName(c.curUser.Name)
	if err != nil {
		c.jsonResult(enums.JRCodeRequestError, "同步失败", "")
	}
	voiceprint := oM.VoicePrint
	if voiceprint == "未注册" {
		c.jsonResult(enums.JRCodeRequestError, "未注册声纹，请注册声纹", "")
	}

	work, total := models.WorkByUserID(c.curUser.UserId)
	result := make(map[string]interface{})
	var rows []interface{}
	if total == 0 {
		c.jsonResult(enums.JRCodeRequestError, "无法找到工作票", map[string]interface{}{
			"UserId": oM.UserId,
		})
	}

	for _, row := range work {
		filepath := row.WorkcardtaskUrl
		name := path.Base(filepath)

		utils.CopyFile(c.tmpRoot+"/"+name, filepath)

		one := map[string]interface{}{
			"Id":          row.Id,
			"UserId":      row.UserId,
			"WorkpathUrl": row.WorkcardtaskUrl,
			"FileName":    name,
		}
		rows = append(rows, one)
	}

	result["rows"] = rows

	f, err := os.Open(c.tmpRoot)
	if err != nil {
		fmt.Println(err)
		c.jsonResult(enums.JRCodeSucc, "工作票同步失败", result)
	}
	defer f.Close()
	var files = []*os.File{f}
	zipfile := c.homeRoot + "/" + c.curUser.UserId + "_workfile.zip"
	err = utils.Compress(files, zipfile)
	if err != nil {
		fmt.Println(err)
		c.jsonResult(enums.JRCodeSucc, "工作票同步失败", result)
	}
	tmpfile := c.BaseController.tmpRoot + "/"
	fileInfoList, err := ioutil.ReadDir(tmpfile)

	fmt.Println(len(fileInfoList))
	for i := range fileInfoList {
		fmt.Println(fileInfoList[i].Name()) //打印当前文件或目录下的文件或目录名
		os.Chmod(fileInfoList[i].Name(), 0777)
		os.Remove(tmpfile + fileInfoList[i].Name())
	}
	downfile := c.curUser.UserId + "_workfile.zip"
	c.Ctx.Output.Download(zipfile, downfile)

	// quit := make(chan bool)
	// go func() {
	// 	//for {
	// 		select {
	// 		case <- quit:
	// 			fmt.Println("2e2e")
	// 		default:
	// 			// Do other stuff
	// 			c.Ctx.Output.Download(c.tmpRoot+".zip","workfile.zip")
	// 		}
	// 	//}
	// }()
	// close(quit)

	// result["rows"]=rows
	// result["DownloadFile"]="workfile.zip"
	// c.jsonResult(enums.JRCodeSucc,"工作票同步成功",result)

}
