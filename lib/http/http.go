package http

import (
    "bytes"
    "encoding/json"
    "io"
    "io/ioutil"
    "net/http"
    "time"
)

// 发送GET请求
// url:请求地址
// response:请求返回的内容
func GetHttp(url string) (response string) {
    client := http.Client{Timeout: 5 * time.Second}
    resp, error := client.Get(url)
    if error != nil {
        panic(error)
    }
    defer resp.Body.Close()
    var buffer [512]byte
    result := bytes.NewBuffer(nil)
    for {
        n, err := resp.Body.Read(buffer[0:])
        result.Write(buffer[0:n])
        if err != nil && err == io.EOF {
            break
        } else if err != nil {
            panic(err)
        }
    }

    response = result.String()
    return
}

// 发送POST请求
// url:请求地址，data:POST请求提交的数据
// content:请求放回的内容
func PostHttp(url string, data interface{}) (content string) {
    jsonStr, _ := json.Marshal(data)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Add("content-type", "application/json")
    if err != nil {
        panic(err)
    }
    defer req.Body.Close()

    client := &http.Client{Timeout: 5 * time.Second}
    resp, error := client.Do(req)
    if error != nil {
        panic(error)
    }
    defer resp.Body.Close()

    result, _ := ioutil.ReadAll(resp.Body)
    content = string(result)
    return
}