package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"momentia-be/internal/config"
	"momentia-be/model"
)

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	return string(hash)
}

func main() {
	cfg, _ := config.Load()
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{})
	if err != nil {
		panic("database connection failed: " + err.Error())
	}

	// Insert a user first to satisfy foreign key constraints
	users := []model.User{
		{
		ID:           		1,
		Username:     		"worker1",
		Email:				"worker1@gustr.com",
		Msisdn:				"+6281284091230",
		PasswordHash: 		hashPassword("testpassword"),
		},
		{
			ID: 			2,
			Username: 		"worker2",
			Email: 			"worker2@gustr.com",
			Msisdn: 		"+6281284091232",
			PasswordHash: 	hashPassword("testpassword"),
		},
		{
			ID: 			3,
			Username: 		"worker3",
			Email: 			"worker3@gustr.com",
			Msisdn: 		"+6281284091233",
			PasswordHash: 	hashPassword("testpassword"),
		},
		{
			ID: 			4,
			Username: 		"worker4",
			Email: 			"worker4@gustr.com",
			Msisdn: 		"+6281284091234",
			PasswordHash: 	hashPassword("testpassword"),
		},
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users)

	// Reset the sequence so future auto-increment inserts don't collide with seeded IDs
	db.Exec("SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT MAX(id) FROM users))")

	// Insert persons first and capture their real IDs
	person := []model.Person{
		{CreatorUserID: 1, Name: "Alice", Relationship: "Partner", IsPinned: true},
		{CreatorUserID: 1, Name: "Bob", Relationship: "Family", IsPinned: true},
		{CreatorUserID: 1, Name: "Charlie", Relationship: "Friend", IsPinned: false},
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&person)

	// alice := persons[0].ID.String()
	// bob := persons[1].ID.String()
	// charlie := persons[2].ID.String()

	// db.Clauses(clause.OnConflict{DoNothing: true}).Create(&[]model.PersonLikes{
	// 	{PersonID: alice, LikedFood: "Pizza", LikedPlaces: "Beach", LikedColor: "Blue", AdditionalNotes: "Loves sunsets"},
	// 	{PersonID: bob, LikedFood: "Sushi", LikedPlaces: "Mountains", LikedColor: "Green", AdditionalNotes: "Enjoys hiking"},
	// 	{PersonID: charlie, LikedFood: "Burgers", LikedPlaces: "City", LikedColor: "Red", AdditionalNotes: "Fan of nightlife"},
	// })

	// db.Clauses(clause.OnConflict{DoNothing: true}).Create(&[]model.PersonMedicalInfo{
	// 	{PersonID: alice, Name: "Alice", Relationship: "Partner", PhoneNumber: "123-456-7890", BloodType: "A+", Allergies: "Peanuts", AdditionalNotes: "Has asthma"},
	// 	{PersonID: bob, Name: "Bob", Relationship: "Family", PhoneNumber: "987-654-3210", BloodType: "B-", Allergies: "None", AdditionalNotes: "Diabetic"},
	// 	{PersonID: charlie, Name: "Charlie", Relationship: "Friend", PhoneNumber: "555-555-5555", BloodType: "O+", Allergies: "Shellfish", AdditionalNotes: "Vegetarian"},
	// })

	// db.Clauses(clause.OnConflict{DoNothing: true}).Create(&[]model.PersonContacts{
	// 	{PersonID: alice, Email: "alice@example.com", PhoneNumber: "123-456-7890", WhatsApp: "alice_whatsapp", AdditionalNotes: "Best contact method is email"},
	// 	{PersonID: bob, Email: "bob@example.com", PhoneNumber: "987-654-3210", WhatsApp: "bob_whatsapp", AdditionalNotes: "Prefers phone calls"},
	// 	{PersonID: charlie, Email: "chuck@example.com", PhoneNumber: "555-555-5555", WhatsApp: "charlie_whatsapp", AdditionalNotes: "Available after 6 PM"},
	// })

	// db.Clauses(clause.OnConflict{DoNothing: true}).Create(&[]model.PersonDislikes{
	// 	{PersonID: alice, DislikedFoods: "Broccoli", DislikedPlaces: "Crowded places", DislikedColor: "Yellow", AdditionalNotes: "Dislikes loud noises"},
	// 	{PersonID: bob, DislikedFoods: "Spicy food", DislikedPlaces: "Cold places", DislikedColor: "Purple", AdditionalNotes: "Not a fan of surprises"},
	// 	{PersonID: charlie, DislikedFoods: "Fish", DislikedPlaces: "Rural areas", DislikedColor: "Gray", AdditionalNotes: "Prefers structured plans"},
	// })
}
