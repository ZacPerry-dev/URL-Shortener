package main

import (
	"URL-Shortener/internal/database"
	"URL-Shortener/internal/server"
)

func main() {
	db := database.New()
	defer db.CloseConnection()

	s := server.NewServer(db, "./views")
	httpServe := s.CreateHttpServer()

	err := httpServe.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
