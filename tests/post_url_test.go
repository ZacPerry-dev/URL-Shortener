package tests

import (
	"URL-Shortener/internal/database"
	"URL-Shortener/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"URL-Shortener/internal/server"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// Test for:
// Valid response if the url is created in the DB (201)
// Valid response if the url is found in the DB (302)
// Error if passed an invalid URL (400)
// Error if passed no url (Missing field, 400)

func TestPostUrl(t *testing.T) {
	mockDB := &database.MockDatabase{
		GetCollectionFunc: func(s string) *mongo.Collection {
			return &mongo.Collection{}
		},
	}

	s := server.NewServer(mockDB, "../views")

	// mockCollection := &mongo.Collection{}
	// mockCollection.On("FindOne").Return(&mongo.SingleResult{}, nil)
	// mockDB.SetMockCollection(mockCollection)

	// db := s.DB().(*database.MockDatabase)

	// assert.NotNil(t, db)

	mockCollection := &mock.Mock{}
	mockCollection.On("FindOne").Return(&mongo.SingleResult{}, nil)
	// s.db.(*database.MockDatabase).SetMockCollection(mockCollection)

	requestBody := models.BaseUrlInfo{
		Url: "https://www.example.com",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/addURL", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	s.PostURL(res, req)

	var response models.NewUrlInfo
	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// assert that nil is returned
	// assert.Equal(t, http.StatusFound, res.Code)
	// assert.Equal(t, "https://www.example.com", response.LongUrl)
	// assert.NotEmpty(t, response.ShortUrl)
	// assert.NotEmpty(t, response.Key)
}
