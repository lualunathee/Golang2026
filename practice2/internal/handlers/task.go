package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Store struct {
	mu     sync.Mutex
	nextID int
	tasks  map[int]Task
}

func NewStore() *Store {
	return &Store{
		nextID: 1,
		tasks:  make(map[int]Task),
	}
}

// Общий handler для /tasks: внутри разруливаем GET/POST/PATCH
func (s *Store) TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGet(w, r)
	case http.MethodPost:
		s.handlePost(w, r)
	case http.MethodPatch:
		s.handlePatch(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// GET /tasks  или  GET /tasks?id=1
func (s *Store) handleGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	s.mu.Lock()
	defer s.mu.Unlock()

	// без id -> вернуть все
	if idStr == "" {
		out := make([]Task, 0, len(s.tasks))
		for _, t := range s.tasks {
			out = append(out, t)
		}
		writeJSON(w, http.StatusOK, out)
		return
	}

	// с id -> вернуть одну
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	t, ok := s.tasks[id]
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	writeJSON(w, http.StatusOK, t)
}

// POST /tasks  body: {"title":"..."}
func (s *Store) handlePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "invalid title")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	t := Task{
		ID:    s.nextID,
		Title: req.Title,
		Done:  false,
	}
	s.tasks[t.ID] = t
	s.nextID++

	writeJSON(w, http.StatusCreated, t)
}

// PATCH /tasks?id=1  body: {"done":true}
func (s *Store) handlePatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Done *bool `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Done == nil {
		writeError(w, http.StatusBadRequest, "invalid done")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	t.Done = *req.Done
	s.tasks[id] = t

	// по заданию можно вернуть {"updated": true}
	writeJSON(w, http.StatusOK, map[string]bool{"updated": true})
}

// helpers
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
