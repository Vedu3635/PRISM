package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type firebaseSignInRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type firebaseSignInResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}

type firebaseErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// GetFirebaseToken exchanges email+password for a Firebase ID token
// by calling the Firebase Auth REST API.
func GetFirebaseToken(email, password string) (*firebaseSignInResponse, error) {
	apiKey := os.Getenv("FIREBASE_WEB_API_KEY")
	if apiKey == "" {
		return nil, errors.New("FIREBASE_WEB_API_KEY is not set")
	}

	url := fmt.Sprintf(
		"https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s",
		apiKey,
	)

	body, _ := json.Marshal(firebaseSignInRequest{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to reach Firebase: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var fbErr firebaseErrorResponse
		if err := json.Unmarshal(raw, &fbErr); err == nil && fbErr.Error.Message != "" {
			return nil, errors.New(fbErr.Error.Message)
		}
		return nil, errors.New("firebase authentication failed")
	}

	var result firebaseSignInResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Firebase response: %w", err)
	}

	return &result, nil
}
