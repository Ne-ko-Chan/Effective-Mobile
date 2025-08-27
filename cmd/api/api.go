package api

import (
	"database/sql"
	"log"
	"net/http"
	"rest-service/service/subscription"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db: db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/testservice/v1").Subrouter()
	
	subscriptionStore := subscription.NewStore(s.db)
	subscriptionHandler := subscription.NewHandler(subscriptionStore)
	subscriptionHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
