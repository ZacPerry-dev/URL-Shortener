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
	// find the urlVal associated with the given findVal
	// So if findval == "longurl", urlVal == "https://www.google.com"
	// then we would return the url info associated with the long url

	for _, url := range m.Collection {
		// check if the urlVal equals the value of the findVal within the collection
		if findVal == "key" && url.Key == urlVal {
			return url, true, nil
		}
		if findVal == "longurl" && url.LongUrl == urlVal {
			return url, true, nil
		}
		if findVal == "shorturl" && url.ShortUrl == urlVal {
			return url, true, nil
		}
	}

	return models.NewUrlInfo{}, false, nil
}

func (m *MockDB) PostURL(ctx context.Context, newUrl models.NewUrlInfo) error {
	m.Collection[newUrl.Key] = newUrl
	return nil
}

func (m *MockDB) DeleteURL(ctx context.Context, urlKey string) error {
	delete(m.Collection, urlKey)
	return nil
}
