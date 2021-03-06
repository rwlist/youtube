package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rwlist/gjrpc/pkg/gjserver"
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
	"github.com/rwlist/youtube/internal/conf"
	"github.com/rwlist/youtube/internal/global"
	"github.com/rwlist/youtube/internal/isrv"
	"github.com/rwlist/youtube/internal/lists"
	"github.com/rwlist/youtube/internal/logic"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/repos"
	"github.com/rwlist/youtube/internal/rpc"
	"github.com/rwlist/youtube/internal/ytsync"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/youtube/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

//nolint:funlen
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

	b, err := os.ReadFile("client_secret.json")
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
	db = db.Debug()

	err = models.AutoMigrate(db)
	if err != nil {
		log.WithError(err).Fatal("failed to migrate tables")
	}

	usersRepo := repos.NewUsers(db)
	catalogsRepo := repos.NewCatalogLists(db)

	youtubeLikesCursor := ytsync.NewLikesCursor(oauthConfig, usersRepo)

	globalDir := &global.Directory{
		LikedSync:    ytsync.NewLikedSync(youtubeLikesCursor),
		CatalogsRepo: catalogsRepo,
	}

	catalog := lists.NewCatalog(catalogsRepo, db, globalDir)
	hooks := []logic.AuthHook{
		catalog,
	}
	authService := logic.NewAuth(oauthConfig, usersRepo, hooks)
	authServer := isrv.NewAuth(authService)
	youtubeServer := isrv.NewYoutube(oauthConfig)
	catalogServer := isrv.NewListsCatalog(oauthConfig, catalog)
	listServer := isrv.NewList(oauthConfig, catalog)
	handlers := rpc.NewHandlers(authServer, youtubeServer, catalogServer, listServer)

	var rpcHandler jsonrpc.Handler
	rpcHandler = rpc.NewRouter(handlers, rpc.ConvertError).Handle
	rpcHandler = rpc.AuthMiddleware(authService, []string{"auth.oauth"})(rpcHandler)
	rpcHandler = rpc.LogMiddleware()(rpcHandler)

	httpHandler := &gjserver.HandlerHTTP{Handler: rpcHandler} // TODO: add constructor?

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
