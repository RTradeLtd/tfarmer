package mail_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/RTradeLtd/config/v2"
	"github.com/RTradeLtd/tfarmer/mail"
	"github.com/jinzhu/gorm"
)

var (
	testRecipientEmail1 = "postables+test1@rtradetechnologies.com"
	testRecipientEmail2 = "postables+test2@rtradetechnologies.com"
	testRecipientName1  = "postables1"
	testRecipientName2  = "postables2"
	testCfgPath         = "../testenv/config.json"
)

func TestMail(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Sendgrid.APIKey == "" {
		cfg.Sendgrid.APIKey = os.Getenv("SENDGRID_API_KEY")
		cfg.Sendgrid.EmailAddress = "temporal@rtradetechnologies.com"
		cfg.Sendgrid.EmailName = "Temporal TravisCI Test"
	}
	mm, err := mail.NewManager(cfg, db)
	if err != nil {
		t.Fatal(err)
	}
	content := fmt.Sprint("<br>WowSuchEmail<br>WowSuchFormat")
	if _, err = mm.SendEmail(
		"testEmail",
		content,
		"",
		testRecipientName1,
		testRecipientEmail1,
	); err != nil {
		t.Fatal(err)
	}
	if _, err = mm.SendEmail(
		"testEmail",
		content,
		"text/html",
		testRecipientName1,
		testRecipientEmail1,
	); err != nil {
		t.Fatal(err)
	}
	if err = mm.BulkSend(
		"testEmail",
		content,
		"text/html",
		[]string{testRecipientName1, testRecipientName2},
		[]string{testRecipientEmail1, testRecipientEmail2},
	); err != nil {
		t.Fatal(err)
	}
}

func openDatabaseConnection(cfg *config.TemporalConfig) (*gorm.DB, error) {
	dbConnURL := fmt.Sprintf("host=127.0.0.1 port=%s user=postgres dbname=temporal password=%s sslmode=disable",
		cfg.Database.Port, cfg.Database.Password)

	return gorm.Open("postgres", dbConnURL)
}
