package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hoyci/auth-service/service/healthcheck"
	"github.com/hoyci/auth-service/types"
	httpSwagger "github.com/swaggo/http-swagger"
)

type APIServer struct {
	addr string
	// db     *sql.DB
	Router *mux.Router
	Config types.Config
}

func NewApiServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
		// db:     db,
		Router: nil,
		Config: types.Config{},
	}
}

func (s *APIServer) SetupRouter(
	healthCheckHandler *healthcheck.HealthCheckHandler,
) *mux.Router {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	subrouter.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})

	subrouter.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/api/v1/swagger.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	subrouter.HandleFunc("/healthcheck", healthCheckHandler.HandleHealthCheck).Methods(http.MethodGet)

	s.Router = router

	return router
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
