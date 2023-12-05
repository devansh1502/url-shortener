package utils

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGetDomain(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Passing Empty URL",
			args: args{url: ""},
			want: "",
		},
		{
			name: "Get Domain Success",
			args: args{url: "www.youtube.com"},
			want: "youtube.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDomain(tt.args.url); got != tt.want {
				t.Errorf("GetDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortenURL(t *testing.T) {
	t.Run("Get Shortened URL", func(t *testing.T) {
		url := "https://www.google.com"
		regex := "google\\.com/[A-Za-z0-9]+"
		su := ShortenURL(url)
		assert.Matches(t, su, regex)
	})
}
