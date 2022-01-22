package models

import (
	"bytes"
	"encoding/json"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	AccessToken string `gorm:"unique_index"`

	GoogleEmail string
	GoogleID    string        `gorm:"unique_index"`
	GoogleOAuth []byte        `gorm:"type:jsonb"`
	GoogleToken *oauth2.Token `gorm:"-"`
}

func (u *User) PreloadToken() error {
	u.GoogleToken = &oauth2.Token{}
	err := json.Unmarshal(u.GoogleOAuth, u.GoogleToken)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) SaveTokenIfUpdated() (bool, error) {
	newJSON, err := json.Marshal(u.GoogleToken)
	if err != nil {
		return false, err
	}
	if !bytes.Equal(u.GoogleOAuth, newJSON) {
		u.GoogleOAuth = newJSON
		return true, nil
	}
	return false, nil
}
