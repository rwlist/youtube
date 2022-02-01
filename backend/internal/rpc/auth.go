package rpc

import (
	"context"
	"github.com/rwlist/gjrpc/pkg/gjserver"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/proto"
	"net/http"
	"time"
)

type ctxKey string

const ctxUser ctxKey = "user"

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, ctxUser, user)
}

func GetUser(ctx context.Context) *models.User {
	if user, ok := ctx.Value(ctxUser).(*models.User); ok {
		return user
	}
	return nil
}

type authProvider interface {
	VerifyAuth(accessToken string) (*models.User, error)
	UpdateIfRequired(user *models.User)
}

func AuthMiddleware(auth authProvider, exceptions []string) jsonrpc.Middleware {
	exceptionMap := make(map[string]struct{})
	for _, e := range exceptions {
		exceptionMap[e] = struct{}{}
	}

	return func(next jsonrpc.Handler) jsonrpc.Handler {
		return func(req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
			if _, ok := exceptionMap[req.Method]; ok {
				return next(req)
			}

			ctx := req.Context
			httpReq := gjserver.HttpRequest(ctx)
			httpResp := gjserver.HttpResponse(ctx)
			accessToken := AccessTokenFromRequest(httpReq)
			user, err := auth.VerifyAuth(accessToken)
			if err != nil {
				return nil, proto.AuthError.WithData(err.Error())
			}
			defer auth.UpdateIfRequired(user)

			ctx = WithUser(ctx, user)
			req.Context = ctx

			res, rpcErr := next(req)
			if rpcErr != nil && rpcErr.Code == proto.AuthError.Code {
				http.SetCookie(httpResp, &http.Cookie{
					Name:    cookieAccessToken,
					Value:   "",
					Expires: time.Unix(0, 0),
				})
			}

			return res, rpcErr
		}
	}
}
