package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/repos"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	googleOAuth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type Auth struct {
	oauthConfig *oauth2.Config
	repo        *repos.Users
}

func NewAuth(oauthConfig *oauth2.Config, repo *repos.Users) *Auth {
	return &Auth{
		oauthConfig: oauthConfig,
		repo:        repo,
	}
}

func (a *Auth) generateAccessToken() (string, error) {
	const accessTokenLength = 32

	b := make([]byte, accessTokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (a *Auth) findUserOrCreate(tok *oauth2.Token) (*models.User, error) {
	cli := a.oauthConfig.Client(context.Background(), tok)
	srv, err := googleOAuth.NewService(context.Background(), option.WithHTTPClient(cli))
	if err != nil {
		return nil, err
	}

	userInfo, err := srv.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	googleID := userInfo.Id
	user, err := a.repo.FindByGoogleID(googleID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user != nil {
		user.GoogleToken = tok
		a.UpdateIfRequired(user)
		return user, nil
	}

	newAccessToken, err := a.generateAccessToken()
	if err != nil {
		return nil, err
	}

	user = &models.User{
		Model:       gorm.Model{},
		AccessToken: newAccessToken,
		GoogleEmail: userInfo.Email,
		GoogleID:    userInfo.Id,
		GoogleOAuth: nil, // will be filled later
		GoogleToken: tok,
	}
	_, err = user.SaveTokenIfUpdated()
	if err != nil {
		return nil, err
	}

	err = a.repo.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ExchangeCode for Google OAuth
func (a *Auth) ExchangeCode(code string) (accessToken string, err error) {
	tok, err := a.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}

	user, err := a.findUserOrCreate(tok)
	if err != nil {
		return "", err
	}

	return user.AccessToken, nil
}

// VerifyAuth to find user by accessToken
func (a *Auth) VerifyAuth(accessToken string) (*models.User, error) {
	user, err := a.repo.FindByAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	err = user.PreloadToken()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateIfRequired to update user if token was updated
func (a *Auth) UpdateIfRequired(user *models.User) {
	ok, err := user.SaveTokenIfUpdated()
	if err != nil {
		log.WithError(err).Error("failed to update user token")
		return
	}
	if !ok {
		return
	}

	err = a.repo.UpdateGoogleOAuth(user)
	if err != nil {
		log.WithError(err).Error("failed to update user token")
		return
	}
}

func (a *Auth) CreateRedirectURL() string {
	return a.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
}
