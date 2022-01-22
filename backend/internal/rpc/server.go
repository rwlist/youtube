package rpc

import (
	"net/http"
)

const cookieAccessToken = "access_token"

func AccessTokenFromRequest(r *http.Request) string {
	cookie, err := r.Cookie(cookieAccessToken)
	if err != nil {
		return ""
	}
	return cookie.Value
}

type oauthService interface {
	ExchangeCode(code string) (accessToken string, err error)
}

func OAuthHandler(oauthService oauthService, redirect string) http.HandlerFunc {
	if redirect == "" {
		redirect = "/"
	}

	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		accessToken, err := oauthService.ExchangeCode(code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  cookieAccessToken,
			Value: accessToken,
		})
		http.Redirect(w, r, redirect, http.StatusFound)
	}
}
