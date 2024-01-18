package database

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	GetCollection(string) *mongo.Collection
	CloseConnection()
}

type service struct {
	db *mongo.Client
}

// var (
// 	// host     = os.Getenv("DB_HOST")
// 	// port     = os.Getenv("DB_PORT")
// 	database = os.Getenv("DB_URI")
// 	db_name  = os.Getenv("DB_NAME")
// )

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("DB_URI")))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to DB!")

	return &service{
		db: client,
	}
}

func (s *service) GetCollection(collectionName string) *mongo.Collection {
	collection := s.db.Database(os.Getenv("DB_NAME")).Collection(collectionName)

	return collection
}

func (s *service) CloseConnection() {
	if s.db == nil {
		return
	}

	err := s.db.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to DB closed")
}
