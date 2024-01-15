package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/home", s.Home)
	mux.HandleFunc("/", s.GetURL)
	mux.HandleFunc("/addURL", s.PostURL)
	mux.HandleFunc("/deleteURL", s.DeleteURL)

	fileServer := http.FileServer(http.Dir("./views/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return mux
}
