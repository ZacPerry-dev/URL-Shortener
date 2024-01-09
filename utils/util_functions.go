package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"URL-Shortener/internal/models"
)

func ValidateURL(urlString string) (bool, string) {
	parsedURL, err := url.Parse(urlString)

	if err != nil || parsedURL == nil {
		return false, "Trouble Parsing URL"
	}

	if parsedURL.Host == "" {
		return false, "No Host"
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false, "Invalid URL Scheme"
	}

	return true, ""
}

func FindURL(findVal string, urlVal string, urlCollection *mongo.Collection) (models.NewUrlInfo, bool, error) {
	var newUrl models.NewUrlInfo

	result := urlCollection.FindOne(context.Background(), bson.M{findVal: urlVal})

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

func Hashing(baseUrl models.BaseUrlInfo, urlCollection *mongo.Collection) (string, error) {
	var key string

	for {
		key, _ = GenerateHashKey(baseUrl)
		result, _ := FindHashKey(key, urlCollection)
		if !result {
			break
		}
	}

	return key, nil
}

func GenerateHashKey(baseUrl models.BaseUrlInfo) (string, error) {
	hash := sha256.Sum256([]byte(baseUrl.Url))
	hashString := hex.EncodeToString(hash[:])

	start := rand.Intn(len(hashString) - 6)

	hashKey := hashString[start : start+6]

	return hashKey, nil
}

func FindHashKey(hashKey string, urlCollection *mongo.Collection) (bool, error) {
	result := urlCollection.FindOne(context.Background(), bson.M{"key": hashKey})

	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func CreateShortUrl(key string) string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	shortKey := "http://" + host + "/" + key

	return shortKey
}

// Pass an interface arg. Allows you to pass whatever you want to it
// This is kinda dumb in terms of robustness and error checking.
// Did this so i could pass either a string (newUrl.shortURL) or an entire struct to it (newUrl)
func JsonMarshal(newUrl interface{}) ([]byte, error) {
	res, err := json.Marshal(newUrl)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func CreateResponse(w http.ResponseWriter, httpStatus int, res []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(res)
}

func CreateErrorResponse(w http.ResponseWriter, httpMethod string, errorString string, httpStatus int) {
	w.Header().Set("Allow", httpMethod)
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, errorString, httpStatus)
}
