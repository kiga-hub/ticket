package sms

import (
	"math/rand"
	"strconv"
	"time"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
	"errors"
	"Two-Card/utils"
)

func GetRandVerifyCode() string{                                                  
	return strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(899999) + 100000)
} 

//通过秒滴发送一条验证码
func SendSmsByMaiodi(mobile, code string) error{
	client := &http.Client{}
	timestamp := time.Now().Format("20060101150405")
	md5str := "9040b94d711045aeab97d1c3c88f7c21e356d47977db418b92020fc983953f44"+timestamp
	sig := utils.String2md5(md5str)
	post_arg := url.Values{
		"accountSid": {"9040b94d711045aeab97d1c3c88f7c21"},
		"templateid": {"1372951165"}, 
		"param":{code}, 
		"to":{mobile}, 
		"timestamp":{timestamp}, 
		"sig":{sig}, 
	}
	// return nil
	//str_code := "accountSid=9040b94d711045aeab97d1c3c88f7c21&templateid=1372951165&param=816940&to=13261831802&timestamp=16046092341611647&sig=29b2cee11d507e05bdc698c3dbae02e3"
	//读取Api数据
	urls := "https://api.miaodiyun.com/20150822/industrySMS/sendSMS"
	req, err := http.NewRequest("POST", urls, strings.NewReader(post_arg.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	if response["respCode"].(string) == "00000" {
		return nil
	}
	return errors.New(response["respDesc"].(string))
}
func SendSmsByMiaodiMessage(mobile, code string) error{
	client := &http.Client{}
	timestamp := time.Now().Format("20060101150405")
	md5str := "9040b94d711045aeab97d1c3c88f7c21e356d47977db418b92020fc983953f44"+timestamp
	sig := utils.String2md5(md5str)
	post_arg := url.Values{
		"accountSid": {"9040b94d711045aeab97d1c3c88f7c21"},
		"templateid": {"248569"}, 
		"param":{code}, 
		//"smsContent":{code},
		"to":{mobile}, 
		"timestamp":{timestamp}, 
		"sig":{sig}, 
	}
	utils.LogDebug(post_arg)
	// return nil
	//str_code := "accountSid=9040b94d711045aeab97d1c3c88f7c21&templateid=1372951165&param=816940&to=13261831802&timestamp=16046092341611647&sig=29b2cee11d507e05bdc698c3dbae02e3"
	//读取Api数据
	urls := "https://api.miaodiyun.com/20150822/industrySMS/sendSMS"
	req, err := http.NewRequest("POST", urls, strings.NewReader(post_arg.Encode()))
	if err != nil {
		utils.LogError(err)
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		utils.LogError(err)
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	utils.LogDebug(response)
	if response["respCode"].(string) == "00000" {
		return nil
	}
	return errors.New(response["respDesc"].(string))
}