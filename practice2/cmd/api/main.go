package main

import (
	"encoding/json"
	"log"
	"net/http"
	"practice2/internal/handlers"
	"practice2/internal/middleware"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	store := handlers.NewStore()
	var tasksHandler http.Handler = http.HandlerFunc(store.TasksHandler)

	tasksHandler = middleware.Logging("request received")(tasksHandler)
	tasksHandler = middleware.APIKeyAuth("secret12345")(tasksHandler)

	mux.Handle("/tasks", tasksHandler)

	addr := ":8080"
	log.Printf("Server started on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
