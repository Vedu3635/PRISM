package dto

type SignupRequest struct {
	FirebaseUID string  `json:"firebase_uid" binding:"required"`
	Email       string  `json:"email"        binding:"required,email"`
	Password    string  `json:"password"     binding:"required,min=6"`
	Username    string  `json:"username"     binding:"required"`
	FullName    string  `json:"full_name"    binding:"required"`
	Phone       *string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
