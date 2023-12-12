package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/addURL", s.addURL)
	mux.HandleFunc("/getURL", s.getURL)

	return mux
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

// func (s *Server) getCollections(w http.ResponseWriter, r *http.Request) {
// 	collection := s.db.GetCollection("url-mappings")

// 	fmt.Println("Poggers, got the collection: ", collection)
// }

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

func (s *Server) addURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		// Calls write header and write behind the scenes for you. Much easier to just use this I guess.
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Creating new url "))
}

func (s *Server) getURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		// Calls write header and write behind the scenes for you. Much easier to just use this I guess.
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Getting existing url...."))
}
