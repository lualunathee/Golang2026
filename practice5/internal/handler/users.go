package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"assignment5/internal/repository"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 5
	}

	filters := map[string]string{
		"id":         r.URL.Query().Get("id"),
		"name":       r.URL.Query().Get("name"),
		"email":      r.URL.Query().Get("email"),
		"gender":     r.URL.Query().Get("gender"),
		"birth_date": r.URL.Query().Get("birth_date"),
	}

	orderBy := r.URL.Query().Get("order_by")

	resp, err := h.repo.GetPaginatedUsers(page, limit, filters, orderBy)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user1, err1 := strconv.Atoi(r.URL.Query().Get("user1"))
	user2, err2 := strconv.Atoi(r.URL.Query().Get("user2"))
	if err1 != nil || err2 != nil || user1 <= 0 || user2 <= 0 || user1 == user2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid user ids"})
		return
	}

	users, err := h.repo.GetCommonFriends(user1, user2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(users)
}