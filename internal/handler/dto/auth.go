package dto

// request POST /auth/register
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
}

// response POST /auth/register
type RegisterResponse struct {
	ID string `json:"id" example:"f4e28058-5f0b-437f-965f-f4670c8e22a7"`
}

// request POST /auth/login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

// response POST /auth/login
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR..."`
}
