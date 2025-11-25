package status_error

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shrewx/ginx"
	"github.com/shrewx/ginx/pkg/i18nx"
	"github.com/shrewx/ginx/pkg/statuserror"
	"net/http"
)

type ServiceError struct {
	// Error 错误信息
	Error string `json:"error" error:"error"`
	// Code 错误码
	Code string `json:"code" error:"code"`
}

func (s *ServiceError) Match(ctx *gin.Context, err error) bool {
	var commonError statuserror.CommonError
	if errors.As(err, &commonError) {
		return true
	}
	return false
}

func (s *ServiceError) Format(ctx *gin.Context, response ginx.ErrorResponse) (statusCode int, contentType string, body []byte, headers http.Header) {
	err := response.Error()
	// 检查错误是否为 ServiceError 类型
	var serviceErr statuserror.CommonError
	statusCode = http.StatusUnprocessableEntity
	contentType = response.ContentType()
	headers = response.Headers()

	if errors.As(err, &serviceErr) {
		if code, ok := s.StatusCodeMap()[serviceErr.Code()]; ok {
			statusCode = code
		}
		i18Error := serviceErr.Localize(i18nx.Instance(), ctx.GetString("lang"))
		s.Error = i18Error.Value()
		s.Code = fmt.Sprintf("%d", serviceErr.Code())
		body, _ = json.Marshal(s)
	}

	return statusCode, contentType, body, headers
}

func (s *ServiceError) StatusCodeMap() map[int64]int {
	return map[int64]int{
		0:                   http.StatusUnprocessableEntity,
		int64(UserNotFound): 450,
	}
}
