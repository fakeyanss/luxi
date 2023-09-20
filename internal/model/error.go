package model

import "github.com/gin-gonic/gin"

type UsgError struct {
	Code     string
	Message  string
	HttpCode int
}

func (e *UsgError) Error() string {
	return e.Message
}

func (e *UsgError) ToUsgErrorBody(c *gin.Context) *UsgErrorBody {
	body := &UsgErrorBody{}
	body.Code = e.Code
	body.Message = e.Message
	body.RequestId = c.Request.Header.Get(UsgReqIdHeaderKey)
	body.Resource = c.Request.URL.Path
	return body
}

func NewUsgError(code, msg string, httpCode int) *UsgError {
	return &UsgError{
		Code:     code,
		Message:  msg,
		HttpCode: httpCode,
	}
}

func NewInvalidArgument(msg string) *UsgError {
	return NewUsgError("InvalidArgument", msg, 400)
}
