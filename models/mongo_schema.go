package models

type UrlCollection struct {
	URL      string `json:"url" bson:"url"`
	ShortURL string `json:"short_url" bson:"short_url"`
	Domain   string `json:"domain" bson:"domain"`
}

type DomainMetricsCollection struct {
	Domain  string `json:"domain" bson:"domain"`
	Counter int    `json:"counter" bson:"counter"`
}
