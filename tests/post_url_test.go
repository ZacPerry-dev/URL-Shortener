package tests

import (
	"URL-Shortener/internal/database"
	"URL-Shortener/internal/models"
	"URL-Shortener/internal/server"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

// Test for:
// Valid response if the url is created in the DB (201)
// Valid response if the url is found in the DB (302)
// Error if passed an invalid URL (400)
// Error if passed no url (Missing field, 400)

func InsertTestDataDB(t *testing.T, db database.Service) {
	urlCollection := db.GetCollection("url-mappings")

	testData := []interface{}{
		bson.M{"key": "test1", "longUrl": "http://test1.com", "shortUrl": "localhost:8080/abc123"},
		bson.M{"key": "test2", "longUrl": "http://test2.com", "shortUrl": "localhost:8080/def456"},
	}

	_, err := urlCollection.InsertMany(context.TODO(), testData)
	if err != nil {
		panic(err)
	}
}

func CleanDB(t *testing.T, db database.Service) {
	urlCollection := db.GetCollection("url-mappings")

	_, err := urlCollection.DeleteMany(context.TODO(), bson.M{
		"key": bson.M{
			"$in": []string{"test1", "test2"},
		},
	})

	if err != nil {
		panic(err)
	}

	db.CloseConnection()
}

func TestPostURL(t *testing.T) {
	// Setup
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal("Error loading .env file")
	}

	// Print loaded environment variables
	fmt.Println("DB_URI:", os.Getenv("DB_URI"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))

	db := database.New()
	defer CleanDB(t, db)
	InsertTestDataDB(t, db)

	s := server.NewServer(db, "../views")

	// Create a request body
	requestBody := models.BaseUrlInfo{
		Url: "http://www.example.com",
	}

	// Convert the request body to JSON
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// Create a POST request
	req, err := http.NewRequest("POST", "/addURL", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	res := httptest.NewRecorder()

	// Call the PostURL handler
	s.PostURL(res, req)

	// Check the response status code
	if res.Code != http.StatusCreated {
		t.Errorf("Expected status code %d. Got %d", http.StatusCreated, res.Code)
	}

	// You can check other aspects of the response if needed
	// For example, you might want to check the response body for expected content.
	var response models.NewUrlInfo

	fmt.Println(res.Body.String())

	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

}
