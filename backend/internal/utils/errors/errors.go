package errors

import "net/http"

type CustomError interface {
    error
    StatusCode() int
    Message() string
}

type BusinessError struct {
    Msg  string
    Code int
}

func (e *BusinessError) Error() string {
    return e.Msg
}

func (e *BusinessError) StatusCode() int {
    return e.Code
}

func (e *BusinessError) Message() string {
    return e.Msg
}

func NewValidationError(msg string) *BusinessError {
    return &BusinessError{Msg: msg, Code: http.StatusBadRequest}
}

func GetStatusCode(err error) int {
    if bizErr, ok := err.(*BusinessError); ok {
        return bizErr.Code
    }
    return http.StatusInternalServerError
}

var (
    // Auth errors
    ErrUserAlreadyExists  = &BusinessError{Msg: "user already exists", Code: http.StatusConflict}
    ErrInvalidEmail       = &BusinessError{Msg: "invalid email format", Code: http.StatusBadRequest}
    ErrWeakPassword       = &BusinessError{Msg: "password must be at least 8 characters", Code: http.StatusBadRequest}
    ErrInvalidCredentials = &BusinessError{Msg: "invalid credentials", Code: http.StatusForbidden}
    ErrUserNotFound       = &BusinessError{Msg: "user not found", Code: http.StatusForbidden}
    ErrFailedToCheckUser  = &BusinessError{Msg: "failed to check user existence", Code: http.StatusInternalServerError}
    ErrFailedToSetPassword = &BusinessError{Msg: "failed to set password", Code: http.StatusInternalServerError}
    ErrFailedToCreateUser  = &BusinessError{Msg: "failed to create user", Code: http.StatusInternalServerError}
    
    // JWT errors
    ErrTokenGeneration    = &BusinessError{Msg: "failed to generate token", Code: http.StatusInternalServerError}
    ErrInvalidToken       = &BusinessError{Msg: "invalid token", Code: http.StatusUnauthorized}
    ErrTokenExpired       = &BusinessError{Msg: "token expired", Code: http.StatusUnauthorized}
    ErrInvalidSigningMethod = &BusinessError{Msg: "invalid signing method", Code: http.StatusUnauthorized}

    // Device errors
    ErrDeviceNotFound      = &BusinessError{Msg: "device not found", Code: http.StatusNotFound}
    ErrDeviceAlreadyExists = &BusinessError{Msg: "device with this serial number already exists", Code: http.StatusConflict}
    ErrForbidden           = &BusinessError{Msg: "access to this resource is forbidden", Code: http.StatusForbidden}
    ErrDatabaseError       = &BusinessError{Msg: "database error", Code: http.StatusInternalServerError}
)