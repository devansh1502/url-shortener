package main

import (
	"context"
	"url-shortener/api"
	"url-shortener/database"
	"url-shortener/server"
)

func main() {
	ctx := context.Background()
	sI := database.NewStore()
	a := api.NewAPI(ctx, sI)
	serv := server.NewServer(ctx, a)
	serv.Start()
}