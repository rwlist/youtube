package repos

import (
	"github.com/rwlist/youtube/internal/models"
	"gorm.io/gorm"
)

type Users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *Users {
	return &Users{db}
}

func (r *Users) FindByAccessToken(accessToken string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("access_token = ?", accessToken).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Users) Save(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *Users) UpdateGoogleOAuth(user *models.User) error {
	return r.db.Model(user).Update("google_o_auth", user.GoogleOAuth).Error
}

func (r *Users) FindByGoogleID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("google_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
