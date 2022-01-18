package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rwlist/youtube/pkg/conf"
	"github.com/rwlist/youtube/pkg/logic"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	cfg, err := conf.ParseEnv()
	if err != nil {
		log.WithError(err).Fatal("failed to parse config from env")
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(cfg.PrometheusBind, mux)
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("prometheus server error")
		}
	}()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	oauthConfig, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	l := logic.NewLogic(oauthConfig)

	if cfg.PreloadedToken != "" {
		tok := &oauth2.Token{}
		err := json.Unmarshal([]byte(cfg.PreloadedToken), tok)
		if err != nil {
			log.WithError(err).Error("failed to parse preloaded token")
		}
		l.SetToken(tok)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		str := l.Status()
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(str))
	})

	r.Get("/oauth", func(w http.ResponseWriter, r *http.Request) {
		redirect, err := l.DoOauth(r.URL)
		if err != nil {
			log.WithError(err).Error("failed to do oauth")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
	})

	r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
		str := l.ListPlaylists()
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(str))
	})

	r.Get("/liked", func(w http.ResponseWriter, r *http.Request) {
		str := l.ListLikedPlaylist()
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(str))
	})

	log.WithField("addr", cfg.HttpAddr).Info("starting http server")
	err = http.ListenAndServe(cfg.HttpAddr, r)
	if err != nil {
		log.WithError(err).Fatal("http server error")
	}
}
