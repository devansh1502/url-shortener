package interfaces

import (
	"net/http"
	"url-shortener/models"
)

// Store has all functions of the db.go as part of the interface.
type Store interface {
	Create(url, shortURl string) bool
	GetByURL(url string) string
	GetByShortURL(shortUrl string) string
	GetTopThreeDomains() []models.DomainMetricsCollection
}

// API has all functions like shortening and redirect as part of the interface
type API interface {
	RedirectURL(w http.ResponseWriter, r *http.Request)
	UrlShortner(w http.ResponseWriter, r *http.Request)
	Metrics(w http.ResponseWriter, r *http.Request)
}
