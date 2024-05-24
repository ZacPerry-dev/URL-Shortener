package tests

import (
	"URL-Shortener/internal/handlers"
	"URL-Shortener/internal/models"
	"URL-Shortener/tests/mocks"
	"URL-Shortener/utils"
	"bytes"
	"net/http/httptest"
	"testing"
)

// / Test cases for the HandleGetURL handler
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

// Test cases for the HandleShortenURL handler
func TestHandleShortenURLFound(t *testing.T) {
	mockDB := mocks.NewMockDB()
	handlers := handlers.NewHandler(mockDB)

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
	}
}

func TestHandleShortenURLCreate(t *testing.T) {
	mockDB := mocks.NewMockDB()
	handlers := handlers.NewHandler(mockDB)

	baseUrl := models.BaseUrlInfo{
		Url: "https://www.example.com",
	}

	jsonBody, err := utils.JsonMarshal(baseUrl)
	if err != nil {
		t.Errorf("Error marshaling JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handlers.HandleShortenURL().ServeHTTP(rr, req)

	if rr.Code != 201 {
		t.Errorf("Expected status code 201, got %d", rr.Code)
	}
}

// Test cases for HandleDeleteURL handler
func TestHandleDeleteURLSuccess(t *testing.T) {
	mockDB := mocks.NewMockDB()
	handlers := handlers.NewHandler(mockDB)

	baseUrl := models.BaseUrlInfo{
		Url: "http://localhost:8080/123456",
	}

	jsonBody, err := utils.JsonMarshal(baseUrl)
	if err != nil {
		t.Errorf("Error marshaling JSON: %v", err)
	}

	req := httptest.NewRequest("DELETE", "/delete", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handlers.HandleDeleteURL().ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	if rr.Body.String() != "Deleted" {
		t.Errorf("Expected response body to be 'Deleted', got %s", rr.Body.String())
	}
}

func TestHandleDeleteURLFailure(t *testing.T) {
	mockDB := mocks.NewMockDB()
	handlers := handlers.NewHandler(mockDB)

	baseUrl := models.BaseUrlInfo{
		Url: "http://localhost:8080/123457",
	}

	jsonBody, err := utils.JsonMarshal(baseUrl)
	if err != nil {
		t.Errorf("Error marshaling JSON: %v", err)
	}

	req := httptest.NewRequest("DELETE", "/delete", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handlers.HandleDeleteURL().ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Errorf("Expected status code 404, got %d", rr.Code)
	}
}
