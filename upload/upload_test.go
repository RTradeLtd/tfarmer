package upload

import (
	"fmt"
	"testing"
	"time"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
	"github.com/RTradeLtd/gorm"
	"github.com/RTradeLtd/rtfs"
)

const (
	testCID = "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
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
	if err := db.AutoMigrate(&models.Upload{}).Error; err != nil {
		t.Fatal(err)
	}
}

func TestUpload(t *testing.T) {
	// load configuration
	cfg, err := config.LoadConfig("../testenv/config.json")
	if err != nil {
		t.Fatal(err)
	}
	// open db
	db, err := openDatabaseConnection(cfg)
	if err != nil {
		t.Fatal(err)
	}
	// connect to ipfs
	ipfs, err := rtfs.NewManager(
		cfg.IPFS.APIConnection.Host+":"+cfg.IPFS.APIConnection.Port,
		"", 60*time.Minute,
	)
	if err != nil {
		t.Fatal(err)
	}
	// initialize upload farmer
	farmer := NewFarmer(db, ipfs)
	// fake upload1
	upload1, err := farmer.UM.NewUpload(
		testCID, "pin", models.UploadOptions{
			NetworkName:      "public",
			Username:         "testuser1",
			HoldTimeInMonths: 5,
			Encrypted:        false,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer farmer.UM.DB.Unscoped().Delete(upload1)
	// fake upload2
	upload2, err := farmer.UM.NewUpload(
		testCID, "pin", models.UploadOptions{
			NetworkName:      "public",
			Username:         "testuser2",
			HoldTimeInMonths: 5,
			Encrypted:        false,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer farmer.UM.DB.Unscoped().Delete(upload2)
	if uploadCount, err := farmer.TotalUploadsCount(false); err != nil {
		t.Fatal(err)
	} else if uploadCount < 2 {
		t.Fatal("bad upload count recovered")
	}
	if uploadCount, err := farmer.TotalUploadsCount(true); err != nil {
		t.Fatal(err)
	} else if uploadCount != 1 {
		t.Fatal("bad upload count recovered")
	}
	var (
		//expectedUniqueSize =
		expectedNonUniqueSize = 6.094574928283691e-06
		expectedUniqueSize    = 6.094574928283691e-06
	)
	size, err := farmer.AverageUploadSize(false)
	if err != nil {
		t.Fatal(err)
	}
	if size != expectedNonUniqueSize {
		t.Fatal("failed to calculate correct non unique average size")
	}
	size, err = farmer.AverageUploadSize(true)
	if err != nil {
		t.Fatal(err)
	}
	if size != expectedUniqueSize {
		t.Fatal("Failed to calculate correct unique average size")
	}
}

func openDatabaseConnection(cfg *config.TemporalConfig) (*gorm.DB, error) {
	dbConnURL := fmt.Sprintf("host=127.0.0.1 port=%s user=postgres dbname=temporal password=%s sslmode=disable",
		cfg.Database.Port, cfg.Database.Password)

	return gorm.Open("postgres", dbConnURL)
}
