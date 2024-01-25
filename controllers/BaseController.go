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
	controllerName string
	actionName     string
	curUser        models.User
	param          map[string]interface{}
	dataRoot       string
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
 * @description: get user info from cache
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
//   - @description: Get user information (including resource UrlFor) and save it to Redis
//     */
func (c *BaseController) setUser2Cache(user *models.User) error {
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
 * @description: return json result
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
 * @description: json request
 */
func (c *BaseController) jsonRequest() {
	json.Unmarshal(c.Ctx.Input.RequestBody, &c.param)
	// c.param["URI"] = fmt.Sprintf("%s[%s.%s]", c.Ctx.Request.RequestURI, c.controllerName, c.actionName)
}

/**
 * @description: Check mandatory parameters
 */
func (c *BaseController) checkParams(request []string) bool {
	for _, item := range request {
		if _, ok := c.param[item]; !ok {
			return false
		}
	}
	return true
}
