package dto

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

type SignupResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

type SigninResponse struct {
	Token string `json:"token"`
}
