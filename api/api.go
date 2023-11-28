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
	url := r.URL.String()
	finalURL:= strings.Split(url,"/short/")[1]
	shortURL := utils.ShortenURL(finalURL)

	if a.db.Create(finalURL, shortURL) {
		jsonResponse, _ := json.Marshal(map[string]string{"short_url": shortURL})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}
