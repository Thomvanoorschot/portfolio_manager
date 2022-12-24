package infrastructure

import (
	"github.com/Thomvanoorschot/portfolioManager/app/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Router     *mux.Router
	UnitOfWork *data.UnitOfWork
}

func (server *Server) Run() {
	srv := &http.Server{
		Handler:      server.Router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServeTLS("localhost.crt", "localhost.key"))
}

func NewServer(unitOfWork *data.UnitOfWork,
	router *mux.Router) *Server {
	return &Server{
		Router:     router,
		UnitOfWork: unitOfWork,
	}
}
