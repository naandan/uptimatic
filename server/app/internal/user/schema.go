package user

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email,required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"min=6,required"`
	NewPassword string `json:"new_password" validate:"min=6,required"`
}

type UpdateFotoRequest struct {
	FileName string `json:"file_name" validate:"required"`
}

type PresignedUrlRequest struct {
	FileName    string `json:"file_name" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
}
