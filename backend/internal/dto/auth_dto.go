package dto

// @Description Login credentials for user authentication
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securePassword123"`
}

// @Description User registration information
type SignupRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securePassword123"`
}

// @Description JWT token returned upon successful authentication
type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}
