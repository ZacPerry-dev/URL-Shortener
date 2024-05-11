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
	db     *database.Database
}

func NewServer(dbURI, dbName string) (*Server, error) {
	db, err := database.NewDatabase(dbURI, dbName)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	server := &Server{router: mux, db: db}
	routes.AddRoutes(server.router, server.db)

	return server, nil
}

func RunServer(
	ctx context.Context,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// setup db ????
	dbURI := getenv("DB_URI")
	dbName := getenv("DB_NAME")
	port := getenv("PORT")

	// Start the server and listen for incoming requests / errors
	server, err := NewServer(dbURI, dbName)
	if err != nil {
		return err
	}
	defer server.db.CloseConnection()

	errs := make(chan error, 1)

	// listen and server in a goroutine
	go func() {
		fmt.Print("Starting server on port: ", port, "...\n")
		errs <- http.ListenAndServe(":"+port, server.router)
	}()

	// block until either an error or a signal is received
	select {
	case err := <-errs:
		fmt.Fprintf(stderr, "error Case: %v\n", err)
		return err
	case <-ctx.Done():
		return nil
	}
}
