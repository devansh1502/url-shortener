package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	mocks "url-shortener/mocks/interfaces"
	"url-shortener/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAPI(t *testing.T) {
	testSKNotFound(t)
	testShortURLNotFound(t)
	testRedirectURL(t)
	testMethod(t)
	testEmptyURL(t)
	testExistingURL(t)
	testCreateURL(t)
	testCreateURLFailedCase(t)
	testTopThreeDomains(t)
}

func testSKNotFound(t *testing.T) {
	testContext := context.Background()
	testStore := mocks.NewStore(t)
	testAPI := NewAPI(testContext, testStore)

	t.Run("Missing Short Key Redirect", func(t *testing.T) {
		shortKey := ""
		req := httptest.NewRequest(http.MethodGet, "/redirect/"+shortKey, nil)
		w := httptest.NewRecorder()
		testAPI.RedirectURL(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusNotFound)
	})
}

func testShortURLNotFound(t *testing.T) {
	testContext := context.Background()
	testStore := mocks.NewStore(t)
	testAPI := NewAPI(testContext, testStore)

	t.Run("Short URL not Found Redirect", func(t *testing.T) {
		shortKey := "google.com/7378mDnD"
		testStore.On("GetByShortURL", shortKey).Return("")

		req := httptest.NewRequest(http.MethodGet, "/redirect/"+shortKey, nil)
		w := httptest.NewRecorder()
		testAPI.RedirectURL(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusNotFound)
	})
}

func testRedirectURL(t *testing.T) {
	testContext := context.Background()
	testStore := mocks.NewStore(t)
	testAPI := NewAPI(testContext, testStore)

	t.Run("Redirect Success", func(t *testing.T) {
		shortKey := "google.com/7378mDnD"
		testStore.On("GetByShortURL", shortKey).Return("https://www.google.com")

		req := httptest.NewRequest(http.MethodGet, "/redirect/"+shortKey, nil)
		w := httptest.NewRecorder()
		testAPI.RedirectURL(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusMovedPermanently)
	})
}

func testMethod(t *testing.T) {
	testContext := context.Background()
	testStore := mocks.NewStore(t)
	testAPI := NewAPI(testContext, testStore)

	t.Run("Wrong Method", func(t *testing.T) {
		testURL := "www.google.com"

		req := httptest.NewRequest(http.MethodGet, "/short/"+testURL, nil)
		w := httptest.NewRecorder()
		testAPI.UrlShortner(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		exData, _ := json.Marshal(map[string]string{"Error": "Method not Supported!"})
		assert.Equal(t, exData, data)
	})
}

func testEmptyURL(t *testing.T) {
	testContext := context.Background()
	testStore := mocks.NewStore(t)
	testAPI := NewAPI(testContext, testStore)

	t.Run("URL is Empty", func(t *testing.T) {
		testURL := ""
		req := httptest.NewRequest(http.MethodPost, "/short/"+testURL, nil)
		w := httptest.NewRecorder()
		testAPI.UrlShortner(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		exData, _ := json.Marshal(map[string]string{"Error": "URL is Empty!"})
		assert.Equal(t, exData, data)
	})
}

func testExistingURL(t *testing.T) {
	testContext := context.Background()
	testStore := mocks.NewStore(t)
	testAPI := NewAPI(testContext, testStore)

	t.Run("Get Existing URL", func(t *testing.T) {
		testURL := "https://www.google.com"
		testStore.On("GetByURL", testURL).Return("google.com/7378mDnD").Once()

		req := httptest.NewRequest(http.MethodPost, "/short/"+testURL, nil)
		w := httptest.NewRecorder()
		testAPI.UrlShortner(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		exData, _ := json.Marshal(map[string]string{"short_url": "google.com/7378mDnD"})
		assert.Equal(t, exData, data)
	})
}

func testCreateURL(t *testing.T) {
	testContext := context.Background()
	testStore := mocks.NewStore(t)
	testAPI := NewAPI(testContext, testStore)

	t.Run("Create Short URL", func(t *testing.T) {
		testURL := "www.google.com"
		shortURL := "google.com/7378mDnD"
		testStore.On("GetByURL", "https://"+testURL).Return("").Once()
		testStore.On("Create", "https://"+testURL, shortURL).Return(true).Once()

		req := httptest.NewRequest(http.MethodPost, "/short/"+testURL, nil)
		w := httptest.NewRecorder()
		testAPI.UrlShortner(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		exData, _ := json.Marshal(map[string]string{"short_url": "google.com/7378mDnD"})
		assert.Equal(t, exData, data)
	})
}

func testCreateURLFailedCase(t *testing.T) {
	testContext := context.Background()
	testStore := mocks.NewStore(t)
	testAPI := NewAPI(testContext, testStore)

	t.Run("Failed to Create Short URL", func(t *testing.T) {
		testURL := "https://www.google.com"
		testStore.On("GetByURL", testURL).Return("").Once()
		testStore.On("Create", testURL, mock.Anything).Return(false).Once()

		req := httptest.NewRequest(http.MethodPost, "/short/"+testURL, nil)
		w := httptest.NewRecorder()
		testAPI.UrlShortner(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		exData, _ := json.Marshal(map[string]string{"Error": "Failed to Shorten the URl!"})
		assert.Equal(t, exData, data)
	})
}

func testTopThreeDomains(t *testing.T) {
	testContext := context.Background()
	testStore := mocks.NewStore(t)
	testAPI := NewAPI(testContext, testStore)

	t.Run("Top Three Domains", func(t *testing.T) {
		dmc := []models.DomainMetricsCollection{
			{Domain: "youtube.com", Counter: 3},
			{Domain: "google.com", Counter: 2},
			{Domain: "infracloud.com", Counter: 2},
		}
		testStore.On("GetTopThreeDomains").Return(dmc).Once()

		req := httptest.NewRequest(http.MethodGet, "/metrics/", nil)
		w := httptest.NewRecorder()
		testAPI.Metrics(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		exData, _ := json.MarshalIndent(dmc, "", " ")
		assert.Equal(t, exData, data)
	})
}
