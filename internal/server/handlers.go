package server

import (
	"context"
	"encoding/json"
	"fmt"
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

func (s *Server) postURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var baseUrl baseUrlInfo

	// Decode, move to utils later
	err := json.NewDecoder(r.Body).Decode(&baseUrl)
	if err != nil {
		http.Error(w, "Invalid JSON Payload", http.StatusBadRequest)
		return
	}

	if baseUrl.LongUrl == "" {
		http.Error(w, "Missing field: URL", http.StatusBadRequest)
		return
	}

	status, errorString := checkURL(baseUrl.LongUrl)
	if !status {
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	// Check if it exists already
	var urlCollection *mongo.Collection = s.db.GetCollection("url-mappings")
	// var newUrl newUrlInfo

	count, err := urlCollection.CountDocuments(context.Background(), bson.M{"LongUrl": baseUrl.LongUrl})
	if err != nil {
		http.Error(w, "Error with DB. Please try again later...", http.StatusInternalServerError)
		return
	}

	// TODO: Abstract this to handle. If it already exists, then just call the get function and return I guess idk
	if count != 0 {
		fmt.Println("Url has already been converted... here -> ")
		return
	}

	// If not, do the hashing stuff

	// If hash already exists, redo the hashing and continue until it gets a unique one

	// Then, store and return the shortened URL to the user

	w.Write([]byte(baseUrl.LongUrl))
}

// TODO: Need to figure out how I wanna structure this. Move some of the "getting" logic
// from above to this function and return all info if it exists. If not, add and then return idk yet
func (s *Server) getURL(w http.ResponseWriter, r *http.Request) {
}

func checkURL(urlString string) (bool, string) {
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
