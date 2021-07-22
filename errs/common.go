package errs

// 公用错误码

var (
	ErrOK                  = New(0, "OK")
	ErrNotModified         = New(-304, "Not Modified")          // 木有改动
	ErrBadRequest          = New(-400, "Bad Request")           // 客户端请求错误，如参数错误
	ErrUnauthorized        = New(-401, "Unauthorized")          // 未认证
	ErrForbidden           = New(-403, "Forbidden")             // 访问权限不足
	ErrNotFound            = New(-404, "Not Found")             // 啥都木有
	ErrMethodNotAllowed    = New(-405, "Method Not Allowed")    // 不支持该方法
	ErrConflict            = New(-409, "Conflict")              // 冲突
	ErrInternalServerError = New(-500, "Internal Server Error") // 服务器内部错误
	ErrServiceUnavailable  = New(-503, "Service Unavailable")   // 服务暂不可用
	ErrGatewayTimeout      = New(-504, "Gateway Timeout")       // 服务调用超时
	ErrUnknown             = New(-999, "Unknown Error")         // 未知错误
)
