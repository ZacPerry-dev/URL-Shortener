package routes

import (
	"URL-Shortener/internal/database"
	"URL-Shortener/internal/handlers"
	"net/http"
)

func AddRoutes(mux *http.ServeMux, db database.IDataBase) {
	h := handlers.NewHandler(db)

	// Handle requests
	mux.Handle("/home", h.HandleHomePageRequest())
	mux.Handle("/", h.HandleGetURL())
	mux.Handle("/shorten", h.HandleShortenURL())
	mux.Handle("/delete", h.HandleDeleteURL())

	// Serve static files (This is just here for now. Idk where else to put this yet...)
	fileServer := http.FileServer(http.Dir("./views/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
}
