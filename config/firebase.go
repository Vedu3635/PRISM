package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

var FirebaseAuth *auth.Client

func InitFirebase() {

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("Firebase app init error: %v", err)
	}

	FirebaseAuth, err = app.Auth(ctx)
	if err != nil {
		log.Fatalf("Firebase auth init error: %v", err)
	}

	log.Println("Firebase initialized ✅")
}
