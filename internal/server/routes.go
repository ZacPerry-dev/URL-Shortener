package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() {

	s.mux.HandleFunc("/home", s.Home)
	s.mux.HandleFunc("/", s.GetURL)
	s.mux.HandleFunc("/addURL", s.PostURL)
	s.mux.HandleFunc("/deleteURL", s.DeleteURL)

	fileServer := http.FileServer(http.Dir("./views/static/"))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
}
