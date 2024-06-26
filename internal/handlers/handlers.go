package handlers

import (
	"URL-Shortener/internal/database"
	"URL-Shortener/internal/models"
	"context"
	"encoding/json"
	"net/http"

	"URL-Shortener/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	db database.IDataBase
}

func NewHandler(db database.IDataBase) *Handler {
	return &Handler{db: db}
}

// http://localhost:8080/home
func (h *Handler) HandleHomePageRequest() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the files
		files := []string{
			"views/index.html",
			"views/components/url_form.html",
		}

		// Parse
		tmpl := utils.ParseTemplates(files...)

		// Execute
		tmpl.ExecuteTemplate(w, "index", nil)
	})
}

// http://localhost:8080/{key}
func (h *Handler) HandleGetURL() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[1:]

		if key == "" {
			http.Redirect(w, r, "/home", http.StatusNoContent)
			return
		}

		if len(key) != 6 {
			http.Redirect(w, r, "/home", http.StatusBadRequest)
			return
		}

		if r.Method != http.MethodGet {
			utils.CreateErrorResponse(w, http.MethodGet, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var newUrl models.NewUrlInfo
		newUrl, status, err := h.db.GetURL("key", key)

		if status {
			w.Header().Set("Location", newUrl.LongUrl)
			utils.CreateResponse(w, http.StatusFound, []byte(newUrl.LongUrl))

			// I don't think I need this. If you curl with -L, it will follow the redirect
			// But, when using postman it does it anyway. Will keep for now though
			http.Redirect(w, r, newUrl.LongUrl, http.StatusFound)
			return
		}

		if err != nil && err != mongo.ErrNoDocuments {
			utils.CreateErrorResponse(w, http.MethodGet, "Error with the DB. Please Try Again", http.StatusBadRequest)
			return
		}

		utils.CreateErrorResponse(w, http.MethodGet, "Key not found", http.StatusNotFound)

	})
}

// http://localhost:8080/shorten
func (h *Handler) HandleShortenURL() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.CreateErrorResponse(w, http.MethodPost, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			utils.CreateErrorResponse(w, http.MethodPost, "Invalid Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		var baseUrl models.BaseUrlInfo
		var newUrl models.NewUrlInfo

		if err := json.NewDecoder(r.Body).Decode(&baseUrl); err != nil {
			utils.CreateErrorResponse(w, http.MethodPost, "Error Decoding JSON request body", http.StatusBadRequest)
			return
		}

		if baseUrl.Url == "" {
			utils.CreateErrorResponse(w, http.MethodPost, "Missing field: Url", http.StatusBadRequest)
			return
		}

		if status, errorString := utils.ValidateURL(baseUrl.Url); !status {
			utils.CreateErrorResponse(w, http.MethodPost, errorString, http.StatusBadRequest)
			return
		}

		// First, check if the URL already exists in the DB
		newUrl, status, err := h.db.GetURL("longurl", baseUrl.Url)
		if status {
			res, err := utils.JsonMarshal(newUrl)
			if err != nil {
				utils.CreateErrorResponse(w, http.MethodPost, "Error marshaling JSON. Please try again.", http.StatusBadRequest)
				return
			}

			utils.CreateResponse(w, http.StatusFound, res)
			w.Write([]byte("\n<a href='" + newUrl.ShortUrl + "'>" + newUrl.ShortUrl + "</a>\n"))
			return
		}

		if err != nil && err != mongo.ErrNoDocuments {
			utils.CreateErrorResponse(w, http.MethodPost, "Error with the DB. Please Try Again", http.StatusBadRequest)
			return
		}

		var key string

		if err == mongo.ErrNoDocuments {
			key, _ = utils.Hashing(baseUrl, h.db)
		}

		// Then, store and return the shortened URL to the user
		newUrl = models.NewUrlInfo{
			LongUrl:  baseUrl.Url,
			ShortUrl: utils.CreateShortUrl(key),
			Key:      key,
		}

		if err := h.db.PostURL(context.Background(), newUrl); err != nil {
			utils.CreateErrorResponse(w, http.MethodPost, "Error saving in the DB. Please try again.", http.StatusBadRequest)
			return
		}

		res, err := utils.JsonMarshal(newUrl)
		if err != nil {
			utils.CreateErrorResponse(w, http.MethodPost, "Error marshaling JSON. Please try again.", http.StatusBadRequest)
			return
		}

		utils.CreateResponse(w, http.StatusCreated, res)
		w.Write([]byte("\n<a href='" + newUrl.ShortUrl + "'>" + newUrl.ShortUrl + "</a>\n"))
	})
}

// http://localhost:8080/delete
func (h *Handler) HandleDeleteURL() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.CreateErrorResponse(w, http.MethodDelete, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			utils.CreateErrorResponse(w, http.MethodDelete, "Invalid Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		var baseUrl models.BaseUrlInfo
		var newUrl models.NewUrlInfo

		if err := json.NewDecoder(r.Body).Decode(&baseUrl); err != nil {
			utils.CreateErrorResponse(w, http.MethodDelete, "Error Decoding JSON request body", http.StatusBadRequest)
			return
		}

		if baseUrl.Url == "" {
			utils.CreateErrorResponse(w, http.MethodDelete, "Missing field: Url", http.StatusBadRequest)
			return
		}

		if status, errorString := utils.ValidateURL(baseUrl.Url); !status {
			utils.CreateErrorResponse(w, http.MethodDelete, errorString, http.StatusBadRequest)
			return
		}

		// First, check if the URL already exists in the DB
		newUrl, status, err := h.db.GetURL("shorturl", baseUrl.Url)

		// If found, delete
		if status {
			err := h.db.DeleteURL(context.Background(), newUrl.Key)
			if err != nil {
				utils.CreateErrorResponse(w, http.MethodDelete, "Error with the DB. Please Try Again", http.StatusBadRequest)
				return
			}
			utils.CreateResponse(w, http.StatusOK, []byte("Deleted"))
			return
		}

		if err != nil && err != mongo.ErrNoDocuments {
			utils.CreateErrorResponse(w, http.MethodGet, "Error with the DB. Please Try Again", http.StatusBadRequest)
			return
		}

		utils.CreateErrorResponse(w, http.MethodGet, "Key not found", http.StatusNotFound)
	})
}
