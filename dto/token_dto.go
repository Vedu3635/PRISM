package dto

// TokenResponse is returned by POST /auth/token.
// Use the IDToken as: Authorization: Bearer <idToken>
type TokenResponse struct {
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	Email        string `json:"email"`
}
