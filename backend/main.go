package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rwlist/gjrpc/pkg/gjserver"
	"github.com/rwlist/youtube/internal/conf"
	"github.com/rwlist/youtube/internal/isrv"
	"github.com/rwlist/youtube/internal/logic"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/repos"
	"github.com/rwlist/youtube/internal/rpc"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		log.WithError(err).Fatal("unable to read client secret file")
	}

	oauthConfig, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("failed to connect to postgres")
	}
	usersRepo := repos.NewUsers(db)

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.WithError(err).Fatal("failed to migrate tables")
	}

	authService := logic.NewAuth(oauthConfig, usersRepo)
	authServer := isrv.NewAuth(authService)
	youtubeServer := isrv.NewYoutube(oauthConfig)
	handlers := rpc.NewHandlers(authServer, youtubeServer)
	rpcHandler := rpc.NewRouter(handlers)
	httpHandler := &gjserver.HandlerHTTP{
		Handler: rpcHandler.Handle,
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Handle("/rpc", httpHandler)
	r.Get("/oauth", rpc.OAuthHandler(authService, cfg.CustomOAuthRedirect))

	log.WithField("addr", cfg.HttpAddr).Info("starting http server")
	err = http.ListenAndServe(cfg.HttpAddr, r)
	if err != nil {
		log.WithError(err).Fatal("http server error")
	}
}
