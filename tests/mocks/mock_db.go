package mocks

import (
	"URL-Shortener/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MockDB struct {
	Collection map[string]models.NewUrlInfo
}

func NewMockDB() *MockDB {
	return &MockDB{
		Collection: map[string]models.NewUrlInfo{
			"123456": {
				Key:      "123456",
				LongUrl:  "https://www.google.com",
				ShortUrl: "http://localhost:8080/123456",
			},
		},
	}
}

func (m *MockDB) CloseConnection() {
	return
}

func (m *MockDB) GetURLCollection() *mongo.Collection {
	return nil
}

func (m *MockDB) GetURL(findVal, urlVal string) (models.NewUrlInfo, bool, error) {
	url, ok := m.Collection[urlVal]
	if !ok {
		return models.NewUrlInfo{}, false, nil
	}

	return url, true, nil
}

func (m *MockDB) PostURL(ctx context.Context, newUrl models.NewUrlInfo) error {
	m.Collection[newUrl.Key] = newUrl
	return nil
}

func (m *MockDB) DeleteURL(ctx context.Context, urlKey string) error {
	delete(m.Collection, urlKey)
	return nil
}
