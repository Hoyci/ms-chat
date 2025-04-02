package api

import (
	"database/sql"
	"github.com/hoyci/ms-chat/contacts-service/config"
	"github.com/hoyci/ms-chat/contacts-service/services/contacts"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/hoyci/ms-chat/contacts-service/config"
	//"github.com/hoyci/ms-chat/contacts-service/service/auth"
	//"github.com/hoyci/ms-chat/contacts-service/service/healthcheck"
	//"github.com/hoyci/ms-chat/contacts-service/service/user"
	//"github.com/hoyci/ms-chat/contacts-service/utils"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
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

func (s *Server) SetupRouter(contactHandler *contacts.ContactHandler) *mux.Router {
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

	subrouter.HandleFunc("/contacts", contactHandler.HandleCreateContact).Methods(http.MethodPost)

	s.Router = router

	return router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
