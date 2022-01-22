package isrv

import (
	"context"
	"github.com/rwlist/youtube/internal/logic"
	"github.com/rwlist/youtube/internal/proto"
	"github.com/rwlist/youtube/internal/rpc"
)

type Auth struct {
	authService *logic.Auth
}

func NewAuth(authService *logic.Auth) *Auth {
	return &Auth{
		authService: authService,
	}
}

func (a *Auth) Oauth(ctx context.Context) (proto.OAuthResponse, error) {
	redirectURL := a.authService.CreateRedirectURL()
	return proto.OAuthResponse{
		RedirectURL: redirectURL,
	}, nil
}

func (a *Auth) Status(ctx context.Context) (proto.AuthStatus, error) {
	user := rpc.GetUser(ctx)
	return proto.AuthStatus{
		Email: user.GoogleEmail,
	}, nil
}
