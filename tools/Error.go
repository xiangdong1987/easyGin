package tools

const (
	// a == 1 (iota在每个const开头被重设为0)
	ErrorNotFound = 1000 + iota
	ErrorInvalidParams
	ErrorDB
	ErrorDBConfig
	ErrorNetErr
	ErrorImgCaptchaCheckErr
)

var ErrorList = map[int]string{
	ErrorNotFound:           "数据不存在",
	ErrorInvalidParams:      "无效参数",
	ErrorDB:                 "DB请求出错",
	ErrorDBConfig:           "DB配置出错，请检查配置",
	ErrorNetErr:             "网络异常",
	ErrorImgCaptchaCheckErr: "图片验证码不正确",
}

type ReturnError struct {
	Code int    `json:"code" `
	Msg  string `json:"msg" `
}

func (ge ReturnError) Error() string {
	return ge.Msg
}

func (ge ReturnError) Instance(code int) ReturnError {
	ge.Code = code
	ge.Msg = ErrorList[code]
	return ge
}

func (ge ReturnError) Custom(code int, msg string) ReturnError {
	ge.Code = code
	ge.Msg = msg
	return ge
}
