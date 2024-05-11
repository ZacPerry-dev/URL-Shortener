package main

import (
	"URL-Shortener/internal/server"
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()

	getenv := func(key string) string {
		switch key {
		case "PORT":
			return os.Getenv("PORT")
		case "DB_URI":
			return os.Getenv("DB_URI")
		case "DB_NAME":
			return os.Getenv("DB_NAME")
		case "DB_HOST":
			return os.Getenv("DB_HOST")
		case "IDLE_TIMEOUT":
			return os.Getenv("IDLE_TIMEOUT")
		default:
			return ""
		}
	}

	if err := server.RunServer(ctx, getenv, os.Stdin, os.Stdout, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
