package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MockDatabase struct {
	HealthFunc        func() map[string]string
	GetCollectionFunc func(string) *mongo.Collection
	mockCollection    *mongo.Collection
}

// Health implements Service.
func (m *MockDatabase) Health() map[string]string {
	return m.HealthFunc()
}

func (m *MockDatabase) SetMockCollection(collection *mongo.Collection) {
	m.mockCollection = collection
}

func (m *MockDatabase) GetCollection(collectionName string) *mongo.Collection {
	if m.GetCollectionFunc != nil {
		return m.GetCollectionFunc(collectionName)
	}
	return m.mockCollection
}
