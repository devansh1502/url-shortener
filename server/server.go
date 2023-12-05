package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"url-shortener/interfaces"
)

type Server struct {
	ctx context.Context
	a   interfaces.API
}

// Start handles the routes and starts the server
func (serv *Server) Start() {
	ctx, stop := signal.NotifyContext(serv.ctx, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		http.HandleFunc("/redirect/", serv.a.RedirectURL)
		http.HandleFunc("/short/", serv.a.UrlShortner)
		http.HandleFunc("/metrics/", serv.a.Metrics)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-ctx.Done()
	log.Printf("Server Shutting Down!")

	// Gracefull shutdown
	stop()
	return
}

// NewServer returns an entry of the Server struct with values.
// this is further consumed by the Start function
func NewServer(ctx context.Context, api interfaces.API) *Server {
	return &Server{ctx: ctx, a: api}
}
