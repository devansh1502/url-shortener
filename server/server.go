package server

import (
	"context"
	"log"
	"net/http"
	"url-shortener/interfaces"
)

type Server struct {
	ctx context.Context
	a   interfaces.API
}

// Start handles the routes and starts the server
func (serv *Server) Start() {
	http.HandleFunc("/redirect/", serv.a.RedirectURL)
	http.HandleFunc("/short/", serv.a.UrlShortner)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// NewServer returns an entry of the Server struct with values.
// this is further consumed by the Start function
func NewServer(ctx context.Context, api interfaces.API) *Server {
	return &Server{ctx: ctx, a: api}
}
