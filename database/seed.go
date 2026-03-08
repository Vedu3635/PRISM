package database

import (
	"log"

	"github.com/Vedu3635/PRISM.git/models"
	"github.com/google/uuid"
)

func SeedData() {

	db := DB

	userA := models.User{
		ID:          uuid.New(),
		FirebaseUID: "firebase_user_1",
		Email:       "aryan@test.com",
		Username:    "aryan",
		FullName:    "Aryan Patel",
	}

	userB := models.User{
		ID:          uuid.New(),
		FirebaseUID: "firebase_user_2",
		Email:       "bob@test.com",
		Username:    "bob",
		FullName:    "Bob Shah",
	}

	userC := models.User{
		ID:          uuid.New(),
		FirebaseUID: "firebase_user_3",
		Email:       "charlie@test.com",
		Username:    "charlie",
		FullName:    "Charlie Mehta",
	}

	db.FirstOrCreate(&userA, models.User{Email: userA.Email})
	db.FirstOrCreate(&userB, models.User{Email: userB.Email})
	db.FirstOrCreate(&userC, models.User{Email: userC.Email})

	group := models.Group{
		ID:   uuid.New(),
		Name: "Goa Trip",
		Type: "trip",
	}

	db.FirstOrCreate(&group, models.Group{Name: group.Name})

	db.Create(&models.GroupMember{
		ID:      uuid.New(),
		GroupID: group.ID,
		UserID:  userA.ID,
		Role:    "admin",
	})

	db.Create(&models.GroupMember{
		ID:      uuid.New(),
		GroupID: group.ID,
		UserID:  userB.ID,
		Role:    "member",
	})

	db.Create(&models.GroupMember{
		ID:      uuid.New(),
		GroupID: group.ID,
		UserID:  userC.ID,
		Role:    "member",
	})

	log.Println("Seed data inserted")
}
