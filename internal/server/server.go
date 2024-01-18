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

type ServerInterface interface {
	RegisterRoutes()
	CreateHttpServer() *http.Server
	GetDB() database.Service
	PostURL(w http.ResponseWriter, r *http.Request)
	GetURL(w http.ResponseWriter, r *http.Request)
	DeleteURL(w http.ResponseWriter, r *http.Request)
	Home(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	port      int
	db        database.Service
	templates *template.Template
	mux       *http.ServeMux
}

func NewServer(db database.Service, viewsPath string) ServerInterface {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	files := []string{
		fmt.Sprintf("%s/index.html", viewsPath),
		fmt.Sprintf("%s/components/url_form.html", viewsPath),
	}

	NewServer := &Server{
		port:      port,
		db:        db,
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

func (s *Server) GetDB() database.Service {
	return s.db
}
