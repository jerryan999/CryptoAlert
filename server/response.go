// wrap json reponse
package server

const (
	CodeInvalidDataError  = 40002
	CodeParseRequestError = 40201
	CodeDataNotFoundError = 40401
	CodeInternalErr       = 50001
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Code:    0,
		Data:    data,
		Message: "success",
	}
}

func NewErrorResponse(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

func NewInvalidDataErrorResponse() *Response {
	return NewErrorResponse(CodeInvalidDataError, "Invalid data")
}

func NewParseRequestErrorResponse() *Response {
	return NewErrorResponse(CodeParseRequestError, "Parse request error")
}

func NewDataNotFoundErrorResponse() *Response {
	return NewErrorResponse(CodeDataNotFoundError, "Data not found")
}

func NewInternalErrorResponse() *Response {
	return NewErrorResponse(CodeInternalErr, "Internal error")
}
