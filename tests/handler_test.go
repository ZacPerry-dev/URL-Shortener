package tests

import (
	"URL-Shortener/internal/handlers"
	"URL-Shortener/internal/models"
	"URL-Shortener/tests/mocks"
	"URL-Shortener/utils"
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"
)

// / Test cases for the HandleGetURL function
func TestHandleGetURLFound(t *testing.T) {
	mockDB := mocks.NewMockDB()
	handlers := handlers.NewHandler(mockDB)

	req := httptest.NewRequest("GET", "/123456", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetURL().ServeHTTP(rr, req)

	if rr.Code != 302 {
		t.Errorf("Expected status code 302, got %d", rr.Code)
	}
}

func TestHandleGetURLNotFound(t *testing.T) {
	mockDB := mocks.NewMockDB()
	handlers := handlers.NewHandler(mockDB)

	req := httptest.NewRequest("GET", "/123457", nil)
	rr := httptest.NewRecorder()

	handlers.HandleGetURL().ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Errorf("Expected status code 404, got %d", rr.Code)
	}
}

// Test cases for the HandleShortenURL function
func TestHandleShortenURLFound(t *testing.T) {
	mockDB := mocks.NewMockDB()
	handlers := handlers.NewHandler(mockDB)

	// create json body for the request using BaseURLInfo
	baseUrl := models.BaseUrlInfo{
		Url: "https://www.google.com",
	}

	jsonBody, err := utils.JsonMarshal(baseUrl)
	if err != nil {
		t.Errorf("Error marshaling JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handlers.HandleShortenURL().ServeHTTP(rr, req)

	if rr.Code != 302 {
		t.Errorf("Expected status code 302, got %d", rr.Code)
		fmt.Println(rr.Body.String())
		//print the db collection
		fmt.Println(mockDB.Collection)
	}
}
