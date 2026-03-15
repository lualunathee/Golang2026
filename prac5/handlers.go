package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 5
	}

	resp, _ := h.repo.GetPaginatedUsers(page, pageSize)

	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) CommonFriends(w http.ResponseWriter, r *http.Request) {

	u1, _ := strconv.Atoi(r.URL.Query().Get("user1"))
	u2, _ := strconv.Atoi(r.URL.Query().Get("user2"))

	users, _ := h.repo.GetCommonFriends(u1, u2)

	json.NewEncoder(w).Encode(users)
}
