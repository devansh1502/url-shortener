package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"url-shortener/interfaces"
	"url-shortener/utils"
)

type API struct {
	ctx context.Context
	db  interfaces.Store
}

func NewAPI(ctx context.Context, db interfaces.Store) interfaces.API {
	a := &API{
		ctx: ctx,
		db:  db,
	}
	return a
}

// RedirectURL redirects the shortened url to the original url
func (a *API) RedirectURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Redirect URL", r.URL.String())
	shortKey := r.URL.Path[len("/redirect/"):]
	if shortKey == "" {
		http.Error(w, "Short key is missing", http.StatusNotFound)
		return
	}

	url := a.db.GetByShortURL(shortKey)
	if url == "" {
		http.Error(w, "Shorten URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

// URLShortner returns a shorten url of the original url
func (a *API) UrlShortner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse, _ := json.Marshal(map[string]string{"Error": "Method not Supported!"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)

		return
	}

	url := r.URL.String()
	finalUrl := strings.Split(url, "/short/")[1]
	if finalUrl == "" {
		jsonResponse, _ := json.Marshal(map[string]string{"Error": "URL is Empty!"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)

		return
	}

	if !strings.Contains(finalUrl, "https://") {
		finalUrl = "https://" + finalUrl
	}
	shortUrl := utils.ShortenURL(finalUrl)

	existingURL := a.db.GetByURL(finalUrl)
	if existingURL != "" {
		jsonResponse, _ := json.Marshal(map[string]string{"short_url": existingURL})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)

		return
	}

	created := a.db.Create(finalUrl, shortUrl)
	if !created {
		jsonResponse, _ := json.Marshal(map[string]string{"Error": "Failed to Shorten the URl!"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)

		return
	}

	jsonResponse, _ := json.Marshal(map[string]string{"short_url": shortUrl})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)

	return
}

// Metrics returns the top three domains
func (a *API) Metrics(w http.ResponseWriter, r *http.Request) {
	topThree := a.db.GetTopThreeDomains()

	// Using Marshal Indent for formatting the JSON Response
	jsonResponse, _ := json.MarshalIndent(topThree, "", " ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
