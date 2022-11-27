package errno

func Is(err error, code int) bool {
	return err.Error() == GetMessage(code)
}

func GetMessage(code int) string {
	return codeMessage[code]
}

var codeMessage = map[int]string{
	ERR_CACHE_EXIST:            "缓冲已存在",
	ERR_INTERNAL_SERVER_WRONG:  "服务器出错",
	SUCCESS:                    "请求成功",
	ERR_ACCOUNT_PASSWORD_WRONG: "学号或密码错误",
	ERR_JSON_BIND_WRONG:        "数据绑定失败",
	ERR_CACHE_NOT_EXIST:        "缓存不存在",
}

var (
	SUCCESS                   = 10000
	ERR_INTERNAL_SERVER_WRONG = 99999
	
	ERR_CACHE_EXIST     = 20001
	ERR_CACHE_NOT_EXIST = 20002
	
	ERR_ACCOUNT_PASSWORD_WRONG = 30001
	ERR_USER_EXIST             = 30002
	ERR_TOKEN_NOT_EXIST        = 30003
	ERR_TOKEN_TYPE_WRONG       = 30004
	ERR_TOKEN_WRONG            = 30005
	ERR_TOKEN_EXPIRED          = 30006
	
	ERR_JSON_BIND_WRONG = 40001
)
