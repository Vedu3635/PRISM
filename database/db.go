package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Vedu3635/PRISM.git/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	DB = db

	runMigrations()

	log.Println("Database connected")
}

func runMigrations() {

	err := DB.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.GroupMember{},
		&models.Transaction{},
		&models.TransactionSplit{},
		&models.Settlement{},
		&models.Payment{},
		&models.Balance{},
		&models.AuditLog{},
		&models.Notification{},
		&models.SettlementSplit{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migrated successfully")
}
