package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type DBStore struct {
	db *sql.DB
}

func NewDBStore() (*DBStore, error) {
	db, err := sql.Open("sqlite", "idempotency.db")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	_, err = db.Exec("PRAGMA busy_timeout = 5000;")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS idempotency_keys (
		key TEXT PRIMARY KEY,
		status TEXT NOT NULL,
		response_code INTEGER,
		response_body TEXT
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &DBStore{db: db}, nil
}
func (s *DBStore) StartOrGet(key string) (string, int, string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", 0, "", err
	}

	result, err := tx.Exec(
		"INSERT OR IGNORE INTO idempotency_keys(key, status) VALUES (?, ?)",
		key,
		"processing",
	)
	if err != nil {
		tx.Rollback()
		return "", 0, "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return "", 0, "", err
	}

	if rowsAffected == 1 {
		err = tx.Commit()
		if err != nil {
			return "", 0, "", err
		}

		return "new", 0, "", nil
	}

	var status string
	var responseCode sql.NullInt64
	var responseBody sql.NullString

	err = tx.QueryRow(
		"SELECT status, response_code, response_body FROM idempotency_keys WHERE key = ?",
		key,
	).Scan(&status, &responseCode, &responseBody)

	if err != nil {
		tx.Rollback()
		return "", 0, "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", 0, "", err
	}

	if status == "processing" {
		return "processing", 0, "", nil
	}

	if status == "completed" {
		return "completed", int(responseCode.Int64), responseBody.String, nil
	}

	return status, 0, "", nil
}

func (s *DBStore) Finish(key string, statusCode int, body string) error {
	_, err := s.db.Exec(
		"UPDATE idempotency_keys SET status = ?, response_code = ?, response_body = ? WHERE key = ?",
		"completed",
		statusCode,
		body,
		key,
	)

	return err
}

func IdempotencyMiddleware(store *DBStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")

		if key == "" {
			http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
			return
		}

		status, code, body, err := store.StartOrGet(key)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if status == "processing" {
			http.Error(w, "Duplicate request in progress", http.StatusConflict)
			return
		}

		if status == "completed" {
			fmt.Println("Cached response returned")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			w.Write([]byte(body))
			return
		}

		recorder := httptest.NewRecorder()

		next.ServeHTTP(recorder, r)

		responseBody := recorder.Body.String()

		err = store.Finish(key, recorder.Code, responseBody)
		if err != nil {
			http.Error(w, "Failed to save response", http.StatusInternalServerError)
			return
		}

		for k, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Add(k, value)
			}
		}

		w.WriteHeader(recorder.Code)
		w.Write([]byte(responseBody))
	})
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing started")

	time.Sleep(2 * time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"paid","amount":1000,"transaction_id":"uuid-12345"}`))

	fmt.Println("Processing completed")
}

func main() {
	store, err := NewDBStore()
	if err != nil {
		log.Fatal(err)
	}

	handler := IdempotencyMiddleware(store, http.HandlerFunc(paymentHandler))

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &http.Client{}

	var wg sync.WaitGroup

	idempotencyKey := "same-key-123"

	fmt.Println("Sending 10 simultaneous requests with the same Idempotency-Key")

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(requestNumber int) {
			defer wg.Done()

			req, err := http.NewRequest(http.MethodPost, server.URL, nil)
			if err != nil {
				fmt.Println("Request error:", err)
				return
			}

			req.Header.Set("Idempotency-Key", idempotencyKey)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Client error:", err)
				return
			}
			defer resp.Body.Close()

			bodyBytes, _ := io.ReadAll(resp.Body)

			fmt.Printf("Request %d received status: %d, body: %s\n", requestNumber, resp.StatusCode, string(bodyBytes))
		}(i)
	}

	wg.Wait()

	fmt.Println()
	fmt.Println("Sending one more request after processing is completed")

	req, err := http.NewRequest(http.MethodPost, server.URL, nil)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}

	req.Header.Set("Idempotency-Key", idempotencyKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Client error:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	fmt.Printf("Final repeated request received status: %d, body: %s\n", resp.StatusCode, string(bodyBytes))
}
