package dtos

type BaseResponse struct {
	Code  int            `json:"code"`
	Data  interface{}    `json:"data"`
	Error *ErrorResponse `json:"error,omitempty"`
}

// ErrorResponse struct
type ErrorResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	ErrorDetails string `json:"error_details"`
}
