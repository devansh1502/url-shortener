package database

import (
	"reflect"
	"testing"
	"url-shortener/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateDB(t *testing.T) {
	testStore := NewStore()
	t.Run("Create Success", func(t *testing.T) {
		url := "https://www.google.com"
		shortUrl := "google.com/7378mDnD"

		val := testStore.Create(url, shortUrl)
		assert.Equal(t, val, true)
	})
	t.Run("Create Existing Domain Success", func(t *testing.T) {
		url := "https://www.google.com/1234"
		shortUrl := "google.com/7378mDnD"

		val := testStore.Create(url, shortUrl)
		assert.Equal(t, val, true)
	})
}

func TestDB_GetByURL(t *testing.T) {
	testStore := NewStore()
	t.Run("Get By URL Success", func(t *testing.T) {
		url := "https://www.google.com"
		shortUrl := "google.com/7378mDnD"
		testStore.Create(url, shortUrl)
		val := testStore.GetByURL(url)
		assert.Equal(t, val, shortUrl)
	})

	t.Run("Empty URL", func(t *testing.T) {
		url := ""
		val := testStore.GetByURL(url)
		assert.Equal(t, val, "")
	})
}

func TestDB_GetByShortURL(t *testing.T) {
	testStore := NewStore()
	t.Run("Get By Short URL Success", func(t *testing.T) {
		url := "https://www.google.com"
		shortUrl := "google.com/7378mDnD"
		testStore.Create(url, shortUrl)
		val := testStore.GetByShortURL(shortUrl)
		assert.Equal(t, val, url)
	})

	t.Run("Empty Short URL", func(t *testing.T) {
		ShortURL := ""
		val := testStore.GetByShortURL(ShortURL)
		assert.Equal(t, val, "")
	})
}
func TestDB_GetTopThreeDomains(t *testing.T) {
	dmc := []models.DomainMetricsCollection{
		{Domain: "youtube.com", Counter: 3},
		{Domain: "google.com", Counter: 2},
		{Domain: "infracloud.com", Counter: 2},
	}
	type fields struct {
		metricsMap map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   []models.DomainMetricsCollection
	}{
		{
			name: "Top Three Domains Success",
			fields: fields{
				metricsMap: map[string]int{"google.com": 2, "infracloud.com": 2, "youtube.com": 3, "facebook": 1},
			},
			want: dmc,
		},
		{
			name: "Less than three Domains Success",
			fields: fields{
				metricsMap: map[string]int{"infracloud.com": 2, "youtube.com": 3},
			},
			want: []models.DomainMetricsCollection{
				{Domain: "youtube.com", Counter: 3},
				{Domain: "infracloud.com", Counter: 2},
			},
		},
		{
			name: "Metrics Map Empty",
			fields: fields{
				metricsMap: map[string]int{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				metricsMap: tt.fields.metricsMap,
			}
			if got := db.GetTopThreeDomains(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.GetTopThreeDomains() = %v, want %v", got, tt.want)
			}
		})
	}
}
