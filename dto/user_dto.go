package dto

type CreateUserRequest struct {
	FirebaseUID string  `json:"firebase_uid"`
	Email       string  `json:"email"      binding:"required,email"`
	Username    string  `json:"username"   binding:"required"`
	FullName    string  `json:"full_name"  binding:"required"`
	Password    *string `json:"password"`
	Phone       *string `json:"phone"`
}

type UpdateUserRequest struct {
	Username     *string `json:"username"`
	FullName     *string `json:"full_name"`
	Phone        *string `json:"phone"`
	CurrencyPref *string `json:"currency_pref"`
	Password     *string `json:"password"`
}
