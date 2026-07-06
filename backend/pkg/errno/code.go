package errno

// 业务错误码定义
// 格式：模块前缀 + 错误编号
// 100xx - 通用错误
// 200xx - 用户模块

var (
	// 通用
	OK          = &Errno{Code: 0, Message: "成功"}
	ErrParam    = &Errno{Code: 10001, Message: "参数错误"}
	ErrInternal = &Errno{Code: 10002, Message: "内部错误"}
	ErrDB       = &Errno{Code: 10003, Message: "数据库错误"}
	ErrNotFound = &Errno{Code: 10004, Message: "资源不存在"}

	// 认证模块
	ErrUnauthorized = &Errno{Code: 10005, Message: "未登录或Token已过期"}
	ErrTokenInvalid = &Errno{Code: 10006, Message: "Token无效"}

	// 用户模块
	ErrUserNotFound = &Errno{Code: 20001, Message: "用户不存在"}
	ErrUserExists   = &Errno{Code: 20002, Message: "用户已存在"}
)

type Errno struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Errno) Error() string {
	return e.Message
}

func (e *Errno) WithMessage(msg string) *Errno {
	return &Errno{Code: e.Code, Message: msg}
}

func Decode(err error) (int, string) {
	if err == nil {
		return 0, "成功"
	}
	if e, ok := err.(*Errno); ok {
		return e.Code, e.Message
	}
	return 10002, err.Error()
}
