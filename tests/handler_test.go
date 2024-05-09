package tests

import (
	"URL-Shortener/internal/database"
	"URL-Shortener/internal/server"
	"net/http"
	"testing"

	"github.com/matryer/is"
)

func TestGetURLHandler(t *testing.T) {
	is := is.New(t)

	db := database.New()
	defer db.CloseConnection()

	s := server.NewServer(db, "./views")
	httpServe := s.CreateHttpServer()

	err := httpServe.ListenAndServe()
	defer httpServe.Close()
	if err != nil {
		panic(err)
	}

	_, errr := http.NewRequest("GET", "http://localhost:8000/", nil)
	if errr != nil {
		t.Error("Expected response from server, got error")
	}

	is.NoErr(errr)
}
