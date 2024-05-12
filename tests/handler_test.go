package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"URL-Shortener/internal/server"

	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Helper function to drop in some test data into the local DB instance
func AddTestData(ctx context.Context, dbURI, dbName, port string) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(dbName).Collection("url-mappings")
	_, err = collection.InsertOne(ctx, bson.M{
		"longurl":  "https://www.example.com",
		"shorturl": "https://localhost:" + port + "/56f945",
		"key":      "56f945",
	})
	if err != nil {
		return err
	}

	return nil
}

// Helper function to clean up the test data from the local DB instance
func CleanTestData(ctx context.Context, dbURI, dbName string) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	collection := client.Database(dbName).Collection("url-mappings")
	_, err = collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}

/*
Tests the HandleGetURL function. Ensures that it is able to get and redirect to the correct URL.
*/
func TestHandleGetURL(t *testing.T) {
	is := is.New(t)
	getenv := func(key string) string {
		switch key {
		case "PORT":
			return "3000"
		case "DB_URI":
			return "mongodb://localhost:27017"
		case "DB_NAME":
			return "url-shortener-tests"
		case "DB_HOST":
			return "localhost"
		case "IDLE_TIMEOUT":
			return "1800"
		default:
			return ""
		}
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	defer cancel()

	go func() {
		server.RunServer(ctx, getenv, os.Stdin, os.Stdout, os.Stderr)
	}()

	// Add the test data to the local DB instance
	err := AddTestData(ctx, getenv("DB_URI"), getenv("DB_NAME"), getenv("PORT"))
	is.NoErr(err)

	// Create the request
	requestURL := "http://localhost:" + getenv("PORT") + "/56f945"
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	// Ensure it found the URL and redirected successfully
	is.Equal(res.StatusCode, http.StatusOK)
	is.NoErr(err)

	// Clean up the test data in the local DB instance
	err = CleanTestData(ctx, getenv("DB_URI"), getenv("DB_NAME"))
	is.NoErr(err)
}

/*
Tests the HandleShortenURL function. Ensures that it is able to shorten a URL and return the shortened URL.
*/
func TestHandleShortenURL(t *testing.T) {
	is := is.New(t)
	getenv := func(key string) string {
		switch key {
		case "PORT":
			return "3001"
		case "DB_URI":
			return "mongodb://localhost:27017"
		case "DB_NAME":
			return "url-shortener-tests"
		case "DB_HOST":
			return "localhost"
		case "IDLE_TIMEOUT":
			return "1800"
		default:
			return ""
		}
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	defer cancel()

	go func() {
		server.RunServer(ctx, getenv, os.Stdin, os.Stdout, os.Stderr)
	}()

	// Create the request
	requestURL := "http://localhost:" + getenv("PORT") + "/shorten"
	requestBody := `{"url": "https://www.example.com"}`
	res, err := http.Post(requestURL, "application/json", strings.NewReader(requestBody))
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	// Ensure it successfully shortened the URL
	is.Equal(res.StatusCode, http.StatusCreated)
	is.NoErr(err)

	// Clean up the test data in the local DB instance
	err = CleanTestData(ctx, getenv("DB_URI"), getenv("DB_NAME"))
	is.NoErr(err)
}

/*
Tests the HandleDeleteURL function. Ensures that it is able to delete a URL from the DB.
*/
func TestHandleDeleteURL(t *testing.T) {
	is := is.New(t)
	getenv := func(key string) string {
		switch key {
		case "PORT":
			return "3002"
		case "DB_URI":
			return "mongodb://localhost:27017"
		case "DB_NAME":
			return "url-shortener-tests"
		case "DB_HOST":
			return "localhost"
		case "IDLE_TIMEOUT":
			return "1800"
		default:
			return ""
		}
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	defer cancel()

	go func() {
		server.RunServer(ctx, getenv, os.Stdin, os.Stdout, os.Stderr)
	}()

	// Clean and add data into the db
	err := CleanTestData(ctx, getenv("DB_URI"), getenv("DB_NAME"))
	is.NoErr(err)
	err = AddTestData(ctx, getenv("DB_URI"), getenv("DB_NAME"), getenv("PORT"))
	is.NoErr(err)

	// Construct the request
	requestURL := "http://localhost:" + getenv("PORT") + "/delete"
	requestBody := struct {
		Url string `json:"url"`
	}{
		Url: "https://localhost:3002/56f945",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("DELETE", requestURL, strings.NewReader(string(jsonBody)))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// Ensure the URL was found and deleted successfully
	is.Equal(res.StatusCode, http.StatusOK)
	is.NoErr(err)
}
