package models

import "ticket/enums"

// JsonResult Base struuct used to return ajax requests
type JsonResult struct {
	Code enums.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Data interface{}          `json:"data"`
}

// BaseQueryParam struuct used for querying
type BaseQueryParam struct {
	Sort   string `json:"Sort"`
	Order  string `json:"Order"`
	Offset int    `json:"Offset"`
	Limit  int    `json:"Limit"`
}
