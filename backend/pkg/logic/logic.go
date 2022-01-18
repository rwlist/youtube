package logic

import (
	"context"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"net/url"
)

type Logic struct {
	oauthConfig    *oauth2.Config
	token          *oauth2.Token
	client         *http.Client
	youtubeService *youtube.Service
}

func NewLogic(oauthConfig *oauth2.Config) *Logic {
	return &Logic{
		oauthConfig: oauthConfig,
	}
}

func (l *Logic) Status() string {
	res := "ok<br>"

	if l.token == nil {
		res += "no token, not logged in<br>"
		googleRedirectURL := l.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
		res += "<a href=\"" + googleRedirectURL + "\">log in</a>"
	} else {
		res += "token ok, logged in<br>"
		res += "do something, API call, whatever<br>"
		tokenJSON, err := json.Marshal(l.token)
		if err != nil {
			res += "error marshalling token: " + err.Error()
		} else {
			res += "token: " + string(tokenJSON) + "<br>"
		}
		res += "<br><br>"
		res += "<a href=\"/list\">list playlists</a><br>"
		res += "<a href=\"/liked\">get liked videos</a><br>"
	}

	return res
}

func (l *Logic) ListPlaylists() string {
	res := "ok<br>"

	call := l.youtubeService.Playlists.List([]string{"contentDetails", "id", "localizations", "player", "snippet", "status"})
	call.Id("LL", "WL")
	//call.Mine(true)
	resp, err := call.Do()
	if err != nil {
		res += "error: " + err.Error()
	} else {
		res += "response:<br><br>" + "<pre>" + spew.Sdump(resp) + "</pre>"
	}

	for _, playlist := range resp.Items {
		res += "<br>" + playlist.Snippet.Title
	}

	return res
}

func (l *Logic) ListLikedPlaylist() string {
	res := "ok<br>"

	call := l.youtubeService.PlaylistItems.List([]string{"contentDetails", "id", "snippet", "status"})
	call.PlaylistId("LL")
	call.MaxResults(50)
	resp, err := call.Do()
	if err != nil {
		res += "error: " + err.Error()
	} else {
		res += "response:<br><br>" + "<pre>" + spew.Sdump(resp) + "</pre>"
	}

	for _, playlistItem := range resp.Items {
		res += "<br>" + playlistItem.Snippet.Title
	}

	return res
}

func (l *Logic) DoOauth(url *url.URL) (redirectURL string, err error) {
	code := url.Query().Get("code")
	tok, err := l.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}

	return "/", l.SetToken(tok)
}

func (l *Logic) SetToken(tok *oauth2.Token) error {
	var err error
	l.token = tok
	l.client = l.oauthConfig.Client(context.Background(), tok)
	l.youtubeService, err = youtube.NewService(context.Background(), option.WithHTTPClient(l.client))
	return err
}
