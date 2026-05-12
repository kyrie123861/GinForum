package controllers

type MyCode int64

// 在go语言中没有enum的关键字来实现枚举
// 但是可以直接使用const + iota来实现
// 好处是保持语言简单，且利用现有机制就能实现enum而不用引入enum这个关键字

const (
	CodeSuccess         MyCode = 1000 //MyCode = 1000 + iota
	CodeInvalidParams   MyCode = 1001
	CodeUserExist       MyCode = 1002
	CodeUserNotExist    MyCode = 1003
	CodeInvalidPassword MyCode = 1004
	CodeServerBusy      MyCode = 1005

	CodeInvalidToken      MyCode = 1006
	CodeInvalidAuthFormat MyCode = 1007
	CodeNotLogin          MyCode = 1008
)

var msgFlags = map[MyCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken:      "无效的Token",
	CodeInvalidAuthFormat: "认证格式有误",
	CodeNotLogin:          "未登录",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if !ok {
		msg = msgFlags[CodeServerBusy]
	}
	return msg
}
