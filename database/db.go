package database

import (
	"log"
	"sort"
	"url-shortener/interfaces"
	"url-shortener/models"
	"url-shortener/utils"
)

// DB struct contains the url,ShortURL map and metrics map data types
type DB struct {
	urlMap     map[string]string
	metricsMap map[string]int
}

// NewStore returns an entry of the Store interface
func NewStore() interfaces.Store {
	db := &DB{
		urlMap:     make(map[string]string),
		metricsMap: make(map[string]int),
	}
	return db
}

// Create this function is used to add the entry in the maps for url and shortURl
// and also adds the entry for the metric in the metrics map
func (db *DB) Create(url, shortURl string) bool {
	domain := utils.GetDomain(url)
	db.urlMap[url] = shortURl

	if db.metricsMap[domain] != 0 {
		value := db.metricsMap[domain] + 1
		db.metricsMap[domain] = value
	} else {
		db.metricsMap[domain] = 1
	}

	return true
}

// GetByURL this function is used to get the value of the shortURL w.r.t to the URL
func (db *DB) GetByURL(url string) string {
	if len(url) < 1 {
		log.Println("Url can not be empty!")
		return ""
	}

	return db.urlMap[url]
}

// GetByShortURL this function is used to get the value of the url w.r.t to the ShortURL
func (db *DB) GetByShortURL(shortUrl string) string {
	if len(shortUrl) < 1 {
		log.Println("Url can not be empty!")
		return ""
	}

	for k, v := range db.urlMap {
		if v == shortUrl {
			return k
		}
	}
	return ""
}

// GetTopThreeDomains lists down the top three most hit domains
func (db *DB) GetTopThreeDomains() []models.DomainMetricsCollection {
	keys := make([]string, 0, len(db.metricsMap))
	for k := range db.metricsMap {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return db.metricsMap[keys[i]] > db.metricsMap[keys[j]]
	})

	var dmc []models.DomainMetricsCollection
	i := 0
	for _, k := range keys {
		if i == 3 {
			return dmc
		}
		dmc = append(dmc, models.DomainMetricsCollection{Domain: k, Counter: db.metricsMap[k]})
		i++
	}
	return nil
}
