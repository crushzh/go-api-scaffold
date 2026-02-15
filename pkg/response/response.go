package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the unified API response structure
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData holds paginated data
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// PageQuery is the pagination request
type PageQuery struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Keyword  string `form:"keyword" json:"keyword,omitempty"`
}

// Normalize normalizes pagination parameters
func (p *PageQuery) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// Offset returns the database offset
func (p *PageQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// ========================
// Error code definitions
// ========================
// Format: XYYY (X=category, YYY=sequence)
// 0 success | 1 params | 2 resource | 3 business | 4 auth | 5 system

const (
	CodeSuccess = 0

	// 1xxx Client errors
	CodeParamError = 1001
	CodeParamType  = 1002
	CodeJSONError  = 1003

	// 2xxx Resource errors
	CodeNotFound = 2001
	CodeConflict = 2002

	// 3xxx Business errors (extensible per domain)
	CodeBizError = 3001

	// 4xxx Authentication errors
	CodeUnauthorized = 4001
	CodeForbidden    = 4002
	CodeTokenExpired = 4003

	// 5xxx System errors
	CodeInternal = 5001
	CodeDatabase = 5002
	CodeTimeout  = 5005
)

// ========================
// Success responses
// ========================

// OK returns a success response without data
func OK(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
	})
}

// Success returns a success response with data
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage returns a success response with custom message
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// SuccessPage returns a paginated success response
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// ========================
// Error responses
// ========================

// Fail returns a business failure (HTTP 200, non-zero business code)
func Fail(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeInternal,
		Message: message,
	})
}

// Error returns a response with specified error code
func Error(c *gin.Context, code int, message string) {
	c.JSON(getHTTPStatus(code), Response{
		Code:    code,
		Message: message,
	})
}

// ParamError returns a parameter error response
func ParamError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeParamError,
		Message: message,
	})
}

// NotFound returns a resource not found response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
	})
}

// Conflict returns a resource conflict response
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Code:    CodeConflict,
		Message: message,
	})
}

// Unauthorized returns an unauthenticated response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeUnauthorized,
		Message: message,
	})
}

// Forbidden returns a forbidden response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: message,
	})
}

// ServerError returns an internal server error response
func ServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeInternal,
		Message: message,
	})
}

// DatabaseError returns a database error response
func DatabaseError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeDatabase,
		Message: "database error",
	})
}

// getHTTPStatus maps error codes to HTTP status codes
func getHTTPStatus(code int) int {
	switch {
	case code == CodeSuccess:
		return http.StatusOK
	case code >= 1000 && code < 2000:
		return http.StatusBadRequest
	case code == CodeNotFound:
		return http.StatusNotFound
	case code == CodeConflict:
		return http.StatusConflict
	case code >= 2000 && code < 3000:
		return http.StatusBadRequest
	case code >= 3000 && code < 4000:
		return http.StatusOK // Business errors use 200
	case code == CodeUnauthorized || code == CodeTokenExpired:
		return http.StatusUnauthorized
	case code == CodeForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
