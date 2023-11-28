package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"regexp"
)

// GetDomain returns the domain name from the url
func GetDomain(url string) string {
	if url == "" {
		return ""
	}
	m := regexp.MustCompile(`\.?([^.]*.com)`)
	domain := m.FindStringSubmatch(url)[1]

	return domain
}

// ShortenURL returns the shortened url with the domain name as the prefix
func ShortenURL(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	shortURL := base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]

	return GetDomain(url) + "/" + shortURL

}
