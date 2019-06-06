package user

import (
	"fmt"
	"testing"

	"github.com/RTradeLtd/config/v2"
	"github.com/RTradeLtd/database/v2/models"
	"github.com/jinzhu/gorm"
)

func TestMigration(t *testing.T) {
	cfg, err := config.LoadConfig("../testenv/config.json")
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Usage{}).Error; err != nil {
		t.Fatal(err)
	}
}

func Test_User(t *testing.T) {
	// open config
	cfg, err := config.LoadConfig("../testenv/config.json")
	if err != nil {
		t.Fatal(err)
	}
	// open db connection
	db, err := openDatabaseConnection(cfg)
	if err != nil {
		t.Fatal(err)
	}
	// initialize farmer
	farmer := NewFarmer(db)

	// create test user1
	user1, err := farmer.UM.NewUserAccount(
		"testuser1",
		"password123",
		"testuser1@example.org",
	)
	if err != nil {
		t.Fatal(err)
	}
	usage1, err := farmer.US.FindByUserName("testuser1")
	if err != nil {
		t.Fatal(err)
	}
	// defer the user model delete
	defer farmer.UM.DB.Unscoped().Delete(user1)
	// defer the usage model delete
	defer farmer.US.DB.Unscoped().Delete(usage1)

	// create test user2 and set to light tier
	user2, err := farmer.UM.NewUserAccount(
		"testuser2",
		"password123",
		"testuser2@example.org",
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := farmer.US.UpdateTier(user2.UserName, models.Paid); err != nil {
		t.Fatal(err)
	}
	usage2, err := farmer.US.FindByUserName("testuser2")
	if err != nil {
		t.Fatal(err)
	}
	// defer the user model delete
	defer farmer.UM.DB.Unscoped().Delete(user2)
	// defer the usage model delete
	defer farmer.US.DB.Unscoped().Delete(usage2)

	// create test user3 and set to plus tier
	user3, err := farmer.UM.NewUserAccount(
		"testuser3",
		"password123",
		"testuser3@example.org",
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := farmer.US.UpdateTier(user3.UserName, models.Paid); err != nil {
		t.Fatal(err)
	}
	usage3, err := farmer.US.FindByUserName("testuser3")
	if err != nil {
		t.Fatal(err)
	}
	// defer the user model delete
	defer farmer.UM.DB.Unscoped().Delete(user3)
	// defer the usage model delete
	defer farmer.US.DB.Unscoped().Delete(usage3)

	// find registered users
	users, err := farmer.RegisteredUsers()
	if err != nil {
		t.Fatal(err)
	}

	// search for test users
	var (
		foundTestUser1 bool
		foundTestUser2 bool
		foundTestUser3 bool
	)
	for _, v := range users {
		if v.UserName == "testuser1" {
			foundTestUser1 = true
		} else if v.UserName == "testuser2" {
			foundTestUser2 = true
		} else if v.UserName == "testuser3" {
			foundTestUser3 = true
		}
	}
	if !foundTestUser1 || !foundTestUser2 || !foundTestUser3 {
		t.Fatal("failed to find correct users")
	}

	// reset found variables
	foundTestUser1, foundTestUser2, foundTestUser3 = false, false, false

	// find free users
	users, err = farmer.FreeUsers()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range users {
		if v.UserName == "testuser1" {
			usage, err := farmer.US.FindByUserName(v.UserName)
			if err != nil {
				t.Fatal(err)
			}
			if usage.UserName == "testuser1" && usage.Tier == models.Free {
				foundTestUser1 = true
			}
		}
	}
	if !foundTestUser1 {
		t.Fatal("failed to find testuser1 from free tier search")
	}

	// find paid users
	users, err = farmer.PaidUsers()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range users {
		if v.UserName == "testuser2" {
			usage, err := farmer.US.FindByUserName(v.UserName)
			if err != nil {
				t.Fatal(err)
			}
			if usage.UserName == "testuser2" && usage.Tier == models.Paid {
				foundTestUser2 = true
			}
		} else if v.UserName == "testuser3" {
			usage, err := farmer.US.FindByUserName(v.UserName)
			if err != nil {
				t.Fatal(err)
			}
			if usage.UserName == "testuser3" && usage.Tier == models.Paid {
				foundTestUser3 = true
			}
		}
	}
	if !foundTestUser2 || !foundTestUser3 {
		t.Fatal("failed to find testuser2from paid tier search")
	}

	// reset found variables
	foundTestUser1, foundTestUser2, foundTestUser3 = false, false, false

	// get active users in the last 24 hours
	users, err = farmer.ActiveUsers24Hours()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range users {
		if v.UserName == "testuser1" {
			foundTestUser1 = true
		} else if v.UserName == "testuser2" {
			foundTestUser2 = true
		} else if v.UserName == "testuser3" {
			foundTestUser3 = true
		}
	}
	if !foundTestUser1 || !foundTestUser2 || !foundTestUser3 {
		t.Fatal("failed to find correct users")
	}

	// reset found variables
	foundTestUser1, foundTestUser2, foundTestUser3 = false, false, false

	// get active usage in the last 24 hours
	users, err = farmer.ActiveUsage24Hours()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range users {
		if v.UserName == "testuser1" {
			foundTestUser1 = true
		} else if v.UserName == "testuser2" {
			foundTestUser2 = true
		} else if v.UserName == "testuser3" {
			foundTestUser3 = true
		}
	}
	if !foundTestUser1 || !foundTestUser2 || !foundTestUser3 {
		t.Fatal("failed to find correct users")
	}
}

func openDatabaseConnection(cfg *config.TemporalConfig) (*gorm.DB, error) {
	dbConnURL := fmt.Sprintf("host=127.0.0.1 port=%s user=postgres dbname=temporal password=%s sslmode=disable",
		cfg.Database.Port, cfg.Database.Password)

	return gorm.Open("postgres", dbConnURL)
}
