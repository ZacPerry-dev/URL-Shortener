package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*

TODO:
Container
Tests
Error checking (return JSON for testing purposes)
Util cleanup
*/

func (s *Server) RedirectURL(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	if key == "" {
		fmt.Fprint(w, "No key provided")
		return
	}

	if len(key) != 6 {
		fmt.Fprint(w, "Invalid key provided")
		return
	}

	if r.Method != http.MethodGet {
		CreateErrorResponse(w, http.MethodGet, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUrl newUrlInfo
	var urlCollection *mongo.Collection = s.db.GetCollection("url-mappings")

	newUrl, status, err := FindURL("key", key, urlCollection)

	if status {
		w.Header().Set("Location", newUrl.LongUrl)
		CreateResponse(w, http.StatusFound, []byte(newUrl.LongUrl))

		// I don't think I techinally need this. If you curl with -L, it will follow the redirect
		// But, when using postman it does it anyway. Will keep for now though
		http.Redirect(w, r, newUrl.LongUrl, http.StatusFound)
		return
	}

	if err != nil && err != mongo.ErrNoDocuments {
		CreateErrorResponse(w, http.MethodGet, "Error with the DB. Please Try Again", http.StatusBadRequest)
		return
	}

	CreateErrorResponse(w, http.MethodGet, "Key not found", http.StatusNotFound)
}

func (s *Server) PostURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		CreateErrorResponse(w, http.MethodPost, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		CreateErrorResponse(w, http.MethodPost, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var baseUrl baseUrlInfo
	var newUrl newUrlInfo
	var urlCollection *mongo.Collection = s.db.GetCollection("url-mappings")

	if err := json.NewDecoder(r.Body).Decode(&baseUrl); err != nil {
		CreateErrorResponse(w, http.MethodPost, "Error Decoding JSON request body", http.StatusBadRequest)
		return
	}

	if baseUrl.Url == "" {
		CreateErrorResponse(w, http.MethodPost, "Missing field: Url", http.StatusBadRequest)
		return
	}

	if status, errorString := ValidateURL(baseUrl.Url); !status {
		CreateErrorResponse(w, http.MethodPost, errorString, http.StatusBadRequest)
		return
	}

	// First, check if the URL already exists in the DB
	newUrl, status, err := FindURL("longurl", baseUrl.Url, urlCollection)
	if status {
		res, err := JsonMarshal(newUrl)
		if err != nil {
			CreateErrorResponse(w, http.MethodPost, "Error marshaling JSON. Please try again.", http.StatusBadRequest)
			return
		}

		CreateResponse(w, http.StatusFound, res)
		return
	}

	if err != nil && err != mongo.ErrNoDocuments {
		CreateErrorResponse(w, http.MethodPost, "Error with the DB. Please Try Again", http.StatusBadRequest)
		return
	}

	var key string

	if err == mongo.ErrNoDocuments {
		key, _ = Hashing(baseUrl, urlCollection)
	}

	// Then, store and return the shortened URL to the user
	newUrl = newUrlInfo{
		LongUrl:  baseUrl.Url,
		ShortUrl: CreateShortUrl(key),
		Key:      key,
	}

	if _, err := urlCollection.InsertOne(context.Background(), newUrl); err != nil {
		CreateErrorResponse(w, http.MethodPost, "Error saving in the DB. Please try again.", http.StatusBadRequest)
		return
	}

	res, err := JsonMarshal(newUrl)
	if err != nil {
		CreateErrorResponse(w, http.MethodPost, "Error marshaling JSON. Please try again.", http.StatusBadRequest)
		return
	}

	CreateResponse(w, http.StatusCreated, res)
}

func (s *Server) DeleteURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		CreateErrorResponse(w, http.MethodDelete, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		CreateErrorResponse(w, http.MethodDelete, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var baseUrl baseUrlInfo
	var newUrl newUrlInfo
	var urlCollection *mongo.Collection = s.db.GetCollection("url-mappings")

	if err := json.NewDecoder(r.Body).Decode(&baseUrl); err != nil {
		CreateErrorResponse(w, http.MethodDelete, "Error Decoding JSON request body", http.StatusBadRequest)
		return
	}

	if baseUrl.Url == "" {
		CreateErrorResponse(w, http.MethodDelete, "Missing field: Url", http.StatusBadRequest)
		return
	}

	if status, errorString := ValidateURL(baseUrl.Url); !status {
		CreateErrorResponse(w, http.MethodDelete, errorString, http.StatusBadRequest)
		return
	}

	newUrl, status, err := FindURL("shorturl", baseUrl.Url, urlCollection)

	// If found, delete
	if status {
		_, err := urlCollection.DeleteOne(context.Background(), bson.M{"key": newUrl.Key})
		if err != nil {
			CreateErrorResponse(w, http.MethodDelete, "Error with the DB. Please Try Again", http.StatusBadRequest)
			return
		}
		CreateResponse(w, http.StatusOK, []byte("Deleted"))
		return
	}

	if err != nil && err != mongo.ErrNoDocuments {
		CreateErrorResponse(w, http.MethodGet, "Error with the DB. Please Try Again", http.StatusBadRequest)
		return
	}

	CreateErrorResponse(w, http.MethodGet, "Key not found", http.StatusNotFound)
}

// I CANT DECIDE BUT MAYBE MOVE THESE TO UTILS
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

func FindURL(findVal string, urlVal string, urlCollection *mongo.Collection) (newUrlInfo, bool, error) {
	var newUrl newUrlInfo

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
