package models

const (
	// Error codes for frontend i18n support
	ErrCodeUnauthorized       = 10001 // Unauthorized (general)
	ErrCodeBadRequest         = 10002 // Bad Request
	ErrCodeForbidden          = 10003 // Forbidden
	ErrCodeNotFound           = 10004 // Not Found
	ErrCodeInternalError      = 10005 // Internal Error
	ErrCodeValidation         = 10006 // Validation Error
	ErrCodeConflict           = 10007 // Conflict
	ErrCodeMissingAuthHeader  = 10008 // Missing Auth Header
	ErrCodeInvalidAuthHeader  = 10009 // Invalid Auth Header
	ErrCodeInvalidToken       = 10010 // Invalid Token
	ErrCodeInvalidTokenClaims = 10011 // Invalid Token Claims
	ErrCodeMissingUserID      = 10012 // Missing User ID
)

// APIResponse represents a standardized API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an error in API response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// PaginatedData represents paginated data structure
type PaginatedData struct {
	Page         int         `json:"page"`
	PageSize     int         `json:"pageSize"`
	TotalElement int64       `json:"totalElement"`
	Data         interface{} `json:"data"`
}

// SuccessResponse creates a success response
func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
}

// PaginatedResponse creates a paginated response
func PaginatedResponse(page, pageSize int, totalElement int64, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data: PaginatedData{
			Page:         page,
			PageSize:     pageSize,
			TotalElement: totalElement,
			Data:         data,
		},
	}
}
