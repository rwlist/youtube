package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rwlist/gjrpc/pkg/gjserver"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
	"github.com/rwlist/youtube/internal/conf"
	"github.com/rwlist/youtube/internal/isrv"
	"github.com/rwlist/youtube/internal/logic"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/repos"
	"github.com/rwlist/youtube/internal/rpc"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/youtube/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

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

	oauthConfig, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope, oauth2.OpenIDScope, oauth2.UserinfoEmailScope)
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
	listsServer := isrv.NewLists(oauthConfig)
	handlers := rpc.NewHandlers(authServer, youtubeServer, listsServer)

	var rpcHandler jsonrpc.Handler = rpc.NewRouter(handlers).Handle
	rpcHandler = rpc.AuthMiddleware(authService, []string{"auth.oauth"})(rpcHandler)
	rpcHandler = rpc.LogMiddleware()(rpcHandler)

	httpHandler := &gjserver.HandlerHTTP{rpcHandler} // TODO: add constructor?

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Handle("/rpc", httpHandler)
	r.Get("/oauth", rpc.OAuthHandler(authService, ""))

	if cfg.ProxyStatic != "" {
		remote, err := url.Parse(cfg.ProxyStatic)
		if err != nil {
			log.WithError(err).Fatal("failed to parse proxy static url")
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			r.Host = remote.Host
			proxy.ServeHTTP(w, r)
		})
	}

	log.WithField("addr", cfg.HttpAddr).Info("starting http server")
	err = http.ListenAndServe(cfg.HttpAddr, r)
	if err != nil {
		log.WithError(err).Fatal("http server error")
	}
}
