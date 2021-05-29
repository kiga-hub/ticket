/*
 * @Author: jinquan
 * @Date: 2019-11-08 22:45:27
 * @LastEditTime: 2020-04-03 10:28:13
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \Two-Card\enums\Common.go
 */
package enums

type JsonResultCode int

const (
	JRCodeSucc   JsonResultCode = 200
	JRCodeFailed                = 600
	JRCodeError                 = 500
	JRCodeWarn                  = 501
	JRCode302                   = 302 //跳转至地址
	JRCode401                   = 401 //未授权访问
	JRCodeParamError			= 504 //请求参数无效
	JRCodeRequestError			= 400 //请求无效
)

