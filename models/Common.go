package models

import "Two-Card/enums"

// JsonResult 用于返回ajax请求的基类
type JsonResult struct {
	Code enums.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Data interface{}          `json:"data"`
}

// BaseQueryParam 用于查询的类
type BaseQueryParam struct {
	Sort   string `json:"Sort"`
	Order  string `json:"Order"`
	Offset int    `json:"Offset"`
	Limit  int    `json:"Limit"`
}
