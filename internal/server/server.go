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
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	files := []string{
		"./views/index.html",
		"./views/components/url_form.html",
	}

	NewServer := &Server{
		port:      port,
		db:        database.New(),
		templates: utils.ParseTemplates(files...),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
