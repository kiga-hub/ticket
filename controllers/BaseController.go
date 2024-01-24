package controllers

import (
	"fmt"
	"time"

	"encoding/json"
	"ticket/enums"
	"ticket/models"
	"ticket/utils"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
)

// BaseController .
type BaseController struct {
	beego.Controller
	controllerName     string
	actionName         string
	curUser            models.User
	param              map[string]interface{}
	dataRoot           string
	homeRoot           string
	workvoiceRoot      string
	operatevoiceRoot   string
	workRoot           string
	operateRoot        string
	videoRoot          string
	tmpRoot            string
	workordinary       string
	workcommandword    string
	operateordinary    string
	operatecommandword string
}

/**
 * @description: -
 */
func (c *BaseController) Prepare() {
	// Get the controller name and execution function name corresponding to the requested API
	c.controllerName, c.actionName = c.GetControllerAndAction()

	// get data from redis. set user info
	c.adapterCacheUserInfo()

	// -
	c.jsonRequest()

	// -
	c.dataRoot = beego.AppConfig.String("dataroot")

	// -
	c.homeRoot = beego.AppConfig.String("homeroot")
	c.workvoiceRoot = beego.AppConfig.String("workvoiceroot")
	c.operatevoiceRoot = beego.AppConfig.String("operatevoiceroot")
	c.workRoot = beego.AppConfig.String("workroot")
	c.operateRoot = beego.AppConfig.String("operateroot")
	c.videoRoot = beego.AppConfig.String("videoroot")
	c.tmpRoot = beego.AppConfig.String("tmproot")
	c.workordinary = beego.AppConfig.String("workordinaryroot")
	c.workcommandword = beego.AppConfig.String("workcommandwordroot")
	c.operateordinary = beego.AppConfig.String("operateordinaryroot")
	c.operatecommandword = beego.AppConfig.String("operatecommandwordroot")

}

/**
 * @description: generate Token string
 * @param {user info}
 * @return: Token string
 */
func (c *BaseController) getTokenStr(user *models.User) string {
	// key
	keyInfo := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()"

	// map2key
	info := map[string]interface{}{}
	info["UserId"] = user.UserId
	dataByte, _ := json.Marshal(info)
	var dataStr = string(dataByte)

	// Token expire time
	expiresAt := time.Now().Unix() + user.Expires

	// using Claim to save json
	data := jwt.StandardClaims{Subject: dataStr, ExpiresAt: expiresAt}
	tokenInfo := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	// generate token string
	tokenStr, _ := tokenInfo.SignedString([]byte(keyInfo))
	//utils.LogDebug(fmt.Sprintf("uid:%s,expires:%d,token:%s", user.Id, user.Expires,tokenStr))
	return tokenStr
}

/**
 * @description: vertification Token
 */
func (c *BaseController) checkTokenStr(tokenStr string) bool {
	// convert token to object type
	tokenInfo, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
		return token, nil
	})

	// verify
	err := tokenInfo.Claims.Valid()
	if err != nil {
		return false
	}

	finToken := tokenInfo.Claims.(jwt.MapClaims)

	// verify
	if !finToken.VerifyExpiresAt(time.Now().Unix(), true) {
		return false
	}

	// -
	sub := finToken["sub"].(string)
	var tokenData map[string]interface{}
	json.Unmarshal([]byte(sub), &tokenData)
	if _, ok := tokenData["UserId"]; !ok {
		return false
	}

	// -
	if c.curUser.UserId != tokenData["UserId"].(string) {
		return false
	}
	return true
}

/**
 * @description: checkLogin 判断用户是否登录，在BaseController.Prepare()后执行
 * @param {type}
 * @return:
 */
func (c *BaseController) checkLogin() {
	// 判断用户
	if c.curUser.UserId == "" {
		c.jsonResult(enums.JRCodeSucc, "未登录", map[string]string{
			"error": "用户信息未找到",
		})
	}
}

// checkLogin判断用户是否有权访问某地址，无权则会跳转到错误页面
// 一定要在BaseController.Prepare()后执行
// 会调用checkLogin
// 传入的参数为忽略权限控制的Action
func (c *BaseController) checkAuthor(ignores ...string) {
	//如果Action在忽略列表里，则直接通用
	for _, ignore := range ignores {
		if ignore == c.actionName {
			return
		}
	}

	// 登录检查
	c.checkLogin()

	//Token检查
	token := c.Ctx.Input.Header("Token")
	if !c.checkTokenStr(c.Ctx.Input.Header("Token")) {
		utils.DelCache(token)
		c.jsonResult(enums.JRCodeSucc, "未登录", map[string]string{
			"error": "token错误",
		})
	}
}

/**
 * @description: 从sesssion里取用户信息, 没有使用
 */
func (c *BaseController) adapterSessionUserInfo() {
	a := c.GetSession("user")
	if a != nil {
		c.curUser = a.(models.User)
	}
}

/**
 * @description: 从cache里取用户信息
 */
func (c *BaseController) adapterCacheUserInfo() {
	token := c.Ctx.Input.Header("Token")
	if token == "" {
		return
	}
	user := new(models.User)
	if err := utils.GetCache(token, user); err != nil {
		return
	}
	if user != nil {
		c.curUser = *user
	}
}

// /**
//   - @description: 获取用户信息（包括资源UrlFor）保存至redis
//     */
func (c *BaseController) setUser2Cache(user *models.User) error {
	// 获取这个用户能获取到的所有资源列表
	// resourceList := models.ResourceTreeGridByUserId(user.Id, 1000)
	// for _, item := range resourceList {
	// 	user.ResourceUrlForList = append(user.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	// }
	// orm.NewOrm().LoadRelated(user, "Organizations")
	if err := utils.SetCache(user.Token, user, int(user.Expires)); err != nil {
		utils.LogError(err)
	}
	return nil
}

/**
 * @description: 返回结果JSON
 */
func (c *BaseController) jsonResult(code enums.JsonResultCode, msg string, data interface{}) {
	r := &models.JsonResult{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	c.Data["json"] = r
	if code == 500 {
		jsonStr, _ := json.Marshal(data)
		utils.LogError(fmt.Sprintf("%s - %s - %s", string(c.Ctx.Input.RequestBody), msg, jsonStr))
	} else if code == 501 {
		jsonStr, _ := json.Marshal(data)
		utils.LogWarning(fmt.Sprintf("%s - %s - %s", string(c.Ctx.Input.RequestBody), msg, jsonStr))
	} else if code == 600 {
		jsonStr, _ := json.Marshal(data)
		utils.LogInfo(fmt.Sprintf("%s - %s - %s", string(c.Ctx.Input.RequestBody), msg, jsonStr))
	}
	c.ServeJSON()
	c.StopRun()
}

/**
 * @description: 请求参数JSON获取
 */
func (c *BaseController) jsonRequest() {
	json.Unmarshal(c.Ctx.Input.RequestBody, &c.param)
	// c.param["URI"] = fmt.Sprintf("%s[%s.%s]", c.Ctx.Request.RequestURI, c.controllerName, c.actionName)
}

/**
 * @description: 检查必选参数
 */
func (c *BaseController) checkParams(request []string) bool {
	for _, item := range request {
		if _, ok := c.param[item]; !ok {
			return false
		}
	}
	return true
}
