package dto

import "net/http"

type ErrorCode string

const (
    ErrorCodeInvalidRequest     ErrorCode = "INVALID_REQUEST"
    ErrorCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
    ErrorCodeValidationFailed   ErrorCode = "VALIDATION_FAILED"
    ErrorCodeUserExists         ErrorCode = "USER_ALREADY_EXISTS"
    ErrorCodeDatabaseError      ErrorCode = "DATABASE_ERROR"
    ErrorCodeInternalError      ErrorCode = "INTERNAL_ERROR"
    ErrorCodeUserNotFound       ErrorCode = "USER_NOT_FOUND"
    ErrorCodeTokenInvalid       ErrorCode = "TOKEN_INVALID"
    ErrorCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
    ErrorCodeInvalidSigningMethod ErrorCode = "INVALID_SIGNING_METHOD"
)

// @Description Generic error response with a simple message
type ErrorResponse struct {
	Error string `json:"error" example:"Error description"`
}

// @Description Resposta de erro detalhada com código, mensagem e detalhes adicionais
type DetailedErrorResponse struct {
    Code    ErrorCode `json:"code" example:"VALIDATION_FAILED"`
    Message string    `json:"message" example:"Invalid email format"`
    Details string    `json:"details,omitempty" example:"The provided email address is not valid"`
}

// @Description Example for a 403 Invalid Credentials response
type ForbiddenErrorResponse struct {
	Code    string `json:"code" example:"FORBIDDEN"`
	Message string `json:"message" example:"Invalid credentials"`
	Details string `json:"details" example:"Please check your credentials and try again"`
}

// @Description Example for a 409 User Already Exists response
type ConflictErrorResponse struct {
	Code    string `json:"code" example:"CONFLICT"`
	Message string `json:"message" example:"user already exists"`
	Details string `json:"details" example:"Please check your input and try again"`
}

// @Description Example for a 400 Invalid Input response
type BadRequestErrorResponse struct {
	Code    string `json:"code" example:"BAD_REQUEST"`
	Message string `json:"message" example:"invalid email format"`
	Details string `json:"details" example:"Please check your input and try again"`
}

// @Description Example for a 500 Internal Server Error response
type InternalServerErrorResponse struct {
	Code    string `json:"code" example:"INTERNAL_SERVER_ERROR"`
	Message string `json:"message" example:"Internal server error"`
	Details string `json:"details" example:"failed to create user"`
}

func ErrorCodeFromStatusCode(statusCode int) ErrorCode {
    switch statusCode {
    case http.StatusBadRequest:
        return ErrorCodeValidationFailed
    case http.StatusForbidden:
        return ErrorCodeInvalidCredentials
    case http.StatusConflict:
        return ErrorCodeUserExists
    case http.StatusUnauthorized:
        return ErrorCodeTokenInvalid
    case http.StatusInternalServerError:
        return ErrorCodeDatabaseError
    default:
        return ErrorCodeInternalError
    }
}

func ErrorCodeFromError(err error) ErrorCode {
    switch err.Error() {
    case "token expired":
        return ErrorCodeTokenExpired
    case "invalid signing method":
        return ErrorCodeInvalidSigningMethod
    case "invalid token":
        return ErrorCodeTokenInvalid
    default:
        return ErrorCodeInternalError
    }
}