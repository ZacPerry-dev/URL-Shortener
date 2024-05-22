// server/server.go
package server

import (
	"URL-Shortener/internal/database"
	"URL-Shortener/internal/routes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
)

type Server struct {
	router *http.ServeMux
	db     database.IDataBase
}

func NewServer(db database.IDataBase) *Server {
	mux := http.NewServeMux()
	server := &Server{router: mux, db: db}
	routes.AddRoutes(server.router, server.db)
	return server
}

func RunServer(
	ctx context.Context,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	dbURI := getenv("DB_URI")
	dbName := getenv("DB_NAME")
	port := getenv("PORT")

	db, err := database.NewDatabase(dbURI, dbName)
	if err != nil {
		return err
	}
	defer db.CloseConnection()

	server := NewServer(db)

	errs := make(chan error, 1)

	go func() {
		fmt.Print("Starting server on port: ", port, "...\n")
		errs <- http.ListenAndServe(":"+port, server.router)
	}()

	select {
	case err := <-errs:
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return err
	case <-ctx.Done():
		return nil
	}
}
