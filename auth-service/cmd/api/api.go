package api

import (
	"database/sql"
	"net/http"

	coreUtils "github.com/hoyci/ms-chat/core/utils"

	"github.com/gorilla/mux"
	"github.com/hoyci/ms-chat/auth-service/config"
	"github.com/hoyci/ms-chat/auth-service/service/auth"
	"github.com/hoyci/ms-chat/auth-service/service/healthcheck"
	"github.com/hoyci/ms-chat/auth-service/service/user"
	coreMiddlewares "github.com/hoyci/ms-chat/core/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	addr   string
	db     *sql.DB
	Router *mux.Router
	Config config.Config
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr:   addr,
		db:     db,
		Router: nil,
		Config: config.Envs,
	}
}

func (s *Server) SetupRouter(
	healthCheckHandler *healthcheck.HealthCheckHandler,
	userHandler *user.UserHandler,
	authHandler *auth.AuthHandler,
) *mux.Router {
	coreUtils.InitLogger()
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	subrouter.HandleFunc(
		"/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "docs/swagger.json")
		},
	)

	subrouter.PathPrefix("/swagger/").Handler(
		httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/api/v1/swagger.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		),
	).Methods(http.MethodGet)

	subrouter.HandleFunc("/healthcheck", healthCheckHandler.HandleHealthCheck).Methods(http.MethodGet)

	subrouter.HandleFunc("/auth", authHandler.HandleUserLogin).Methods(http.MethodPost)
	subrouter.HandleFunc("/auth/refresh", authHandler.HandleRefreshToken).Methods(http.MethodPost)

	subrouter.HandleFunc("/users", userHandler.HandleCreateUser).Methods(http.MethodPost)
	subrouter.Handle(
		"/users", coreMiddlewares.AuthMiddleware(http.HandlerFunc(userHandler.HandleGetUserByID), config.Envs.PublicKeyAccess),
	).Methods(http.MethodGet)
	subrouter.Handle(
		"/users",
		coreMiddlewares.AuthMiddleware(http.HandlerFunc(userHandler.HandleUpdateUserByID), config.Envs.PublicKeyAccess),
	).Methods(http.MethodPut)
	subrouter.Handle(
		"/users",
		coreMiddlewares.AuthMiddleware(http.HandlerFunc(userHandler.HandleDeleteUserByID), config.Envs.PublicKeyAccess),
	).Methods(http.MethodDelete)

	s.Router = router

	return router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
