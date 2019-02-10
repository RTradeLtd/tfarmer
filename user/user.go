package user

import (
	"github.com/RTradeLtd/database/models"
	"github.com/RTradeLtd/gorm"
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
	if err := f.US.DB.Model(&models.Usage{}).Find(&usages).Error; err != nil {
		return nil, err
	}
	users := []models.User{}
	for _, v := range usages {
		if v.Tier == models.Free {
			user, err := f.UM.FindByUserName(v.UserName)
			if err != nil {
				return nil, err
			}
			users = append(users, *user)
		}
	}
	return users, nil
}

// LightUsers is used to retrieve all light users
func (f *Farmer) LightUsers() ([]models.User, error) {
	usages := []models.Usage{}
	if err := f.US.DB.Model(&models.Usage{}).Find(&usages).Error; err != nil {
		return nil, err
	}
	users := []models.User{}
	for _, v := range usages {
		if v.Tier == models.Light {
			user, err := f.UM.FindByUserName(v.UserName)
			if err != nil {
				return nil, err
			}
			users = append(users, *user)
		}
	}
	return users, nil
}

// PlusUsers is used to retrieve all plus users
func (f *Farmer) PlusUsers() ([]models.User, error) {
	usages := []models.Usage{}
	if err := f.US.DB.Model(&models.Usage{}).Find(&usages).Error; err != nil {
		return nil, err
	}
	users := []models.User{}
	for _, v := range usages {
		if v.Tier == models.Plus {
			user, err := f.UM.FindByUserName(v.UserName)
			if err != nil {
				return nil, err
			}
			users = append(users, *user)
		}
	}
	return users, nil
}

/*
// UsersActive24Hours is used to get all active users in the last 24 hours.
func (f *Farmer) UsersActive24Hours() ([]models.User, error) {
	users := []models.User{}
	// db.Where("updated_at > ?", lastWeek).Find(&users)
	return users, nil
}
*/
