package schema

type RegisterRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"min=6,required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"min=6,required"`
}

type ResendVerificationEmailRequest struct {
	Email string `json:"email" validate:"email,required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"email,required"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"min=6,required"`
}
