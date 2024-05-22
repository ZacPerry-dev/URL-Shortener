package database

import (
	"URL-Shortener/internal/models"
	"context"
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDataBase interface {
	CloseConnection()
	GetURLCollection() *mongo.Collection
	GetURL(findVal, urlVal string) (models.NewUrlInfo, bool, error)
	PostURL(ctx context.Context, newUrl models.NewUrlInfo) error
	DeleteURL(ctx context.Context, urlKey string) error
}

type Database struct {
	db     *mongo.Client
	dbName string
}

func NewDatabase(dbUri, dbName string) (IDataBase, error) {
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

func (d *Database) GetURLCollection() *mongo.Collection {
	collection := d.db.Database(d.dbName).Collection("url-mappings")

	return collection
}

// TODO: Implement GET (FIND from utils), POST, DELETE here so it's easy to mock in my tests
func (d *Database) GetURL(findVal, urlVal string) (models.NewUrlInfo, bool, error) {
	var newUrl models.NewUrlInfo

	collection := d.GetURLCollection()
	result := collection.FindOne(context.Background(), bson.M{findVal: urlVal})

	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return newUrl, false, err
		}

		return newUrl, false, err
	}

	if err := result.Decode(&newUrl); err != nil {
		return newUrl, false, err
	}

	return newUrl, true, nil
}

func (d *Database) PostURL(ctx context.Context, newUrl models.NewUrlInfo) error {
	collection := d.GetURLCollection()

	if _, err := collection.InsertOne(ctx, newUrl); err != nil {
		return err
	}

	return nil
}

func (d *Database) DeleteURL(ctx context.Context, urlKey string) error {
	collection := d.GetURLCollection()

	if _, err := collection.DeleteOne(ctx, bson.M{"key": urlKey}); err != nil {
		return err
	}

	return nil
}
