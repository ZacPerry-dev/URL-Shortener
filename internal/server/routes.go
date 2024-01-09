package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.RedirectURL)
	mux.HandleFunc("/addURL", s.PostURL)
	mux.HandleFunc("/deleteURL", s.DeleteURL)
	return mux
}
