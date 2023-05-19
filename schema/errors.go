package schema

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	// Auth
	ErrAuthenticationRequired int = 10000
	ErrTokenInvalid           int = 10001
	ErrTokenExpired           int = 10002
	ErrEmailOrPasswordInvalid int = 10003
	// Requests
	ErrUnparsableJSON       int = 40000
	ErrFieldRequired        int = 40001
	ErrFieldMinChar         int = 40002
	ErrFieldMaxChar         int = 40003
	ErrFieldInvalid         int = 40004
	ErrPasswordUnmatch      int = 40011
	ErrPasswordMinChar      int = 40012
	ErrPasswordMaxChar      int = 40013
	ErrUsernameAlreadyTaken int = 40021
	// Resource general
	ErrResourceNotFound int = 60000
	// Internal
	ErrInternalServer int = 90000
)
