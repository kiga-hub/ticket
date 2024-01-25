package enums

type JsonResultCode int

const (
	JRCodeSucc         JsonResultCode = 200
	JRCodeParamError                  = 504 
	JRCodeRequestError                = 400 
)
