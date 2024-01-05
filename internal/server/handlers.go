package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/* TODO
Need to add some new handlers
Going to store in DB like so..

{
  "_id": ObjectId("unique_id"),
  "smallKey": "value1",
  "longKey": "value2"
}

1. Create a short URL
- Error check the passed URL
- First, check if the short URL key already exists
  - If so, handle accordingly
- Otherwise, call util hash function
- Determine if this hash key is in the DB. If so, regenerate (another function)
- Create new short URL
- Piece it together and store within the DB (mapping the short to real url)

2. Handler to check if it exists already
- call db and check if the fullform URL is already stored.

*/

func (s *Server) PostURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var baseUrl baseUrlInfo
	var newUrl newUrlInfo
	var urlCollection *mongo.Collection = s.db.GetCollection("url-mappings")

	// MOVE TO UTILS LATER I GUESS
	if err := json.NewDecoder(r.Body).Decode(&baseUrl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		w.Write([]byte(newUrl.ShortUrl))
		return
	}

	if err == mongo.ErrNoDocuments {
		// post
		w.Write([]byte("Does not exists. Creating in the DB"))
		// Call hash function
		// etc..
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If not, do the hashing stuff

	// If hash already exists, redo the hashing and continue until it gets a unique one

	// Then, store and return the shortened URL to the user

	w.Write([]byte(baseUrl.LongUrl))
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

func FindURL(baseUrl baseUrlInfo, urlCollection *mongo.Collection) (newUrlInfo, bool, error) {
	var newUrl newUrlInfo

	result := urlCollection.FindOne(context.Background(), bson.M{"LongUrl": baseUrl.LongUrl})

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

// hashing function
// Sha256 to convert the long url to a hash string
// Then, i want to generate a random index within the string and take a sequence of 6 chars
// THis will then be the key to be used for the short url, alongside being the ID in the db
func GenerateHashKey() {}

// Check the DB and see if the hash exists
func FindHashKey() {}

// Function for redirect
func RedirectURL() {}
