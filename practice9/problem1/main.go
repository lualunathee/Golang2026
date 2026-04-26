package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"time"
)

type PaymentClient struct {
	HTTPClient *http.Client
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return true
		}
		return true
	}

	if resp == nil {
		return false
	}

	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	case 401, 404:
		return false
	default:
		return false
	}
}

func (c *PaymentClient) CalculateBackoff(attempt int) time.Duration {
	backoff := c.BaseDelay * time.Duration(1<<attempt)

	if backoff > c.MaxDelay {
		backoff = c.MaxDelay
	}

	jitter := time.Duration(rand.Int63n(int64(backoff)))
	return jitter
}

func (c *PaymentClient) ExecutePayment(ctx context.Context, url string) error {
	for attempt := 0; attempt < c.MaxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			return err
		}

		resp, err := c.HTTPClient.Do(req)

		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()

			var result map[string]string
			json.NewDecoder(resp.Body).Decode(&result)

			fmt.Printf("Attempt %d: Success! Response: %v\n", attempt+1, result)
			return nil
		}

		if !IsRetryable(resp, err) {
			return fmt.Errorf("non-retryable error, status: %v, err: %v", resp.StatusCode, err)
		}

		if attempt == c.MaxRetries-1 {
			break
		}

		delay := c.CalculateBackoff(attempt)

		fmt.Printf("Attempt %d failed: waiting %v...\n", attempt+1, delay)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("payment failed after %d attempts", c.MaxRetries)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount <= 3 {
			fmt.Printf("Server: request %d -> 503 Service Unavailable\n", requestCount)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}

		fmt.Printf("Server: request %d -> 200 OK\n", requestCount)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))

	defer server.Close()

	client := PaymentClient{
		HTTPClient: &http.Client{},
		MaxRetries: 5,
		BaseDelay:  500 * time.Millisecond,
		MaxDelay:   5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.ExecutePayment(ctx, server.URL)

	if err != nil {
		fmt.Println("Final result:", err)
	} else {
		fmt.Println("Final result: payment completed successfully")
	}
}
