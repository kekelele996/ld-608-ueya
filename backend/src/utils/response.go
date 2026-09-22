package utils

import "github.com/gin-gonic/gin"

// APIError is the single error envelope produced by service/controller
// wrappers and rendered by errorHandlerMiddleware.
type APIError struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	HTTPCode  int         `json:"-"`
}

func (e *APIError) Error() string { return e.Message }

// NewAPIError builds an envelope with explicit HTTP status.
func NewAPIError(httpCode int, code, message string, details interface{}) *APIError {
	return &APIError{HTTPCode: httpCode, Code: code, Message: message, Details: details}
}

// Pagination parses page/page_size query params shared by every list route.
func Pagination(c *gin.Context) (int, int) {
	page, size := 1, 50
	if v := c.Query("page"); v != "" {
		if n := atoi(v); n > 0 {
			page = n
		}
	}
	if v := c.Query("page_size"); v != "" {
		if n := atoi(v); n > 0 && n <= 200 {
			size = n
		}
	}
	return page, size
}

func atoi(value string) int {
	n := 0
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return 0
		}
		n = n*10 + int(ch-'0')
	}
	return n
}
