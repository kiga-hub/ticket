package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

type JsonResultCode int

// JsonResult 用于返回ajax请求的基类
type JsonResult struct {
	Code JsonResultCode `json:"code"`
	Msg  string         `json:"msg"`
	Data interface{}    `json:"data"`
}

//GetToken 平台登录Token.
func GetToken(host string, info interface{}, platform string) (string, error) {

	send, _ := json.Marshal(info)
	resp, err := http.Post(host+platform+"/home/login", "application/json", strings.NewReader(string(send)))
	fmt.Println(platform)
	if err != nil {
		fmt.Println("Post Error")
		return "", err
	}
	defer resp.Body.Close()
	var token string
	body, _ := ioutil.ReadAll(resp.Body)
	var m JsonResult
	k := string(body)
	if err := json.Unmarshal([]byte(k), &m); err == nil {
		token = m.Data.(map[string]interface{})["Token"].(string)
	}
	return token, nil
}

//CreateRequestHelper 请求结果.
func CreateRequestHelper(host string, Token string, page map[string]interface{}, platform string, router string) (JsonResult, int, error) {
	client := &http.Client{}
	var dataString string
	var m JsonResult
	dataType, _ := json.Marshal(page)
	dataString = string(dataType)
	postaddress := host + platform + "/" + router
	req, err := http.NewRequest(
		"POST", postaddress,
		strings.NewReader(dataString),
	)
	if err != nil {
		fmt.Println("handle error", err)
		return m, int(m.Code), err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")
	req.Header.Set("Token", Token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return m, int(m.Code), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("handle error", err)
		return m, int(m.Code), err
	}

	k := string(body)
	if err := json.Unmarshal([]byte(k), &m); err != nil {
		fmt.Println(m.Code)
		return m, int(m.Code), err
	}
	//orderid := m.Data.(map[string]interface{})["rows"].([]interface{})[0].(map[string]interface{})["Id"]
	return m, int(m.Code), nil
}

//DownLoadFile voice file.
func DownLoadFile(url string, filename string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

//GetAddr .
func GetAddr() string {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return "Erorr"
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}
