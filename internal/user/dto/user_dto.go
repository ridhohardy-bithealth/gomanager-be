package dto

type AuthRequest struct {
	Password string `json:"password" validate:"required,min=8,max=32"`
	Email    string `json:"email" validate:"required,email,min=1,max=255"`
	Action   string `json:"action" validate:"required,oneof=create login"`
}

type AuthResponse struct {
	Email       string `json:"email"`
	AccessToken string `json:"token"`
}

type UserResponse struct {
	Email           string `json:"email"`
	Username        string `json:"name"`
	UserImageUri    string `json:"userImageUri"`
	CompanyName     string `json:"companyName"`
	CompanyImageUri string `json:"companyImageUri"`
}

type UpdateUserRequest struct {
	Email           *string `json:"email" validate:"required,email,min=1,max=255"`
	Username        *string `json:"name" validate:"required,min=4,max=52"`
	UserImageUri    *string `json:"userImageUri" validate:"required"`
	CompanyName     *string `json:"companyName" validate:"required,min=4,max=52"`
	CompanyImageUri *string `json:"companyImageUri" validate:"required"`
}
