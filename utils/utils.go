package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"log"
	"regexp"
)

// GetDomain returns the domain name from the url
func GetDomain(url string) string {
	if url == "" {
		log.Print("URL Cannot be empty")
		return ""
	}
	m := regexp.MustCompile(`\.?([^.]*.com)`)
	domain := m.FindStringSubmatch(url)[1]

	return domain
}

// ShortenURL returns the shortened url with the domain name as the prefix
func ShortenURL(url string) string {
	hash := sha1.New()
	_, err := hash.Write([]byte(url))
	if err != nil {
		log.Printf("Unable to write hash. %v", err)
		return ""
	}
	shortURL := base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]

	return GetDomain(url) + "/" + shortURL

}
