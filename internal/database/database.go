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
}

type service struct {
	db *mongo.Client
}

var (
	// host     = os.Getenv("DB_HOST")
	// port     = os.Getenv("DB_PORT")
	database = os.Getenv("DB_URI")
	db_name  = os.Getenv("DB_NAME")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(database))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to DB!")

	return &service{
		db: client,
	}
}

func (s *service) GetCollection(collectionName string) *mongo.Collection {
	collection := s.db.Database(db_name).Collection(collectionName)

	return collection
}
