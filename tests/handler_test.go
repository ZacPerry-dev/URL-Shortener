package tests

import (
	"URL-Shortener/tests/mocks"
	"testing"
)

func TestMockImplementation(t *testing.T) {
	mockDB := mocks.NewMockDB()
	// handlers := handlers.NewHandler(mockDB)

	// test that mock db collection is not empty
	if len(mockDB.Collection) == 0 {
		t.Error("MockDB collection is empty")
	}

	// test that mock db collection has the correct key
	if _, ok := mockDB.Collection["123456"]; !ok {
		t.Error("MockDB collection does not have the correct key")
	}

}

// TODO: Add more tests for the mock DB implementation
