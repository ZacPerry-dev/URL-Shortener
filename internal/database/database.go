package database

import (
	"context"
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	db     *mongo.Client
	dbName string
}

func NewDatabase(dbUri, dbName string) (*Database, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbUri))
	if err != nil {
		fmt.Println("Error connecting to DB!")
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected to DB!")
	return &Database{
		db:     client,
		dbName: dbName,
	}, nil
}

func (d *Database) CloseConnection() {
	if d.db == nil {
		return
	}

	err := d.db.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to DB closed successfully!")
}

func (d *Database) GetCollection(collectionName string) *mongo.Collection {
	collection := d.db.Database(d.dbName).Collection(collectionName)

	return collection
}
