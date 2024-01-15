package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"URL-Shortener/internal/database"
	"URL-Shortener/utils"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port      int
	db        database.Service
	templates *template.Template
	mux       *http.ServeMux
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	files := []string{
		"./views/index.html",
		"./views/components/url_form.html",
	}

	NewServer := &Server{
		port:      port,
		db:        database.New(),
		templates: utils.ParseTemplates(files...),
		mux:       http.NewServeMux(),
	}

	NewServer.RegisterRoutes()

	return NewServer
}

func (s *Server) CreateHttpServer() *http.Server {
	addr := fmt.Sprintf(":%d", s.port)

	server := &http.Server{
		Addr:         addr,
		Handler:      s.mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
