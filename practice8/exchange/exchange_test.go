package exchange

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	t.Run("Successful scenario", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"base":"USD","target":"EUR","rate":0.85}`)
		}))
		defer server.Close()

		svc := NewExchangeService(server.URL)
		rate, err := svc.GetRate("USD", "EUR")

		assert.NoError(t, err)
		assert.Equal(t, 0.85, rate)
	})

	t.Run("API Business Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error":"invalid currency pair"}`)
		}))
		defer server.Close()

		svc := NewExchangeService(server.URL)
		_, err := svc.GetRate("USD", "UNKNOWN")

		assert.EqualError(t, err, "api error: invalid currency pair")
	})

	t.Run("Malformed JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"base":"USD", "target":"EUR", "rate":0.85`) // Truncated
		}))
		defer server.Close()

		svc := NewExchangeService(server.URL)
		_, err := svc.GetRate("USD", "EUR")

		assert.ErrorContains(t, err, "decode error")
	})

	t.Run("Slow Response/Timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond) // Simulate slow response
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		svc := NewExchangeService(server.URL)
		svc.Client.Timeout = 1 * time.Millisecond // Force timeout quicker for test

		_, err := svc.GetRate("USD", "EUR")
		assert.ErrorContains(t, err, "network error")
	})

	t.Run("Server Panic / 500 Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"base":"","target":"","rate":0}`)
		}))
		defer server.Close()

		svc := NewExchangeService(server.URL)
		_, err := svc.GetRate("USD", "EUR")

		assert.EqualError(t, err, "unexpected status: 500")
	})

	t.Run("Empty Body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		svc := NewExchangeService(server.URL)
		_, err := svc.GetRate("USD", "EUR")

		assert.ErrorContains(t, err, "decode error")
	})
}