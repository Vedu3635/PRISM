package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

func InitFirebase() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Firebase app init error: %v", err)
	}

	FirebaseAuth, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Firebase auth init error: %v", err)
	}

	log.Println("Firebase initialized ✅")
}
