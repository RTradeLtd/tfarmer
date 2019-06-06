package user

import (
	"time"

	"github.com/RTradeLtd/database/v2/models"
	"github.com/jinzhu/gorm"
)

// used to scrape user related data

// Farmer is the user farmer for Temporal
type Farmer struct {
	UM *models.UserManager
	US *models.UsageManager
}

// NewFarmer instantiates our user farmer class
func NewFarmer(db *gorm.DB) *Farmer {
	return &Farmer{
		UM: models.NewUserManager(db),
		US: models.NewUsageManager(db),
	}
}

// RegisteredUsers is used to retrieve all registered users
func (f *Farmer) RegisteredUsers() ([]models.User, error) {
	users := []models.User{}
	if err := f.UM.DB.Model(&models.User{}).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// FreeUsers is used to retrieve all free users
func (f *Farmer) FreeUsers() ([]models.User, error) {
	usages := []models.Usage{}
	if err := f.US.DB.Where("tier = ?", models.Free).Find(&usages).Error; err != nil {
		return nil, err
	}
	users := []models.User{}
	for _, v := range usages {
		user, err := f.UM.FindByUserName(v.UserName)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

// PaidUsers is used to retrieve all paid users
func (f *Farmer) PaidUsers() ([]models.User, error) {
	usages := []models.Usage{}
	if err := f.US.DB.Where("tier = ?", models.Paid).Find(&usages).Error; err != nil {
		return nil, err
	}
	users := []models.User{}
	for _, v := range usages {
		user, err := f.UM.FindByUserName(v.UserName)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

// ActiveUsers24Hours is used to get all active users in the last 24 hours.
// this doesn't consider actual usage of the platform
//  and is simply a login-in based metrics
func (f *Farmer) ActiveUsers24Hours() ([]models.User, error) {
	users := []models.User{}
	tt := time.Now().Add(time.Hour * -24)
	if err := f.UM.DB.Where("updated_at > ?", tt).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// ActiveUsage24Hours is used to get all active users in the last 24 hours
// based off actual usage of the platform, and considered things like:
// pubsub, key management, ipns, and upload usage
func (f *Farmer) ActiveUsage24Hours() ([]models.User, error) {
	usages := []models.Usage{}
	tt := time.Now().Add(time.Hour * -24)
	if err := f.US.DB.Where("updated_at > ?", tt).Find(&usages).Error; err != nil {
		return nil, err
	}
	users := []models.User{}
	for _, v := range usages {
		user, err := f.UM.FindByUserName(v.UserName)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}
