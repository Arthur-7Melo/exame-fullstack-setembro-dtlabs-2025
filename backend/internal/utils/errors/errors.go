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
)