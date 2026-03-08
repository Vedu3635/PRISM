package dto

type CreateUserRequest struct {
	FirebaseUID string  `json:"firebase_uid"`
	Email       string  `json:"email"`
	Username    string  `json:"username"`
	FullName    string  `json:"full_name"`
	Password    *string `json:"password"`
	Phone       *string `json:"phone"`
}
