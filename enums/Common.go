package enums

type JsonResultCode int

const (
	JRCodeSucc         JsonResultCode = 200
	JRCodeFailed                      = 600
	JRCodeError                       = 500
	JRCodeWarn                        = 501
	JRCode302                         = 302 
	JRCode401                         = 401 
	JRCodeParamError                  = 504 
	JRCodeRequestError                = 400 
)
