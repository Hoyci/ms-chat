package api

import (
	"net/http"

	"github.com/gorilla/mux"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
	"github.com/hoyci/ms-chat/message-service/config"
	"github.com/hoyci/ms-chat/message-service/service/healthcheck"
	httpSwagger "github.com/swaggo/http-swagger"
)

type APIServer struct {
	addr   string
	Router *mux.Router
	Config config.Config
}

func NewApiServer(addr string) *APIServer {
	return &APIServer{
		addr:   addr,
		Router: nil,
		Config: config.Config{},
	}
}

func (s *APIServer) SetupRouter(
	healthCheckHandler *healthcheck.HealthCheckHandler,
	// roomHandler *room.RoomHandler,
) *mux.Router {
	coreUtils.InitLogger()
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

	// subrouter.HandleFunc("/rooms", roomHandler.HandleCreateRoom).Methods(http.MethodPost)
	// subrouter.HandleFunc("/rooms", roomHandler.HandleGetRoomByID).Methods(http.MethodGet)

	s.Router = router

	return router
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
