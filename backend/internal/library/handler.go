package library

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) ListBooks(w http.ResponseWriter, _ *http.Request) {
	books, err := h.svc.GetBooks()
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}
	json.NewEncoder(w).Encode(books)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var b Book

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "invalid request", 400)
		return
	}

	created, err := h.svc.CreateBook(b)
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}
