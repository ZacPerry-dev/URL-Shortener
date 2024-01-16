package main

import (
	"URL-Shortener/internal/database"
	"URL-Shortener/internal/server"
)

func main() {
	db := database.New()
	s := server.NewServer(db, "./views")

	httpServe := s.CreateHttpServer()

	err := httpServe.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
