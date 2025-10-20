package status_error

import "net/http"

//go:generate toolx gen error -p error_codes -c StatusError
//go:generate toolx gen errorYaml -p error_codes -o ../i18n -c StatusError
type StatusError int

const (
	// @errZH 请求参数错误
	// @errEN bad request
	BadRequest StatusError = http.StatusBadRequest*1e8 + iota + 1
	// @errZH 用户ID无效，必须为正整数
	// @errEN invalid user ID, must be a positive integer
	InvalidUserID
)

const (
	// @errZH 未授权，请先授权
	// @errEN unauthorized
	Unauthorized StatusError = http.StatusUnauthorized*1e8 + iota + 1
)

const (
	// @errZH 禁止操作
	// @errEN forbidden
	Forbidden StatusError = http.StatusForbidden*1e8 + iota + 1
)

const (
	// @errZH 资源未找到
	// @errEN not found
	NotFound StatusError = http.StatusNotFound*1e8 + iota + 1
	// @errZH 用户不存在，用户ID：{{.UserID}}
	// @errEN user not found, user ID: {{.UserID}}
	UserNotFound
)

const (
	// @errZH 资源冲突
	// @errEN conflict
	Conflict StatusError = http.StatusConflict*1e8 + iota + 1
)

const (
	// @errZH 操作存在异常，请稍后重试
	// @errEN operation failed, please retry later
	StatusBadGateway StatusError = http.StatusBadGateway*1e8 + iota + 1
	// @errZH 数据操作存在异常，请稍后重试
	// @errEN data operation failed, please retry later
	DataOperationFailed
)
