package server

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
)

func (s *Server) PostURL(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var baseUrl baseUrlInfo
	var newUrl newUrlInfo
	var urlCollection *mongo.Collection = s.db.GetCollection("url-mappings")

	// MOVE TO UTILS LATER I GUESS
	if err := json.NewDecoder(r.Body).Decode(&baseUrl); err != nil {
		http.Error(w, "Error Decoding JSON request body", http.StatusBadRequest)
		return
	}

	if baseUrl.LongUrl == "" {
		http.Error(w, "Missing field: URL", http.StatusBadRequest)
		return
	}

	if status, errorString := ValidateURL(baseUrl.LongUrl); !status {
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	// check if it exists
	// REFACTOR ALL THIS GROSS STUFF
	newUrl, status, err := FindURL(baseUrl, urlCollection)
	if status {
		http.Redirect(w, r, newUrl.LongUrl, http.StatusFound)
		return
	}

	if err != nil && err != mongo.ErrNoDocuments {
		http.Error(w, "Error with the DB. Please Try Again", http.StatusBadRequest)
		return
	}

	var key string
	if err == mongo.ErrNoDocuments {
		w.Write([]byte("Does not exists. Creating in the DB\n"))
		key, _ = Hashing(baseUrl, urlCollection)
	}

	// Then, store and return the shortened URL to the user
	newUrl = newUrlInfo{
		LongUrl:  baseUrl.LongUrl,
		ShortUrl: CreateShortUrl(key),
		Key:      key,
	}

	if _, err := urlCollection.InsertOne(context.Background(), newUrl); err != nil {
		http.Error(w, "Error saving in the DB. Please try again.", http.StatusBadRequest)
	}

	res, err := JsonMarshal(newUrl)
	if err != nil {
		http.Error(w, "Error marshaling JSON. Please try again.", http.StatusBadRequest)
		return
	}

	CreateResponse(w, http.StatusCreated, res)
}

// I CANT DECIDE BUT MAYBE MOVE THESE TO UTILS TOO
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

// Refactor this
func FindURL(baseUrl baseUrlInfo, urlCollection *mongo.Collection) (newUrlInfo, bool, error) {
	var newUrl newUrlInfo

	result := urlCollection.FindOne(context.Background(), bson.M{"longurl": baseUrl.LongUrl})

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

func Hashing(baseUrl baseUrlInfo, urlCollection *mongo.Collection) (string, error) {
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

func GenerateHashKey(baseUrl baseUrlInfo) (string, error) {
	hash := sha256.Sum256([]byte(baseUrl.LongUrl))
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

func JsonMarshal(newUrl newUrlInfo) ([]byte, error) {
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
