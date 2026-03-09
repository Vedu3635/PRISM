package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Vedu3635/PRISM.git/models"
)

func Run(db *gorm.DB) error {

	// Skip if data already exists
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		fmt.Println("Seed data already exists, skipping...")
		return nil
	}

	fmt.Println("Seeding database...")

	// ─── Users ────────────────────────────────────────────────────────────────

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	passwordHash := string(hash)

	alice := models.User{
		ID:           uuid.New(),
		FirebaseUID:  "firebase-alice-001",
		Email:        "alice@example.com",
		Username:     "alice",
		FullName:     "Alice Sharma",
		PasswordHash: passwordHash,
		Phone:        strPtr("+919876543210"),
		CurrencyPref: "INR",
		IsVerified:   true,
		IsDeleted:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	bob := models.User{
		ID:           uuid.New(),
		FirebaseUID:  "firebase-bob-002",
		Email:        "Jack@example.com",
		Username:     "Jack",
		FullName:     "Jack Mehta",
		PasswordHash: passwordHash,
		Phone:        strPtr("+919123456789"),
		CurrencyPref: "INR",
		IsVerified:   true,
		IsDeleted:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	charlie := models.User{
		ID:           uuid.New(),
		FirebaseUID:  "firebase-charlie-003",
		Email:        "chase@example.com",
		Username:     "chase",
		FullName:     "chases Patel",
		PasswordHash: passwordHash,
		Phone:        strPtr("+919000000001"),
		CurrencyPref: "INR",
		IsVerified:   true,
		IsDeleted:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	users := []models.User{alice, bob, charlie}
	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("seed users: %w", err)
	}
	fmt.Printf("  ✓ Created %d users\n", len(users))

	// ─── Groups ───────────────────────────────────────────────────────────────

	goaTrip := models.Group{
		ID:          uuid.New(),
		CreatedBy:   alice.ID,
		Name:        "Dubai Trip 2025",
		Description: strPtr("Beach trip with friends"),
		Type:        "trip",
		Currency:    "INR",
		InviteCode:  uuid.New().String()[:8],
		IsActive:    true,
		IsPersonal:  false,
		CreatedAt:   time.Now(),
	}

	flatmates := models.Group{
		ID:          uuid.New(),
		CreatedBy:   bob.ID,
		Name:        "Flatmates",
		Description: strPtr("Monthly shared expenses"),
		Type:        "home",
		Currency:    "INR",
		InviteCode:  uuid.New().String()[:8],
		IsActive:    true,
		IsPersonal:  false,
		CreatedAt:   time.Now(),
	}

	groups := []models.Group{goaTrip, flatmates}
	if err := db.Create(&groups).Error; err != nil {
		return fmt.Errorf("seed groups: %w", err)
	}
	fmt.Printf("  ✓ Created %d groups\n", len(groups))

	// ─── Group Members ────────────────────────────────────────────────────────

	members := []models.GroupMember{
		// Goa Trip — Alice is admin, Bob and Charlie are members
		{ID: uuid.New(), GroupID: goaTrip.ID, UserID: alice.ID, Role: "admin", JoinedAt: time.Now()},
		{ID: uuid.New(), GroupID: goaTrip.ID, UserID: bob.ID, Role: "member", JoinedAt: time.Now()},
		{ID: uuid.New(), GroupID: goaTrip.ID, UserID: charlie.ID, Role: "member", JoinedAt: time.Now()},

		// Flatmates — Bob is admin, Alice is member
		{ID: uuid.New(), GroupID: flatmates.ID, UserID: bob.ID, Role: "admin", JoinedAt: time.Now()},
		{ID: uuid.New(), GroupID: flatmates.ID, UserID: alice.ID, Role: "member", JoinedAt: time.Now()},
	}

	if err := db.Create(&members).Error; err != nil {
		return fmt.Errorf("seed members: %w", err)
	}
	fmt.Printf("  ✓ Created %d group members\n", len(members))

	// ─── Transactions + Splits + Balances ─────────────────────────────────────

	// Transaction 1: Alice paid for hotel (Goa Trip), split equally among 3
	hotelAmount := 6000.00
	hotelShare := hotelAmount / 3

	hotel := models.Transaction{
		ID:           uuid.New(),
		GroupID:      goaTrip.ID,
		PaidBy:       alice.ID,
		Title:        "Hotel Booking",
		Amount:       hotelAmount,
		Currency:     "INR",
		Category:     strPtr("accommodation"),
		SplitType:    "equal",
		Notes:        strPtr("Beachside resort, 2 nights"),
		Status:       "active",
		TransactedAt: time.Now().Add(-48 * time.Hour),
		CreatedAt:    time.Now().Add(-48 * time.Hour),
	}

	// Transaction 2: Bob paid for dinner (Goa Trip), split equally among 3
	dinnerAmount := 1500.00
	dinnerShare := dinnerAmount / 3

	dinner := models.Transaction{
		ID:           uuid.New(),
		GroupID:      goaTrip.ID,
		PaidBy:       bob.ID,
		Title:        "Beach Dinner",
		Amount:       dinnerAmount,
		Currency:     "INR",
		Category:     strPtr("food"),
		SplitType:    "equal",
		Notes:        strPtr("Seafood restaurant"),
		Status:       "active",
		TransactedAt: time.Now().Add(-24 * time.Hour),
		CreatedAt:    time.Now().Add(-24 * time.Hour),
	}

	// Transaction 3: Bob paid rent (Flatmates), split equally between 2
	rentAmount := 20000.00
	rentShare := rentAmount / 2

	rent := models.Transaction{
		ID:           uuid.New(),
		GroupID:      flatmates.ID,
		PaidBy:       bob.ID,
		Title:        "March Rent",
		Amount:       rentAmount,
		Currency:     "INR",
		Category:     strPtr("rent"),
		SplitType:    "equal",
		Notes:        strPtr("Monthly rent"),
		Status:       "active",
		TransactedAt: time.Now().Add(-72 * time.Hour),
		CreatedAt:    time.Now().Add(-72 * time.Hour),
	}

	transactions := []models.Transaction{hotel, dinner, rent}
	if err := db.Create(&transactions).Error; err != nil {
		return fmt.Errorf("seed transactions: %w", err)
	}
	fmt.Printf("  ✓ Created %d transactions\n", len(transactions))

	// ─── Transaction Splits ───────────────────────────────────────────────────

	splits := []models.TransactionSplit{
		// Hotel splits
		{ID: uuid.New(), TransactionID: hotel.ID, UserID: alice.ID, OwedAmount: hotelShare},
		{ID: uuid.New(), TransactionID: hotel.ID, UserID: bob.ID, OwedAmount: hotelShare},
		{ID: uuid.New(), TransactionID: hotel.ID, UserID: charlie.ID, OwedAmount: hotelShare},

		// Dinner splits
		{ID: uuid.New(), TransactionID: dinner.ID, UserID: alice.ID, OwedAmount: dinnerShare},
		{ID: uuid.New(), TransactionID: dinner.ID, UserID: bob.ID, OwedAmount: dinnerShare},
		{ID: uuid.New(), TransactionID: dinner.ID, UserID: charlie.ID, OwedAmount: dinnerShare},

		// Rent splits
		{ID: uuid.New(), TransactionID: rent.ID, UserID: bob.ID, OwedAmount: rentShare},
		{ID: uuid.New(), TransactionID: rent.ID, UserID: alice.ID, OwedAmount: rentShare},
	}

	if err := db.Create(&splits).Error; err != nil {
		return fmt.Errorf("seed splits: %w", err)
	}
	fmt.Printf("  ✓ Created %d transaction splits\n", len(splits))

	// ─── Balances ─────────────────────────────────────────────────────────────

	// From hotel: Bob owes Alice 2000, Charlie owes Alice 2000
	// From dinner: Alice owes Bob 500, Charlie owes Bob 500
	// From rent: Alice owes Bob 10000
	// Net: Bob owes Alice 2000 - 500 = 1500 → Alice owes Bob 500 - 2000 = net Bob owes Alice 1500
	// Alice owes Bob 10000 (rent) + 500 (dinner) - 2000 (hotel) = 8500 net

	balances := []models.Balance{
		// Goa Trip balances
		{ID: uuid.New(), GroupID: goaTrip.ID, FromUserID: bob.ID, ToUserID: alice.ID, NetAmount: hotelShare - dinnerShare}, // Bob owes Alice net 1500
		{ID: uuid.New(), GroupID: goaTrip.ID, FromUserID: charlie.ID, ToUserID: alice.ID, NetAmount: hotelShare},           // Charlie owes Alice 2000
		{ID: uuid.New(), GroupID: goaTrip.ID, FromUserID: charlie.ID, ToUserID: bob.ID, NetAmount: dinnerShare},            // Charlie owes Bob 500

		// Flatmates balances
		{ID: uuid.New(), GroupID: flatmates.ID, FromUserID: alice.ID, ToUserID: bob.ID, NetAmount: rentShare}, // Alice owes Bob 10000
	}

	if err := db.Create(&balances).Error; err != nil {
		return fmt.Errorf("seed balances: %w", err)
	}
	fmt.Printf("  ✓ Created %d balances\n", len(balances))

	fmt.Println("Seeding complete!")
	fmt.Println()
	fmt.Println("  Seeded UUIDs (use these in Postman):")
	fmt.Printf("  alice.ID    = %s\n", alice.ID)
	fmt.Printf("  bob.ID      = %s\n", bob.ID)
	fmt.Printf("  charlie.ID  = %s\n", charlie.ID)
	fmt.Printf("  goaTrip.ID  = %s\n", goaTrip.ID)
	fmt.Printf("  flatmates.ID = %s\n", flatmates.ID)

	return nil
}

func strPtr(s string) *string {
	return &s
}
