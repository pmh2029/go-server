package dtos

type BaseResponse struct {
	Code    int            `json:"code"`
	Data    interface{}    `json:"data"`
	Message interface{}    `json:"message"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

// ErrorResponse struct
type ErrorResponse struct {
	ErrorDetails interface{} `json:"error_details"`
}
