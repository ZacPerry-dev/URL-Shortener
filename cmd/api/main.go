package main

import "URL-Shortener/internal/server"

func main() {
	s := server.NewServer()

	httpServe := s.CreateHttpServer()

	err := httpServe.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
